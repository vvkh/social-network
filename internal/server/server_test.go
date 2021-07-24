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

	friendshipMock "github.com/vvkh/social-network/internal/domain/friendship/mocks"
	profilesEntity "github.com/vvkh/social-network/internal/domain/profiles/entity"
	profilesMock "github.com/vvkh/social-network/internal/domain/profiles/mocks"
	"github.com/vvkh/social-network/internal/domain/users/entity"
	"github.com/vvkh/social-network/internal/domain/users/mocks"
)

func TestRoutesSmoke(t *testing.T) {
	tests := []struct {
		method        string
		route         string
		wantStatus    int
		wantHeaders   map[string]string
		wantToContain string
	}{
		{
			method:      "GET",
			route:       "/",
			wantStatus:  http.StatusFound,
			wantHeaders: map[string]string{"Location": "/login/"},
		},
		{
			method:        "GET",
			route:         "/login/",
			wantStatus:    http.StatusOK,
			wantToContain: "<title>Login</title>",
		},
		{
			method:        "GET",
			route:         "/register/",
			wantStatus:    http.StatusOK,
			wantToContain: "<title>Register</title>",
		},
		{
			method:      "GET",
			route:       "/profiles/",
			wantStatus:  http.StatusFound,
			wantHeaders: map[string]string{"Location": "/login/"},
		},
		{
			method:      "GET",
			route:       "/profiles/1/",
			wantStatus:  http.StatusFound,
			wantHeaders: map[string]string{"Location": "/login/"},
		},
		{
			method:      "GET",
			route:       "/friends/",
			wantStatus:  http.StatusFound,
			wantHeaders: map[string]string{"Location": "/login/"},
		},
	}
	for _, test := range tests {
		t.Run(test.route, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			usersUseCase := mocks.NewMockUseCase(ctrl)
			profilesUseCase := profilesMock.NewMockUseCase(ctrl)
			friendshipUseCase := friendshipMock.NewMockUseCase(ctrl)
			s := New(":80", "../../templates", usersUseCase, profilesUseCase, friendshipUseCase)

			request := httptest.NewRequest(test.method, test.route, nil)
			responseWriter := httptest.NewRecorder()
			s.Handle(responseWriter, request)

			response := responseWriter.Result()
			require.Equal(t, test.wantStatus, response.StatusCode)

			for k, v := range response.Header {
				fmt.Println(k, v)
			}
			for header, value := range test.wantHeaders {
				require.Equal(t, value, response.Header.Get(header))
			}

			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			require.Contains(t, string(body), test.wantToContain)
		})
	}
}

func TestRoutesSmokeWithAuthentication(t *testing.T) {
	tests := []struct {
		method        string
		route         string
		profileID     uint64
		userID        uint64
		wantStatus    int
		wantHeaders   map[string]string
		wantToContain string
	}{
		{
			method:      "GET",
			route:       "/",
			profileID:   2,
			wantStatus:  http.StatusFound,
			wantHeaders: map[string]string{"Location": "/profiles/2/"},
		},
		{
			method:        "GET",
			route:         "/login/",
			wantStatus:    http.StatusOK,
			wantToContain: "<title>Login</title>",
		},
		{
			method:        "GET",
			route:         "/register/",
			wantStatus:    http.StatusOK,
			wantToContain: "<title>Register</title>",
		},
		{
			method:        "GET",
			profileID:     2,
			route:         "/profiles/",
			wantStatus:    http.StatusOK,
			wantToContain: `<a href="/profiles/2/">John Snow</a>`,
		},
		{
			method:        "GET",
			profileID:     2,
			route:         "/profiles/2/",
			wantStatus:    http.StatusOK,
			wantToContain: "<h1>John Snow</h1>",
		},
	}
	for _, test := range tests {
		t.Run(test.route, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			usersUseCase := mocks.NewMockUseCase(ctrl)
			usersUseCase.EXPECT().DecodeToken(gomock.Any(), gomock.Any()).Return(entity.AccessToken{UserID: test.userID, ProfileID: test.profileID}, nil)

			sampleProfile := profilesEntity.Profile{
				ID:        test.profileID,
				UserID:    test.userID,
				FirstName: "John",
				LastName:  "Snow",
			}
			profilesUseCase := profilesMock.NewMockUseCase(ctrl)
			profilesUseCase.EXPECT().GetByID(gomock.Any(), test.profileID).Return([]profilesEntity.Profile{sampleProfile}, nil).AnyTimes()
			profilesUseCase.EXPECT().GetByUserID(gomock.Any(), test.userID).Return([]profilesEntity.Profile{sampleProfile}, nil).AnyTimes()
			profilesUseCase.EXPECT().ListProfiles(gomock.Any()).Return([]profilesEntity.Profile{sampleProfile}, nil).AnyTimes()
			friendshipUseCase := friendshipMock.NewMockUseCase(ctrl)
			s := New(":80", "../../templates", usersUseCase, profilesUseCase, friendshipUseCase)

			request := httptest.NewRequest(test.method, test.route, nil)
			request.AddCookie(&http.Cookie{
				Name:     "token",
				Value:    "secret",
				Path:     "/",
				HttpOnly: true,
			})
			responseWriter := httptest.NewRecorder()
			s.Handle(responseWriter, request)

			response := responseWriter.Result()
			require.Equal(t, test.wantStatus, response.StatusCode)

			for k, v := range response.Header {
				fmt.Println(k, v)
			}
			for header, value := range test.wantHeaders {
				require.Equal(t, value, response.Header.Get(header))
			}

			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			require.Contains(t, string(body), test.wantToContain)
		})
	}
}

