package friend_test

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

func TestStopFriendship(t *testing.T) {
	tests := []struct {
		name         string
		self         entity.Profile
		url          string
		wantRedirect string
		profileID    uint64
	}{
		{
			name: "stop_friendship",
			self: entity.Profile{
				ID:     1,
				UserID: 2,
			},
			profileID:    3,
			url:          "/friends/3/stop/",
			wantRedirect: "/profiles/3/",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			friendshipUseCase := mocks.NewMockUseCase(ctrl)
			friendshipUseCase.EXPECT().StopFriendship(gomock.Any(), test.self.ID, test.profileID).Return(nil)
			log, _ := zap.NewDevelopment()
			s := server.New(log.Sugar(), ":80", "../../../templates", nil, nil, friendshipUseCase, nil)
			request := httptest.NewRequest("POST", test.url, nil)
			request = request.WithContext(middlewares.AddProfileToCtx(request.Context(), test.self))
			responseWriter := httptest.NewRecorder()

			s.Handle(responseWriter, request)

			response := responseWriter.Result()
			require.Equal(t, http.StatusFound, response.StatusCode)
			require.Equal(t, test.wantRedirect, response.Header.Get("Location"))
		})
	}
}
