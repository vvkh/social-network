package cookies

import (
	"net/http"
	"time"
)

const (
	key = "token"
)

var (
	EmptyAuthCookie = authCookie("", time.Unix(0, 0)) // expire immediately
)

func AuthCookie(token string) *http.Cookie {
	return authCookie(token, time.Time{})
}

func ReadAuthCookie(r *http.Request) (*http.Cookie, error) {
	return r.Cookie(key)
}

func authCookie(token string, expiredAt time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     key,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Expires:  expiredAt,
	}
}