func TestRoutesSmokeWithInvalidAuthenticationToken(t *testing.T) {
	tests := []struct {
		method     string
		route      string
		mocksToken entity.AccessToken

		getProfileMockResponse []profilesEntity.Profile
		getProfileMockErr      error

		wantStatus    int
		wantHeaders   map[string]string
		wantToContain string
	}{
		{
			method: "GET",
			route:  "/",
			mocksToken: entity.AccessToken{
				UserID:    1,
				ProfileID: 2,
			},
			getProfileMockResponse: []profilesEntity.Profile{},
			wantStatus:             http.StatusFound,
			wantHeaders: map[string]string{
				"Location":   "/login/",
				"Set-Cookie": "token=; Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT; HttpOnly",
			},
		},
		{
			method: "GET",
			route:  "/",
			mocksToken: entity.AccessToken{
				UserID:    1,
				ProfileID: 2,
			},
			getProfileMockResponse: []profilesEntity.Profile{
				{
					ID:     3, // different profile id
					UserID: 1,
				},
			},
			wantStatus: http.StatusFound,
			wantHeaders: map[string]string{
				"Location":   "/login/",
				"Set-Cookie": "token=; Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT; HttpOnly",
			},
		},
		{
			method: "GET",
			route:  "/",
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
			wantStatus: http.StatusFound,
			wantHeaders: map[string]string{
				"Location":   "/login/",
				"Set-Cookie": "token=; Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT; HttpOnly",
			},
		},
		{
			method: "GET",
			route:  "/",
			mocksToken: entity.AccessToken{
				UserID:    1,
				ProfileID: 2,
			},
			getProfileMockErr: errors.New("failed"),
			wantStatus:        http.StatusFound,
			wantHeaders: map[string]string{
				"Location":   "/login/",
				"Set-Cookie": "", // not set as we don't expire token if there was an error while getting profiles
			},
		},
	}
	for _, test := range tests {
		t.Run(test.route, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			usersUseCase := mocks.NewMockUseCase(ctrl)
			usersUseCase.EXPECT().DecodeToken(gomock.Any(), gomock.Any()).Return(test.mocksToken, nil)
			profilesUseCase := profilesMock.NewMockUseCase(ctrl)
			profilesUseCase.EXPECT().GetByUserID(gomock.Any(), test.mocksToken.UserID).Return(test.getProfileMockResponse, test.getProfileMockErr)
			friendshipUseCase := friendshipMock.NewMockUseCase(ctrl)
			s := New(":80", "../../templates", usersUseCase, profilesUseCase, friendshipUseCase)

			request := httptest.NewRequest(test.method, test.route, nil)
			request.AddCookie(&http.Cookie{
				Name:     "token",
				Value:    "secret",
				Path:     "/",
				HttpOnly: true,
			})
			responseWriter := httptest.NewRecorder()
			s.Handle(responseWriter, request)

			response := responseWriter.Result()
			require.Equal(t, test.wantStatus, response.StatusCode)

			for k, v := range response.Header {
				fmt.Println(k, v)
			}
			for header, value := range test.wantHeaders {
				require.Equal(t, value, response.Header.Get(header))
			}

			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			require.Contains(t, string(body), test.wantToContain)
		})
	}
}

