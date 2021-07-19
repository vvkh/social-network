package logout

import (
	"net/http"

	"github.com/vvkh/social-network/internal/cookies"
)

func HandleGet(redirectPath string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		http.SetCookie(writer, cookies.EmptyAuthCookie)
		http.Redirect(writer, request, redirectPath, http.StatusFound)
	}
}
