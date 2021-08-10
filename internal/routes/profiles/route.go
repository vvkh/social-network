package profiles

import (
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/vvkh/social-network/internal/domain/profiles"
	"github.com/vvkh/social-network/internal/middlewares"
	"github.com/vvkh/social-network/internal/templates"
)

const (
	defaultLimit = 10
)

func Handle(log *zap.SugaredLogger, profilesUseCase profiles.UseCase, templates *templates.Templates) http.HandlerFunc {
	render := templates.Add("profiles.gohtml").Parse()

	return func(writer http.ResponseWriter, request *http.Request) {
		self, _ := middlewares.ProfileFromCtx(request.Context())

		if err := request.ParseForm(); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		firstNameFilter := request.Form.Get("first_name")
		lastNamePrefix := request.Form.Get("last_name")
		limit, err := strconv.Atoi(request.Form.Get("limit"))
		if err != nil {
			limit = defaultLimit
		}

		profiles, err := profilesUseCase.ListProfiles(request.Context(), firstNameFilter, lastNamePrefix, limit)

		if err != nil {
			log.Errorw("error while listing profiles", "err", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		context := Context{
			Self:     dtoFromModel(self),
			Profiles: dtoFromModels(profiles),
			Filters: Filters{
				FirstName: firstNameFilter,
				LastName:  lastNamePrefix,
			},
		}

		err = render(writer, context)
		if err != nil {
			log.Errorw("error while rendering profiles", "err", err)
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}
}
