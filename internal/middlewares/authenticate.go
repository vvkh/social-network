package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/vvkh/social-network/internal/domain/users"
)

const (
	CookieKeyToken = "token"
	CookiePath     = "/"
	CtxKeyToken    = 1
)

func AuthenticateUser(users users.UseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			encodedToken, err := r.Cookie(CookieKeyToken)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := r.Context()
			token, err := users.DecodeToken(ctx, encodedToken.Value)
			if err != nil {
				// TODO: extract cookie
				cookie := http.Cookie{
					Name:     CookieKeyToken,
					Path:     CookiePath,
					HttpOnly: true,
					Expires:  time.Unix(0, 0), // expire immediately
				}
				w.Header().Set("Set-Cookie", cookie.String())
				next.ServeHTTP(w, r)
				return
			}

			ctx = context.WithValue(ctx, CtxKeyToken, token)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
