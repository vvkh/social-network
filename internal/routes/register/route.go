package register

import (
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/vvkh/social-network/internal/domain/users"
	"github.com/vvkh/social-network/internal/templates"
)

func HandleGet(templates *templates.Templates) http.HandlerFunc {
	render := templates.Add("register.gohtml").Parse()

	return func(writer http.ResponseWriter, request *http.Request) {
		err := render(writer, nil)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func HandlePost(log *zap.SugaredLogger, useCase users.UseCase, redirectPath string, templates *templates.Templates) http.HandlerFunc {
	render := templates.Add("register.gohtml").Parse()
	return func(writer http.ResponseWriter, request *http.Request) {
		if err := request.ParseForm(); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_ = render(writer, Dto{
				Error: "bad form data",
			})
			return
		}

		age, err := strconv.Atoi(request.Form.Get("age"))
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_ = render(writer, Dto{
				Error: "age must be a positive number",
			})
			return
		}
		username := request.Form.Get("username")
		password := request.Form.Get("password")
		firstName := request.Form.Get("first_name")
		lastName := request.Form.Get("last_name")
		location := request.Form.Get("location")
		sex := request.Form.Get("sex")
		about := request.Form.Get("about")

		_, _, err = useCase.CreateUser(request.Context(), username, password, firstName, lastName, uint8(age), location, sex, about)
		if err != nil {
			log.Errorw("error while creating user", "err", err)
			writer.WriteHeader(http.StatusInternalServerError)
			_ = render(writer, Dto{
				Error: "user creation failed",
			})
			return
		}

		http.Redirect(writer, request, redirectPath, http.StatusFound)
	}
}
