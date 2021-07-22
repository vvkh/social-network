package index

import (
	"fmt"
	"net/http"

	"github.com/vvkh/social-network/internal/middlewares"
)

func Handle() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		profile, _ := middlewares.ProfileFromCtx(request.Context())
		http.Redirect(writer, request, fmt.Sprintf("/profiles/%d/", profile.ID), http.StatusFound)
	}
}
