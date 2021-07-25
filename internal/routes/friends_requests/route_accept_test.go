package friends_requests_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/vvkh/social-network/internal/domain/friendship/mocks"
	"github.com/vvkh/social-network/internal/domain/profiles/entity"
	"github.com/vvkh/social-network/internal/middlewares"
	"github.com/vvkh/social-network/internal/server"
)

func TestAcceptFriendshipRequest(t *testing.T) {
	tests := []struct {
		name      string
		self      entity.Profile
		url       string
		profileID uint64
	}{
		{
			name: "accept_request",
			self: entity.Profile{
				ID:     1,
				UserID: 2,
			},
			profileID: 3,
			url:       "/friends/requests/3/accept/",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			friendshipUseCase := mocks.NewMockUseCase(ctrl)
			friendshipUseCase.EXPECT().AcceptRequest(gomock.Any(), test.profileID, test.self.ID).Return(nil)

			s := server.New(":80", "../../../templates", nil, nil, friendshipUseCase)
			request := httptest.NewRequest("POST", test.url, nil)
			request = request.WithContext(middlewares.AddProfileToCtx(request.Context(), test.self))
			responseWriter := httptest.NewRecorder()

			s.Handle(responseWriter, request)

			response := responseWriter.Result()
			require.Equal(t, http.StatusFound, response.StatusCode)
			require.Equal(t, "/friends/requests/", response.Header.Get("Location"))
		})
	}
}
