package login

import (
	"net/http"

	"github.com/vvkh/social-network/internal/templates"
	"github.com/vvkh/social-network/internal/users"
)

func HandleGet(templates *templates.Templates) http.HandlerFunc {
	render := templates.Add("login.gohtml").Parse()

	return func(writer http.ResponseWriter, request *http.Request) {
		_ = render(writer, nil)
	}
}

func HandlePost(useCase users.UseCase) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// TODO: error handling
		_ = request.ParseForm()

		username := request.Form.Get("username")
		password := request.Form.Get("password")
		token, err := useCase.Login(request.Context(), username, password)
		if err != nil {
			writer.WriteHeader(http.StatusForbidden)
			return
		}

		// TODO: extract cookie generation
		cookie := http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
		}
		writer.Header().Set("Set-Cookie", cookie.String())
		writer.WriteHeader(http.StatusOK)
	}
}
