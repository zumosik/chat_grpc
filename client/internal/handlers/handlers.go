package handlers

import (
	"github.com/a-h/templ"
	"net/http"
)

// Render is helper func that renders a component.
func Render(w http.ResponseWriter, r *http.Request, c templ.Component) error {
	return c.Render(r.Context(), w)
}
