package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

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
			s := New(":80", "../../templates", usersUseCase)

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
		wantStatus    int
		wantHeaders   map[string]string
		wantToContain string
	}{
		{
			method:      "GET",
			route:       "/",
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
			route:         "/profiles/",
			wantStatus:    http.StatusOK,
			wantToContain: "<title>Profiles</title>",
		},
		{
			method:        "GET",
			route:         "/profiles/1/",
			wantStatus:    http.StatusOK,
			wantToContain: "<title>Profile</title>",
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
			usersUseCase.EXPECT().DecodeToken(gomock.Any(), gomock.Any()).Return(entity.AccessToken{UserID: 1, ProfileID: 2}, nil)

			s := New(":80", "../../templates", usersUseCase)

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
