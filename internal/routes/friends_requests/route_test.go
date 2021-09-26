package friends_requests_test

import (
	"io"
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

func TestFriendshipRequestPage(t *testing.T) {
	tests := []struct {
		name               string
		self               entity.Profile
		friendshipRequests []entity.Profile
		wantBody           []string
	}{
		{
			name: "friendship_requests_are_shown_on_page",
			self: entity.Profile{
				ID:     1,
				UserID: 2,
			},
			friendshipRequests: []entity.Profile{
				{
					ID:        3,
					UserID:    4,
					FirstName: "John",
					LastName:  "Doe",
				},
				{
					ID:        5,
					UserID:    6,
					FirstName: "Topsy",
					LastName:  "Cret",
				},
			},
			wantBody: []string{
				"<h1>Friendship requests</h1>",
				`<a href="/profiles/3/">John Doe</a>`,
				`<form method="POST" action="/friends/requests/3/accept/"><input type="submit" value="Accept"></form>`,
				`<form method="POST" action="/friends/requests/3/decline/"><input type="submit" value="Decline"></form>`,
				`<a href="/profiles/5/">Topsy Cret</a>`,
				`<form method="POST" action="/friends/requests/5/accept/"><input type="submit" value="Accept"></form>`,
				`<form method="POST" action="/friends/requests/5/decline/"><input type="submit" value="Decline"></form>`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			friendshipUseCase := mocks.NewMockUseCase(ctrl)
			friendshipUseCase.EXPECT().ListPendingRequests(gomock.Any(), test.self.ID).Return(test.friendshipRequests, nil).AnyTimes()

			request := httptest.NewRequest("GET", "/friends/requests/", nil)
			request = request.WithContext(middlewares.AddProfileToCtx(request.Context(), test.self))
			responseWriter := httptest.NewRecorder()
			log, _ := zap.NewDevelopment()
			s := server.New(log.Sugar(), ":80", "../../../templates", nil, nil, friendshipUseCase)
			s.Handle(responseWriter, request)

			response := responseWriter.Result()
			require.Equal(t, http.StatusOK, response.StatusCode)

			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			for _, wantPart := range test.wantBody {
				require.Contains(t, string(body), wantPart)
			}
		})
	}
}
