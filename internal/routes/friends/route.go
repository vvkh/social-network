package friends

import (
	"net/http"

	"github.com/vvkh/social-network/internal/domain/friendship"
	"github.com/vvkh/social-network/internal/middlewares"
	"github.com/vvkh/social-network/internal/navbar"
	"github.com/vvkh/social-network/internal/templates"
)

func Handle(friendshipUseCase friendship.UseCase, navbar *navbar.Navbar, templates *templates.Templates) http.HandlerFunc {
	render := templates.Add("friends.gohtml").Parse()

	return func(writer http.ResponseWriter, request *http.Request) {
		self, _ := middlewares.ProfileFromCtx(request.Context())
		friends, err := friendshipUseCase.ListFriends(request.Context(), self.ID)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		pendingFriendshipRequests, err := friendshipUseCase.ListPendingRequests(request.Context(), self.ID)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		templateCtx := Contex{
			Navbar:               navbar.GetContext(request.Context()),
			Self:                 dtoFromModel(self),
			Friends:              dtoFromModels(friends),
			PendingRequestsCount: len(pendingFriendshipRequests),
		}
		err = render(writer, templateCtx)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
