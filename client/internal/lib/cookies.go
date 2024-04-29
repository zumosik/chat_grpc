package lib

import (
	"encoding/base64"
	"errors"
	"net/http"
)

var (
	ErrTooLong      = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
)

func WriteCookie(w http.ResponseWriter, c http.Cookie) error {
	c.Value = base64.URLEncoding.EncodeToString([]byte(c.Value))

	if len(c.String()) > 4096 {
		return ErrTooLong
	}

	http.SetCookie(w, &c)

	return nil
}

func ReadCookie(r *http.Request, name string) (string, error) {
	c, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	value, err := base64.URLEncoding.DecodeString(c.Value)
	if err != nil {
		return "", ErrInvalidValue
	}

	return string(value), nil
}
