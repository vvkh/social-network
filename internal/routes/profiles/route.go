package profiles

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/vvkh/social-network/internal/domain/profiles"
	"github.com/vvkh/social-network/internal/middlewares"
	"github.com/vvkh/social-network/internal/templates"
)

func Handle(log *zap.SugaredLogger, profilesUseCase profiles.UseCase, templates *templates.Templates) http.HandlerFunc {
	render := templates.Add("profiles.gohtml").Parse()

	return func(writer http.ResponseWriter, request *http.Request) {
		self, _ := middlewares.ProfileFromCtx(request.Context())

		profiles, err := profilesUseCase.ListProfiles(request.Context())
		if err != nil {
			log.Errorw("error while listing profiles", "err", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		context := Context{
			Self:     dtoFromModel(self),
			Profiles: dtoFromModels(profiles),
		}

		err = render(writer, context)
		if err != nil {
			log.Errorw("error while rendering profiles", "err", err)
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}
}
