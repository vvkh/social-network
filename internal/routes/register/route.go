package register

import (
	"net/http"

	"github.com/vvkh/social-network/internal/users"

	"github.com/vvkh/social-network/internal/templates"
)

func HandleGet(templates *templates.Templates) http.HandlerFunc {
	render := templates.Add("register.gohtml").Parse()

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
		_, _, err := useCase.CreateUser(request.Context(), username, password, "John", "Doe", 18, "USA", "male", "")
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		// TODO: redirect to login
		writer.WriteHeader(http.StatusOK)
	}
}
