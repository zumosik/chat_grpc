package get

import (
	"client/internal/storage/mode"
	"log/slog"
	"net/http"
)

// IndexHandler handles the GET / route.
func IndexHandler(log *slog.Logger, m *mode.ModeStorage) http.HandlerFunc {
	const op = "get.IndexHandler"
	log = log.With(slog.String("op", op))

	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("IndexHandler called")

		m.SetMode(w, false)
	}
}