func TestProfilePage(t *testing.T) {
	tests := []struct {
		name       string
		profile    profilesEntity.Profile
		self       profilesEntity.Profile
		url        string
		wantStatus int
		wantBody   string
	}{
		{
			name: "profile_page_contains_friendship_request",
			profile: profilesEntity.Profile{
				ID:     1,
				UserID: 2,
			},
			self: profilesEntity.Profile{
				ID:     3,
				UserID: 4,
			},
			url:        "/profiles/1/",
			wantStatus: http.StatusOK,
			wantBody:   `<input type="submit" value="request friendship">`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			usersUseCase := mocks.NewMockUseCase(ctrl)
			usersUseCase.EXPECT().DecodeToken(gomock.Any(), gomock.Any()).Return(entity.AccessToken{UserID: test.self.UserID, ProfileID: test.self.ID}, nil)
			profilesUseCase := profilesMock.NewMockUseCase(ctrl)
			profilesUseCase.EXPECT().GetByID(gomock.Any(), test.profile.ID).Return([]profilesEntity.Profile{test.profile}, nil).AnyTimes()
			profilesUseCase.EXPECT().GetByUserID(gomock.Any(), test.self.UserID).Return([]profilesEntity.Profile{test.self}, nil).AnyTimes()
			friendshipUseCase := friendshipMock.NewMockUseCase(ctrl)

			s := New(":80", "../../templates", usersUseCase, profilesUseCase, friendshipUseCase)

			request := httptest.NewRequest("GET", test.url, nil)
			request.AddCookie(&http.Cookie{
				Name:     "token",
				Value:    "secret",
				Path:     "/",
				HttpOnly: true,
			})
			responseWriter := httptest.NewRecorder()
			s.Handle(responseWriter, request)

			response := responseWriter.Result()
			require.Equal(t, test.wantStatus, response.StatusCode)

			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			require.Contains(t, string(body), test.wantBody)
		})
	}
}

func TestProfilePageFriendshipRequest(t *testing.T) {
	tests := []struct {
		name       string
		profile    profilesEntity.Profile
		self       profilesEntity.Profile
		url        string
		wantStatus int
	}{
		{
			name: "submit_friendship_form",
			profile: profilesEntity.Profile{
				ID:     1,
				UserID: 2,
			},
			self: profilesEntity.Profile{
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
			usersUseCase := mocks.NewMockUseCase(ctrl)
			usersUseCase.EXPECT().DecodeToken(gomock.Any(), gomock.Any()).Return(entity.AccessToken{UserID: test.self.UserID, ProfileID: test.self.ID}, nil)
			profilesUseCase := profilesMock.NewMockUseCase(ctrl)
			profilesUseCase.EXPECT().GetByID(gomock.Any(), test.profile.ID).Return([]profilesEntity.Profile{test.profile}, nil).AnyTimes()
			profilesUseCase.EXPECT().GetByUserID(gomock.Any(), test.self.UserID).Return([]profilesEntity.Profile{test.self}, nil).AnyTimes()
			friendshipUseCase := friendshipMock.NewMockUseCase(ctrl)
			friendshipUseCase.EXPECT().CreateRequest(gomock.Any(), test.self.ID, test.profile.ID).Return(nil)

			s := New(":80", "../../templates", usersUseCase, profilesUseCase, friendshipUseCase)

			request := httptest.NewRequest("POST", test.url, nil)
			request.AddCookie(&http.Cookie{
				Name:     "token",
				Value:    "secret",
				Path:     "/",
				HttpOnly: true,
			})
			responseWriter := httptest.NewRecorder()
			s.Handle(responseWriter, request)

			response := responseWriter.Result()
			require.Equal(t, test.wantStatus, response.StatusCode)
		})
	}
}

