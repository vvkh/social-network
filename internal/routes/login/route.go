package login

import (
	"fmt"
	"net/http"

	"github.com/vvkh/social-network/internal/users"

	"github.com/vvkh/social-network/internal/templates"
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

		// TODO: http only cookies
		writer.Header().Set("Set-Cookie", fmt.Sprintf("token=%s", token))
		writer.WriteHeader(http.StatusOK)
	}
}
