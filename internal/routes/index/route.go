package index

import "net/http"

func HandleIndex() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(`hello world`))
	}
}
