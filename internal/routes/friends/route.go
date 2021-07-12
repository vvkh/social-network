package friends

import (
	"net/http"

	"github.com/vvkh/social-network/internal/templates"
)

func Handle(templates *templates.Templates) http.HandlerFunc {
	render := templates.Add("friends.gohtml").Parse()

	return func(writer http.ResponseWriter, request *http.Request) {
		_ = render(writer, nil)
	}
}
