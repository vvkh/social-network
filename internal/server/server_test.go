package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	chatMocks "github.com/vvkh/social-network/internal/domain/chats/mocks"
	friendshipMock "github.com/vvkh/social-network/internal/domain/friendship/mocks"
	profilesEntity "github.com/vvkh/social-network/internal/domain/profiles/entity"
	profilesMock "github.com/vvkh/social-network/internal/domain/profiles/mocks"
	"github.com/vvkh/social-network/internal/domain/users/entity"
	"github.com/vvkh/social-network/internal/domain/users/mocks"
	"github.com/vvkh/social-network/internal/middlewares"
)

var (
	authProtectedRoutes = []struct {
		method string
		url    string
	}{
		{
			method: "GET",
			url:    "/",
		},
		{
			method: "GET",
			url:    "/profiles/",
		},
		{
			method: "GET",
			url:    "/profiles/1/",
		},
		{
			method: "GET",
			url:    "/friends/",
		},
		{
			method: "GET",
			url:    "/friends/requests/",
		},
		{
			method: "POST",
			url:    "/friends/requests/1/create/",
		},
		{
			method: "POST",
			url:    "/friends/requests/1/accept/",
		},
		{
			method: "POST",
			url:    "/friends/requests/1/decline/",
		},
		{
			method: "GET",
			url:    "/chats/",
		},
	}

	routesWithNavbar = []struct {
		method string
		url    string
	}{
		{
			method: "GET",
			url:    "/profiles/",
		},
		{
			method: "GET",
			url:    "/profiles/1/",
		},
		{
			method: "GET",
			url:    "/friends/",
		},
		{
			method: "GET",
			url:    "/friends/requests/",
		},
		{
			method: "GET",
			url:    "/chats/",
		},
	}
)

func TestAuthProtectedRoutesRedirectToLogin(t *testing.T) {
	for _, route := range authProtectedRoutes {
		t.Run(route.method+" "+route.url, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			usersUseCase := mocks.NewMockUseCase(ctrl)
			profilesUseCase := profilesMock.NewMockUseCase(ctrl)
			friendshipUseCase := friendshipMock.NewMockUseCase(ctrl)
			log, _ := zap.NewDevelopment()
			s := New(log.Sugar(), ":80", "../../templates", usersUseCase, profilesUseCase, friendshipUseCase, nil)

			request := httptest.NewRequest(route.method, route.url, nil)
			responseWriter := httptest.NewRecorder()
			s.Handle(responseWriter, request)

			response := responseWriter.Result()
			require.Equal(t, http.StatusFound, response.StatusCode)
			require.Equal(t, "/login/", response.Header.Get("Location"))
		})
	}
}

func TestAuthProtectedRoutesWithInvalidTokenRedirectsToLogin(t *testing.T) {
	tests := []struct {
		name                   string
		mocksToken             entity.AccessToken
		getProfileMockResponse []profilesEntity.Profile
		getProfileMockErr      error
		wantResetCookie        bool
	}{
		{
			mocksToken: entity.AccessToken{
				UserID:    1,
				ProfileID: 2,
			},
			getProfileMockResponse: []profilesEntity.Profile{},
			wantResetCookie:        true,
		},
		{
			mocksToken: entity.AccessToken{
				UserID:    1,
				ProfileID: 2,
			},
			getProfileMockResponse: []profilesEntity.Profile{
				{
					ID:     2, // different user id
					UserID: 3,
				},
			},
			wantResetCookie: true,
		},
		{
			mocksToken: entity.AccessToken{
				UserID:    1,
				ProfileID: 2,
			},
			getProfileMockErr: errors.New("failed"),
			wantResetCookie:   false,
		},
	}
	for _, route := range authProtectedRoutes {
		for _, test := range tests {
			testName := fmt.Sprintf("%s %s %s", route.method, route.url, test.name)
			t.Run(testName, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				usersUseCase := mocks.NewMockUseCase(ctrl)
				usersUseCase.EXPECT().DecodeToken(gomock.Any(), gomock.Any()).Return(test.mocksToken, nil)
				profilesUseCase := profilesMock.NewMockUseCase(ctrl)
				profilesUseCase.EXPECT().GetByID(gomock.Any(), test.mocksToken.ProfileID).Return(test.getProfileMockResponse, test.getProfileMockErr)
				log, _ := zap.NewDevelopment()
				s := New(log.Sugar(), ":80", "../../templates", usersUseCase, profilesUseCase, nil, nil)

				request := httptest.NewRequest(route.method, route.url, nil)
				request.AddCookie(&http.Cookie{
					Name:     "token",
					Value:    "secret",
					Path:     "/",
					HttpOnly: true,
				})
				responseWriter := httptest.NewRecorder()
				s.Handle(responseWriter, request)

				response := responseWriter.Result()
				require.Equal(t, http.StatusFound, response.StatusCode)
				require.Equal(t, "/login/", response.Header.Get("Location"))
				if test.wantResetCookie {
					require.Equal(t, "token=; Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT; HttpOnly", response.Header.Get("Set-Cookie"))
				}
			})
		}
	}
}

