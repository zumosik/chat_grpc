package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"rooms_service/internal/client/auth"
	"rooms_service/internal/config"
	"rooms_service/internal/interceptor"
	"rooms_service/internal/lib/logger/slogpretty"
	"rooms_service/internal/service"
	"rooms_service/internal/storage/sql/postgres"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	// open db
	db, err := sql.Open("postgres", cfg.Storage.PostgresURl)
	if err != nil {
		panic(fmt.Sprintf("failed to open db (%s): %v", cfg.Storage.PostgresURl, err))
	}
	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintf("failed to open db (%s): %v", cfg.Storage.PostgresURl, err))
	}

	// create log
	log := setupLogger(cfg.Env)

	// connect to auth service
	authClient, err := auth.Connect(log, cfg.OtherServices.PrivateAuthServiceURL, &cfg.OtherServices.PrivateAuthServiceCert)
	if err != nil {
		log.Error("failed to connect to auth service", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// create storage
	storage := postgres.New(db, authClient)

	// create grpc server
	serv := service.New(log, storage)

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
		interceptor.TokenMiddleware(authClient),
	), grpc.Creds(creds))

	service.Register(gRPCServer, serv)

	// Start gRPC server
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
	log.Info("Gracefully stopped service")
}

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
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
