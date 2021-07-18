package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/vvkh/social-network/internal/users/mocks"
)

func TestRoutesSmoke(t *testing.T) {
	tests := []struct {
		method        string
		route         string
		wantToContain string
	}{
		{
			method:        "GET",
			route:         "/",
			wantToContain: "hello world",
		},
		{
			method:        "GET",
			route:         "/login/",
			wantToContain: "<title>Login</title>",
		},
		{
			method:        "GET",
			route:         "/register/",
			wantToContain: "<title>Register</title>",
		},
		{
			method:        "GET",
			route:         "/profiles/",
			wantToContain: "<title>Profiles</title>",
		},
		{
			method:        "GET",
			route:         "/profiles/1/",
			wantToContain: "<title>Profile</title>",
		},
		{
			method:        "GET",
			route:         "/friends/",
			wantToContain: "<title>Friends</title>",
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
			if response.StatusCode != http.StatusOK {
				t.Errorf("GET %s status must be 200, got %d", test.route, response.StatusCode)
			}

			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Error("error while parsing response body")
			}

			if !strings.Contains(string(body), test.wantToContain) {
				t.Errorf(`want body to contain "%s", got = "%s"`, test.wantToContain, body)
			}
		})
	}
}
