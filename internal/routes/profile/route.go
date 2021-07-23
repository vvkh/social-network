package profile

import (
	fmt "fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/vvkh/social-network/internal/domain/friendship"
	"github.com/vvkh/social-network/internal/domain/profiles"
	"github.com/vvkh/social-network/internal/middlewares"
	"github.com/vvkh/social-network/internal/templates"
)

func Handle(profiles profiles.UseCase, templates *templates.Templates) http.HandlerFunc {
	render := templates.Add("profile.gohtml").Parse()

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
		profileDto := dtoFromModel(profile[0])

		self, _ := middlewares.ProfileFromCtx(request.Context())
		selfProfileDto := dtoFromModel(self)
		err = render(writer, Context{
			Self:    selfProfileDto,
			Profile: profileDto,
		})
		fmt.Println(err)
	}
}

func HandlePost(friendshipUseCase friendship.UseCase) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		profileID, err := strconv.Atoi(chi.URLParam(request, "profileID"))
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		self, _ := middlewares.ProfileFromCtx(request.Context())
		err = friendshipUseCase.CreateRequest(request.Context(), self.ID, uint64(profileID))
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.Redirect(writer, request, fmt.Sprintf("/profiles/%d/", profileID), http.StatusFound)
	}
}
