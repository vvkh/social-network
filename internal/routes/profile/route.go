package profile

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/vvkh/social-network/internal/domain/profiles"
	"github.com/vvkh/social-network/internal/templates"
)

func Handle(profiles profiles.UseCase, templates *templates.Templates) http.HandlerFunc {
	render := templates.Add("profile.gohtml").Parse()

	// TODO: store profile id in jwt or just use GetByUserID?
	return func(writer http.ResponseWriter, request *http.Request) {
		profileID, err := strconv.Atoi(chi.URLParam(request, "profileID"))
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		profile, err := profiles.GetByID(request.Context(), uint64(profileID))
		if err != nil || len(profile) != 1 {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		dto := dtoFromModel(profile[0])
		_ = render(writer, dto)
	}
}
