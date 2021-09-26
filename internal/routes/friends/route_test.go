package friends_test

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

func TestFriendsPage(t *testing.T) {
	tests := []struct {
		name               string
		self               entity.Profile
		friendshipRequests []entity.Profile
		friends            []entity.Profile
		wantBodyParts      []string
	}{
		{
			name: "friends_are_shown_on_page",
			self: entity.Profile{
				ID:     1,
				UserID: 2,
			},
			friends: []entity.Profile{
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
			wantBodyParts: []string{
				"<h1>Friends</h1>",
				`<a href="/profiles/3/">John Doe</a>`,
				`<a href="/profiles/5/">Topsy Cret</a>`,
			},
		},
		{
			name: "pending_requests_count_shown_on_page",
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
			wantBodyParts: []string{
				"<h1>Friends</h1>",
				`<a href="/friends/requests/">Pending friendship requests (2)</a>`,
			},
		},
		{
			name: "pending_requests_count_shown_in_the_navbar_if_no_requests",
			self: entity.Profile{
				ID:     1,
				UserID: 2,
			},
			friendshipRequests: []entity.Profile{},
			wantBodyParts: []string{
				`<a href="/friends/">Friends</a>`,
			},
		},
		{
			name: "pending_requests_count_shown_in_the_navbar",
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
			wantBodyParts: []string{
				`<a href="/friends/">Friends (2)</a>`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			friendshipUseCase := mocks.NewMockUseCase(ctrl)
			friendshipUseCase.EXPECT().ListFriends(gomock.Any(), test.self.ID).Return(test.friends, nil)
			friendshipUseCase.EXPECT().ListPendingRequests(gomock.Any(), test.self.ID).Return(test.friendshipRequests, nil).AnyTimes()
			log, _ := zap.NewDevelopment()
			s := server.New(log.Sugar(), ":80", "../../../templates", nil, nil, friendshipUseCase)

			request := httptest.NewRequest("GET", "/friends/", nil)
			request = request.WithContext(middlewares.AddProfileToCtx(request.Context(), test.self))
			responseWriter := httptest.NewRecorder()

			s.Handle(responseWriter, request)

			response := responseWriter.Result()
			require.Equal(t, http.StatusOK, response.StatusCode)

			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			for _, wantPart := range test.wantBodyParts {
				require.Contains(t, string(body), wantPart)
			}
		})
	}
}
