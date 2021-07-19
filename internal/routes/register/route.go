package register

import (
	"net/http"
	"strconv"

	"github.com/vvkh/social-network/internal/domain/users"
	"github.com/vvkh/social-network/internal/templates"
)

func HandleGet(templates *templates.Templates) http.HandlerFunc {
	render := templates.Add("register.gohtml").Parse()

	return func(writer http.ResponseWriter, request *http.Request) {
		_ = render(writer, nil)
	}
}

func HandlePost(useCase users.UseCase, redirectPath string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// TODO: error handling
		_ = request.ParseForm()

		username := request.Form.Get("username")
		password := request.Form.Get("password")
		firstName := request.Form.Get("first_name")
		lastName := request.Form.Get("last_name")
		age, _ := strconv.Atoi(request.Form.Get("age"))
		location := request.Form.Get("location")
		sex := request.Form.Get("sex")
		about := request.Form.Get("about")

		_, _, err := useCase.CreateUser(request.Context(), username, password, firstName, lastName, uint8(age), location, sex, about)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(writer, request, redirectPath, http.StatusFound)
	}
}