func TestFriendsPage(t *testing.T) {
	tests := []struct {
		name               string
		self               profilesEntity.Profile
		friendshipRequests []profilesEntity.Profile
		friends            []profilesEntity.Profile
		wantBody           []string
	}{
		{
			name: "friends_are_shown_on_page",
			self: profilesEntity.Profile{
				ID:     1,
				UserID: 2,
			},
			friends: []profilesEntity.Profile{
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
				"<h1>Friends</h1>",
				`<a href="/profiles/3/">John Doe</a>`,
				`<a href="/profiles/5/">Topsy Cret</a>`,
			},
		},
		{
			name: "pending_requests_count_shown_on_page",
			self: profilesEntity.Profile{
				ID:     1,
				UserID: 2,
			},
			friendshipRequests: []profilesEntity.Profile{
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
				"<h1>Friends</h1>",
				`<a href="/friends/requests/">Pending friendship requests (2)</a>`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			usersUseCase := mocks.NewMockUseCase(ctrl)
			usersUseCase.EXPECT().DecodeToken(gomock.Any(), gomock.Any()).Return(entity.AccessToken{UserID: test.self.UserID, ProfileID: test.self.ID}, nil)
			profilesUseCase := profilesMock.NewMockUseCase(ctrl)
			profilesUseCase.EXPECT().GetByUserID(gomock.Any(), test.self.UserID).Return([]profilesEntity.Profile{test.self}, nil).AnyTimes()
			friendshipUseCase := friendshipMock.NewMockUseCase(ctrl)
			friendshipUseCase.EXPECT().ListFriends(gomock.Any(), test.self.ID).Return(test.friends, nil)
			friendshipUseCase.EXPECT().ListPendingRequests(gomock.Any(), test.self.ID).Return(test.friendshipRequests, nil)

			s := New(":80", "../../templates", usersUseCase, profilesUseCase, friendshipUseCase)

			request := httptest.NewRequest("GET", "/friends/", nil)
			request.AddCookie(&http.Cookie{
				Name:     "token",
				Value:    "secret",
				Path:     "/",
				HttpOnly: true,
			})
			responseWriter := httptest.NewRecorder()
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

func TestFriendshipRequestPage(t *testing.T) {
	tests := []struct {
		name               string
		self               profilesEntity.Profile
		friendshipRequests []profilesEntity.Profile
		wantBody           []string
	}{
		{
			name: "friendship_requests_are_shown_on_page",
			self: profilesEntity.Profile{
				ID:     1,
				UserID: 2,
			},
			friendshipRequests: []profilesEntity.Profile{
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
			usersUseCase := mocks.NewMockUseCase(ctrl)
			usersUseCase.EXPECT().DecodeToken(gomock.Any(), gomock.Any()).Return(entity.AccessToken{UserID: test.self.UserID, ProfileID: test.self.ID}, nil)
			profilesUseCase := profilesMock.NewMockUseCase(ctrl)
			profilesUseCase.EXPECT().GetByUserID(gomock.Any(), test.self.UserID).Return([]profilesEntity.Profile{test.self}, nil).AnyTimes()
			friendshipUseCase := friendshipMock.NewMockUseCase(ctrl)
			friendshipUseCase.EXPECT().ListPendingRequests(gomock.Any(), test.self.ID).Return(test.friendshipRequests, nil)

			s := New(":80", "../../templates", usersUseCase, profilesUseCase, friendshipUseCase)

			request := httptest.NewRequest("GET", "/friends/requests/", nil)
			request.AddCookie(&http.Cookie{
				Name:     "token",
				Value:    "secret",
				Path:     "/",
				HttpOnly: true,
			})
			responseWriter := httptest.NewRecorder()
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

func TestAcceptDeclineFriendshipRequest(t *testing.T) {
	tests := []struct {
		name          string
		self          profilesEntity.Profile
		url           string
		expectAccept  bool
		expectDecline bool
		profileID     uint64
	}{
		{
			name: "accept_request",
			self: profilesEntity.Profile{
				ID:     1,
				UserID: 2,
			},
			profileID:    3,
			expectAccept: true,
			url:          "/friends/requests/3/accept/",
		},
		{
			name: "decline_request",
			self: profilesEntity.Profile{
				ID:     1,
				UserID: 2,
			},
			profileID:     3,
			expectDecline: true,
			url:           "/friends/requests/3/decline/",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			usersUseCase := mocks.NewMockUseCase(ctrl)
			usersUseCase.EXPECT().DecodeToken(gomock.Any(), gomock.Any()).Return(entity.AccessToken{UserID: test.self.UserID, ProfileID: test.self.ID}, nil)
			profilesUseCase := profilesMock.NewMockUseCase(ctrl)
			profilesUseCase.EXPECT().GetByUserID(gomock.Any(), test.self.UserID).Return([]profilesEntity.Profile{test.self}, nil).AnyTimes()
			friendshipUseCase := friendshipMock.NewMockUseCase(ctrl)
			if test.expectAccept {
				friendshipUseCase.EXPECT().AcceptRequest(gomock.Any(), test.profileID, test.self.ID).Return(nil)
			}
			if test.expectDecline {
				friendshipUseCase.EXPECT().DeclineRequest(gomock.Any(), test.profileID, test.self.ID).Return(nil)
			}

			s := New(":80", "../../templates", usersUseCase, profilesUseCase, friendshipUseCase)

			request := httptest.NewRequest("POST", test.url, nil)
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
			require.Equal(t, "/friends/requests/", response.Header.Get("Location"))
		})
	}
}
