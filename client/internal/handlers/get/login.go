package get

import (
	"client/internal/lib/log/sl"
	"client/internal/storage/mode"
	"client/internal/views/pages"
	"log/slog"
	"net/http"
)

func LoginHandler(log *slog.Logger, m *mode.ModeStorage) http.HandlerFunc {
	const op = "get.LoginHandler"
	log = log.With(slog.String("op", op))

	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("LoginHandler called")

		c := pages.Login(m.IsDarkMode(r))
		err := c.Render(r.Context(), w)
		if err != nil {
			log.Error("Failed to render page", sl.Err(err))
			return
		}
	}
}
