package register

import (
	"html/template"
	"net/http"
	"path"
)

func Handle(templatesDir string) http.HandlerFunc {
	t := template.Must(template.ParseFiles(
		path.Join(templatesDir, "base.gohtml"),
		path.Join(templatesDir, "register.gohtml"),
	))
	return func(writer http.ResponseWriter, request *http.Request) {
		_ = t.Execute(writer, nil)
	}
}
