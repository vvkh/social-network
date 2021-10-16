package message

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/vvkh/social-network/internal/domain/chats"
	"github.com/vvkh/social-network/internal/middlewares"
)

func HandlePost(log *zap.SugaredLogger, chatsUseCase chats.UseCase) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		chatID, err := strconv.Atoi(chi.URLParam(request, "chatID"))
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		err = request.ParseForm()
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		messageContent := request.Form.Get("message")
		if messageContent == "" {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		self, _ := middlewares.ProfileFromCtx(request.Context())
		err = chatsUseCase.SendMessage(request.Context(), self.ID, uint64(chatID), messageContent)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(writer, request, request.URL.String(), http.StatusFound)
	}
}
