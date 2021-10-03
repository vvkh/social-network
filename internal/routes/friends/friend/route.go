package friend

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/vvkh/social-network/internal/domain/friendship"
	"github.com/vvkh/social-network/internal/middlewares"
)

func HandleStop(friendshipUseCase friendship.UseCase) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		profileID, err := strconv.Atoi(chi.URLParam(request, "profile"))
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		self, _ := middlewares.ProfileFromCtx(request.Context())
		err = friendshipUseCase.StopFriendship(request.Context(), self.ID, uint64(profileID))
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.Redirect(writer, request, fmt.Sprintf("/profiles/%d/", profileID), http.StatusFound)
	}
}
