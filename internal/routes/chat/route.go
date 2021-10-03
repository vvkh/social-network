package chat

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/vvkh/social-network/internal/domain/chats"
	"github.com/vvkh/social-network/internal/domain/profiles"
	profilesEntity "github.com/vvkh/social-network/internal/domain/profiles/entity"
	"github.com/vvkh/social-network/internal/middlewares"
	"github.com/vvkh/social-network/internal/navbar"
	"github.com/vvkh/social-network/internal/templates"
)

const (
	mostPossibleProfilesCountPerChat = 2
)

func Handle(log *zap.SugaredLogger, profilesUseCase profiles.UseCase, chatsUseCase chats.UseCase, navbar *navbar.Navbar, templates *templates.Templates) http.HandlerFunc {
	render := templates.Add("chat.gohtml").Parse()

	return func(writer http.ResponseWriter, request *http.Request) {
		chatID, err := strconv.Atoi(chi.URLParam(request, "chatID"))
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		self, _ := middlewares.ProfileFromCtx(request.Context())

		chat, messages, err := chatsUseCase.ListChatMessages(request.Context(), self.ID, uint64(chatID))
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		allAuthorsIDs := make([]uint64, 0, mostPossibleProfilesCountPerChat)
		for _, message := range messages {
			allAuthorsIDs = append(allAuthorsIDs, message.AuthorProfileID)
		}

		var profiles []profilesEntity.Profile

		if len(allAuthorsIDs) > 0 {
			profiles, err = profilesUseCase.GetByID(request.Context(), allAuthorsIDs...)
			if err != nil {
				log.Warnw("failed to fetch profile names, falling back to showing just ids", "err", err)
				profiles = []profilesEntity.Profile{}
			}
		}

		templateCtx := Contex{
			Navbar:   navbar.GetContext(request.Context()),
			Chat:     chatDroFromModel(chat),
			Messages: messagesDtoFromModel(messages, profiles),
		}
		err = render(writer, templateCtx)
		if err != nil {
			log.Warnw("error while rending template", "err", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
