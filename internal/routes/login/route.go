package login

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/vvkh/social-network/internal/cookies"
	"github.com/vvkh/social-network/internal/domain/users"
	"github.com/vvkh/social-network/internal/templates"
)

func HandleGet(templates *templates.Templates) http.HandlerFunc {
	render := templates.Add("login.gohtml").Parse()

	return func(writer http.ResponseWriter, request *http.Request) {
		_ = render(writer, nil)
	}
}

func HandlePost(log *zap.SugaredLogger, useCase users.UseCase, redirectPath string, templates *templates.Templates) http.HandlerFunc {
	render := templates.Add("login.gohtml").Parse()

	return func(writer http.ResponseWriter, request *http.Request) {
		// TODO: add CSRF token
		if err := request.ParseForm(); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_ = render(writer, Dto{
				Error: "bad form sent",
			})
			return
		}

		username := request.Form.Get("username")
		password := request.Form.Get("password")
		token, err := useCase.Login(request.Context(), username, password)
		if errors.Is(err, users.AuthenticationFailed) {
			writer.WriteHeader(http.StatusForbidden)
			_ = render(writer, Dto{
				Error: "authentication failed",
			})
			return
		}
		if errors.Is(err, users.EmptyCredentials) {
			writer.WriteHeader(http.StatusBadRequest)
			_ = render(writer, Dto{
				Error: "both password and username are required",
			})
			return
		}
		if err != nil {
			log.Errorw("error while performing login", "err", err)
			writer.WriteHeader(http.StatusInternalServerError)
			_ = render(writer, Dto{
				Error: "server failed",
			})
			return
		}

		http.SetCookie(writer, cookies.AuthCookie(token))
		http.Redirect(writer, request, redirectPath, http.StatusFound)
	}
}
