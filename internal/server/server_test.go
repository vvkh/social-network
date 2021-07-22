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
			s := New(":80", "../../templates", usersUseCase, profilesUseCase)

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
		{
			method:        "GET",
			route:         "/friends/",
			wantStatus:    http.StatusOK,
			wantToContain: "<title>Friends</title>",
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

			s := New(":80", "../../templates", usersUseCase, profilesUseCase)

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

			s := New(":80", "../../templates", usersUseCase, profilesUseCase)

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
