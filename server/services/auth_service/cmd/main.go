package main

import (
	"auth_service/internal/client/notifications"
	"auth_service/internal/config"
	"auth_service/internal/lib/logger/slogpretty"
	"auth_service/internal/lib/token"
	"auth_service/internal/service"
	"auth_service/internal/storage/sql/postgres"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc/credentials"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	_ "github.com/lib/pq"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	// open db
	db, err := sqlx.Open("postgres", cfg.Storage.PostgresURl)
	if err != nil {
		panic(fmt.Sprintf("failed to open db (%s): %v", cfg.Storage.PostgresURl, err))
	}
	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintf("failed to open db (%s): %v", cfg.Storage.PostgresURl, err))
	}

	// create storage and log
	storage := postgres.New(db)
	log := setupLogger(cfg.Env)
	tokenManager := token.NewManager(cfg.Tokens.TokenSecret, cfg.Tokens.TokenTTL)

	// connect to another services
	notificationManager, err := notifications.Connect(log, cfg.OtherServices.NotificationServiceURL, &cfg.OtherServices.NotificationsCert)
	if err != nil {
		log.Error("failed to connect to notifications service", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// create public auth service
	s := service.New(log, storage, storage, tokenManager, notificationManager)

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
	), grpc.Creds(creds))

	service.Register(gRPCServer, s)

	// Start gRPC server
	go func() {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
		if err != nil {
			log.Error("failed to listen", slog.String("error", err.Error()))
		}

		log.Info("Starting public gRPC server", slog.String("port", fmt.Sprintf(":%d", cfg.GRPC.Port)))

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
