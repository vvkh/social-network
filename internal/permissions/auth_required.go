package permissions

import (
	"net/http"

	"github.com/vvkh/social-network/internal/middlewares"
)

func AuthRequired(redirectPath string) func(http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if _, ok := middlewares.ProfileFromCtx(r.Context()); !ok {
				http.Redirect(w, r, redirectPath, http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
