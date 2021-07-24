package friends_requests

import (
	"net/http"

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
