package middlewares

import (
	"context"
	"net/http"

	"github.com/vvkh/social-network/internal/cookies"
	"github.com/vvkh/social-network/internal/domain/users"
)

type ctxKey int

const (
	CtxKeyToken = ctxKey(1)
)

func AuthenticateUser(users users.UseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			encodedToken, err := cookies.ReadAuthCookie(r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := r.Context()
			token, err := users.DecodeToken(ctx, encodedToken.Value)
			if err != nil {
				http.SetCookie(w, cookies.EmptyAuthCookie)
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
