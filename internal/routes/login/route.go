package login

import (
	"html/template"
	"net/http"
)

func Handle(templates *template.Template) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_ = templates.ExecuteTemplate(writer, "login.gohtml", nil)
	}
}
