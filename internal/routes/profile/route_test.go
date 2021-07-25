package profile_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	friendshipEntity "github.com/vvkh/social-network/internal/domain/friendship/entity"
	friendshipMocks "github.com/vvkh/social-network/internal/domain/friendship/mocks"
	"github.com/vvkh/social-network/internal/domain/profiles/entity"
	profilesMocks "github.com/vvkh/social-network/internal/domain/profiles/mocks"
	"github.com/vvkh/social-network/internal/middlewares"
	"github.com/vvkh/social-network/internal/server"
)

func TestProfilePage(t *testing.T) {
	tests := []struct {
		name             string
		profile          entity.Profile
		self             entity.Profile
		url              string
		friendshipStatus friendshipEntity.FriendshipStatus
		wantStatus       int
		wantBody         []string
	}{
		{
			name: "profile_page_contains_friendship_request",
			profile: entity.Profile{
				ID:     1,
				UserID: 2,
			},
			self: entity.Profile{
				ID:     3,
				UserID: 4,
			},
			url: "/profiles/1/",
			friendshipStatus: friendshipEntity.FriendshipStatus{
				State: friendshipEntity.StateNone,
			},
			wantStatus: http.StatusOK,
			wantBody: []string{
				`<form method="POST" action="/friends/requests/1/create/"><input type="submit" value="request friendship"></form>`,
			},
		},
		{
			name: "has_accept_decline_request_buttons_for_pending_request",
			profile: entity.Profile{
				ID:     1,
				UserID: 2,
			},
			self: entity.Profile{
				ID:     3,
				UserID: 4,
			},
			url: "/profiles/1/",
			friendshipStatus: friendshipEntity.FriendshipStatus{
				RequestedToProfileID:   3,
				RequestedFromProfileID: 1,
				State:                  friendshipEntity.StatePending,
			},
			wantStatus: http.StatusOK,
			wantBody: []string{
				`<form method="POST" action="/friends/requests/1/accept/"><input type="submit" value="Accept"></form>`,
				`<form method="POST" action="/friends/requests/1/decline/"><input type="submit" value="Decline"></form>`,
			},
		},
		{
			name: "waiting_for_friendship_confirmation_button",
			profile: entity.Profile{
				ID:     1,
				UserID: 2,
			},
			self: entity.Profile{
				ID:     3,
				UserID: 4,
			},
			url: "/profiles/1/",
			friendshipStatus: friendshipEntity.FriendshipStatus{
				RequestedToProfileID:   1,
				RequestedFromProfileID: 3,
				State:                  friendshipEntity.StatePending,
			},
			wantStatus: http.StatusOK,
			wantBody: []string{
				`You've already sent friendship request`,
			},
		},
		{
			name: "already_friends",
			profile: entity.Profile{
				ID:     1,
				UserID: 2,
			},
			self: entity.Profile{
				ID:     3,
				UserID: 4,
			},
			url: "/profiles/1/",
			friendshipStatus: friendshipEntity.FriendshipStatus{
				RequestedToProfileID:   1,
				RequestedFromProfileID: 3,
				State:                  friendshipEntity.StateAccepted,
			},
			wantStatus: http.StatusOK,
			wantBody: []string{
				`You are friends`,
				`<form method="POST" action="/friends/1/stop"><input type="submit" value="stop friendship"></form>`,
			},
		},
		{
			name: "request_declined",
			profile: entity.Profile{
				ID:     1,
				UserID: 2,
			},
			self: entity.Profile{
				ID:     3,
				UserID: 4,
			},
			url: "/profiles/1/",
			friendshipStatus: friendshipEntity.FriendshipStatus{
				RequestedToProfileID:   1,
				RequestedFromProfileID: 3,
				State:                  friendshipEntity.StatesDeclined,
			},
			wantStatus: http.StatusOK,
			wantBody: []string{
				`Your friendship request was declined`,
			},
		},
		{
			name: "request_declined_by_self",
			profile: entity.Profile{
				ID:     1,
				UserID: 2,
			},
			self: entity.Profile{
				ID:     3,
				UserID: 4,
			},
			url: "/profiles/1/",
			friendshipStatus: friendshipEntity.FriendshipStatus{
				RequestedToProfileID:   3,
				RequestedFromProfileID: 1,
				State:                  friendshipEntity.StatesDeclined,
			},
			wantStatus: http.StatusOK,
			wantBody: []string{
				`You have declined friendship request`,
				`<form method="POST" action="/friends/requests/1/accept/"><input type="submit" value="Accept"></form>`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			profilesUseCase := profilesMocks.NewMockUseCase(ctrl)
			profilesUseCase.EXPECT().GetByID(gomock.Any(), test.profile.ID).Return([]entity.Profile{test.profile}, nil)
			friendshipUseCase := friendshipMocks.NewMockUseCase(ctrl)
			friendshipUseCase.EXPECT().GetFriendshipStatus(gomock.Any(), test.profile.ID, test.self.ID).Return(test.friendshipStatus, nil)
			s := server.New(":80", "../../../templates", nil, profilesUseCase, friendshipUseCase)

			request := httptest.NewRequest("GET", test.url, nil)
			request = request.WithContext(middlewares.AddProfileToCtx(request.Context(), test.self))
			responseWriter := httptest.NewRecorder()
			s.Handle(responseWriter, request)

			response := responseWriter.Result()
			require.Equal(t, test.wantStatus, response.StatusCode)

			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			for _, bodyPart := range test.wantBody {
				require.Contains(t, string(body), bodyPart)
			}
		})
	}
}
