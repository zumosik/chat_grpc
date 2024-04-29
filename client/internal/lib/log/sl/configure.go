package sl

import (
	"github.com/phsym/console-slog"
	"log/slog"
	"os"
)

const (
	Local = "local"
	Dev   = "dev"
	Prod  = "prod"
)

// MustConfigureLogger configures a logger based on the environment.
// Local: console handler with debug level.
// Dev: JSON handler with info level.
// Prod: JSON handler with error level.
func MustConfigureLogger(env string) *slog.Logger {
	var handler slog.Handler

	switch env {
	case Local:
		handler = console.NewHandler(os.Stdout, &console.HandlerOptions{Level: slog.LevelDebug})
	case Dev:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	case Prod:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelError,
		})
	default:
		panic("unknown environment")
	}

	return slog.New(handler)
}
