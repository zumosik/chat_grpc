package mode

import (
	"client/internal/lib"
	"net/http"
)

const (
	cookieName = "mode"
	darkMode   = "dark"
	lightMode  = "light"
)

type ModeStorage struct {
}

func NewModeStorage() *ModeStorage {
	return &ModeStorage{}
}

func (m *ModeStorage) SwitchMode(r *http.Request, w http.ResponseWriter) error {
	isDark := m.IsDarkMode(r)
	return m.SetMode(w, !isDark)
}

func (m *ModeStorage) IsDarkMode(r *http.Request) bool {
	val, err := lib.ReadCookie(r, cookieName)
	if err != nil {
		return false
	}

	return val == darkMode
}

func (m *ModeStorage) SetMode(w http.ResponseWriter, isDark bool) error {
	val := lightMode
	if isDark {
		val = darkMode
	}

	return lib.WriteCookie(w, http.Cookie{
		Name:  cookieName,
		Value: val,
	})
}
