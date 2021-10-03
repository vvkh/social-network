package requests_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vvkh/social-network/internal/domain/friendship/mocks"
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
			url:        "/friends/requests/1/create/",
			wantStatus: http.StatusFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			friendshipUseCase := mocks.NewMockUseCase(ctrl)
			friendshipUseCase.EXPECT().CreateRequest(gomock.Any(), test.self.ID, test.profile.ID).Return(nil)
			log, _ := zap.NewDevelopment()
			s := server.New(log.Sugar(), ":80", "../../../../templates", nil, nil, friendshipUseCase, nil)

			request := httptest.NewRequest("POST", test.url, nil)
			request = request.WithContext(middlewares.AddProfileToCtx(request.Context(), test.self))

			responseWriter := httptest.NewRecorder()
			s.Handle(responseWriter, request)

			response := responseWriter.Result()
			require.Equal(t, test.wantStatus, response.StatusCode)
		})
	}
}