func TestNavbar(t *testing.T) {
	tests := []struct {
		name                      string
		pendingFriendshipRequests []profilesEntity.Profile
		unreadMessagesCount       int64
		wantBody                  []string
	}{
		{
			name: "pending_requests_count_shown_in_the_navbar",
			pendingFriendshipRequests: []profilesEntity.Profile{
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
				`<a href="/friends/">Friends (2)</a>`,
			},
		},
		{
			name:                      "pending_requests_count_shown_in_the_navbar_if_no_requests",
			pendingFriendshipRequests: []profilesEntity.Profile{},
			wantBody: []string{
				`<a href="/friends/">Friends</a>`,
			},
		},
		{
			name:                "unread_messages_count_is_shown_in_navbar",
			unreadMessagesCount: 10,
			wantBody: []string{
				`<a href="/chats/">Chats (10)</a>`,
			},
		},
		{
			name:                "unread_messages_count_is_not_shown_in_navbar_if_no_messages",
			unreadMessagesCount: 0,
			wantBody: []string{
				`<a href="/chats/">Chats</a>`,
			},
		},
	}

	profile := profilesEntity.Profile{
		ID: 1,
	}

	for _, test := range tests {
		for _, route := range routesWithNavbar {
			t.Run(route.method+" "+route.url, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				usersUseCase := mocks.NewMockUseCase(ctrl)
				profilesUseCase := profilesMock.NewMockUseCase(ctrl)
				profilesUseCase.EXPECT().ListProfiles(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				profilesUseCase.EXPECT().GetByID(gomock.Any(), profile.ID).Return([]profilesEntity.Profile{profile}, nil).AnyTimes()
				friendshipUseCase := friendshipMock.NewMockUseCase(ctrl)
				friendshipUseCase.EXPECT().ListPendingRequests(gomock.Any(), profile.ID).Return(test.pendingFriendshipRequests, nil).AnyTimes()
				friendshipUseCase.EXPECT().GetFriendshipStatus(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				friendshipUseCase.EXPECT().ListFriends(gomock.Any(), gomock.Any()).AnyTimes()
				chatsUseCase := chatMocks.NewMockUseCase(ctrl)
				chatsUseCase.EXPECT().GetUnreadMessagesCount(gomock.Any(), profile.ID).Return(test.unreadMessagesCount, nil)
				chatsUseCase.EXPECT().ListChats(gomock.Any(), profile.ID).AnyTimes()
				log, _ := zap.NewDevelopment()
				s := New(log.Sugar(), ":80", "../../templates", usersUseCase, profilesUseCase, friendshipUseCase, chatsUseCase)

				request := httptest.NewRequest(route.method, route.url, nil)
				request = request.WithContext(middlewares.AddProfileToCtx(request.Context(), profile))
				responseWriter := httptest.NewRecorder()
				s.Handle(responseWriter, request)

				response := responseWriter.Result()
				body, err := io.ReadAll(response.Body)
				require.NoError(t, err)

				for _, bodyPart := range test.wantBody {
					require.Contains(t, string(body), bodyPart)
				}
			})
		}
	}
}
