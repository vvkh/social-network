package friends_requests

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/vvkh/social-network/internal/domain/friendship"
	"github.com/vvkh/social-network/internal/middlewares"
	"github.com/vvkh/social-network/internal/templates"
)

func Handle(friendshipUseCase friendship.UseCase, templates *templates.Templates) http.HandlerFunc {
	render := templates.Add("friends_requests.gohtml").Parse()

	return func(writer http.ResponseWriter, request *http.Request) {
		self, _ := middlewares.ProfileFromCtx(request.Context())
		pendingFriendshipRequests, err := friendshipUseCase.ListPendingRequests(request.Context(), self.ID)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		tempalateCtx := Contex{
			Self:            dtoFromModel(self),
			PendingRequests: dtoFromModels(pendingFriendshipRequests),
		}
		err = render(writer, tempalateCtx)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func HandlePostAccept(friendshipUseCase friendship.UseCase, redirect string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		self, _ := middlewares.ProfileFromCtx(request.Context())
		profileFrom, err := strconv.Atoi(chi.URLParam(request, "profileFrom"))
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		err = friendshipUseCase.AcceptRequest(request.Context(), uint64(profileFrom), self.ID)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.Redirect(writer, request, redirect, http.StatusFound)
	}
}

func HandlePostDecline(friendshipUseCase friendship.UseCase, redirect string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		self, _ := middlewares.ProfileFromCtx(request.Context())
		profileFrom, err := strconv.Atoi(chi.URLParam(request, "profileFrom"))
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		err = friendshipUseCase.DeclineRequest(request.Context(), uint64(profileFrom), self.ID)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.Redirect(writer, request, redirect, http.StatusFound)
	}
}
