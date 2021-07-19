package permissions

import (
	"net/http"

	"github.com/vvkh/social-network/internal/domain/users/entity"
	"github.com/vvkh/social-network/internal/middlewares"
)

func AuthRequired(redirectPath string) func(http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token := r.Context().Value(middlewares.CtxKeyToken)
			if _, ok := token.(entity.AccessToken); !ok {
				http.Redirect(w, r, redirectPath, http.StatusFound)
				return
			}
			next.ServeHTTP(w, r)
		}
	}
}
