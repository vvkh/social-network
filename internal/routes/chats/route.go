package chats

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/vvkh/social-network/internal/domain/chats"
	"github.com/vvkh/social-network/internal/middlewares"
	"github.com/vvkh/social-network/internal/navbar"
	"github.com/vvkh/social-network/internal/templates"
)

func Handle(log *zap.SugaredLogger, chatsUseCase chats.UseCase, navbar *navbar.Navbar, templates *templates.Templates) http.HandlerFunc {
	render := templates.Add("chats.gohtml").Parse()

	return func(writer http.ResponseWriter, request *http.Request) {
		self, _ := middlewares.ProfileFromCtx(request.Context())
		chats, err := chatsUseCase.ListChats(request.Context(), self.ID)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		templateCtx := Contex{
			Navbar: navbar.GetContext(request.Context()),
			Chats:  dtoFromModels(chats),
		}
		err = render(writer, templateCtx)
		if err != nil {
			log.Warnw("error while rending template", "err", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
