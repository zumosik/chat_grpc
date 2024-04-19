package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc/credentials"
	"log/slog"
	"net"
	"notifications_service/internal/config"
	"notifications_service/internal/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/gomail.v2"
)

// TODO: add credentials
func main() {
	cfg := config.MustLoad()

	d := gomail.NewDialer(cfg.Email.SMTP, cfg.Email.Port, cfg.Email.From, cfg.Email.Password)
	log := slog.Default() // TODO

	s := service.New(cfg.Email.From, d, log)

	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			//logging.StartCall, logging.FinishCall,
			logging.PayloadReceived, logging.PayloadSent,
		),
		// Add any other option (check functions starting with logging.With).
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	creds := mustLoadTLSCreds(cfg.GRPC.Certs.CaPath, cfg.GRPC.Certs.CertPath, cfg.GRPC.Certs.KeyPath)

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
	),
		grpc.Creds(creds),
	)

	service.Register(gRPCServer, s)

	go func() {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
		if err != nil {
			log.Error("failed to listen", slog.String("error", err.Error()))
		}

		log.Info("Starting gRPC server", slog.String("port", fmt.Sprintf(":%d", cfg.GRPC.Port)))

		if err := gRPCServer.Serve(l); err != nil {
			log.Error("failed to serve", slog.String("error", err.Error()))
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	gRPCServer.GracefulStop()
	log.Info("Gracefully stopped")
}

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func mustLoadTLSCreds(ca, crt, key string) credentials.TransportCredentials {
	pemClientCA, err := os.ReadFile(ca)
	if err != nil {
		panic(fmt.Sprintf("failed to read certificate file: %v", err))
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		panic("failed to append certificate to pool")
	}

	creds, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		panic(fmt.Sprintf("failed to load TLS keys: %v", err))
	}

	c := &tls.Config{
		Certificates: []tls.Certificate{creds},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(c)
}
