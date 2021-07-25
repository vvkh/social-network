package profile_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	friendshipMocks "github.com/vvkh/social-network/internal/domain/friendship/mocks"
	"github.com/vvkh/social-network/internal/domain/profiles/entity"
	"github.com/vvkh/social-network/internal/middlewares"
	"github.com/vvkh/social-network/internal/server"
)

func TestProfilePageFriendshipRequest(t *testing.T) {
	tests := []struct {
		name       string
		profile    entity.Profile
		self       entity.Profile
		url        string
		wantStatus int
	}{
		{
			name: "submit_friendship_form",
			profile: entity.Profile{
				ID:     1,
				UserID: 2,
			},
			self: entity.Profile{
				ID:     3,
				UserID: 4,
			},
			url:        "/profiles/1/friendship/",
			wantStatus: http.StatusFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			friendshipUseCase := friendshipMocks.NewMockUseCase(ctrl)
			friendshipUseCase.EXPECT().CreateRequest(gomock.Any(), test.self.ID, test.profile.ID).Return(nil)

			s := server.New(":80", "../../../templates", nil, nil, friendshipUseCase)

			request := httptest.NewRequest("POST", test.url, nil)
			request = request.WithContext(middlewares.AddProfileToCtx(request.Context(), test.self))

			responseWriter := httptest.NewRecorder()
			s.Handle(responseWriter, request)

			response := responseWriter.Result()
			require.Equal(t, test.wantStatus, response.StatusCode)
		})
	}
}
