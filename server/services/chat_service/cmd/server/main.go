package main

import (
	"chat_service/internal/client/auth"
	"chat_service/internal/client/rooms"
	"chat_service/internal/config"
	"chat_service/internal/lib/logger/logger/slogpretty"
	"chat_service/internal/service"
	redisSt "chat_service/internal/storage/redis"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	cl := redis.NewClient(&redis.Options{
		Addr:     cfg.Storage.RedisURL,
		Password: cfg.Storage.Password,
		DB:       cfg.Storage.Db,
	})

	storage := redisSt.New(cl)
	log := setupLogger(cfg.Env)

	authService, err := auth.Connect(log, cfg.OtherServices.AuthService, &cfg.OtherServices.AuthCerts)
	if err != nil {
		log.Error("failed to connect to auth service", slog.String("error", err.Error()))
		os.Exit(1)
	}

	roomsService, err := rooms.Connect(log, cfg.OtherServices.RoomsService, &cfg.OtherServices.RoomsCerts)
	if err != nil {
		log.Error("failed to connect to rooms service", slog.String("error", err.Error()))
		os.Exit(1)
	}

	s := service.New(log, authService, roomsService, storage)

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
