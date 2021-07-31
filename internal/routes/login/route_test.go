package login_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vvkh/social-network/internal/domain/users"
	"github.com/vvkh/social-network/internal/domain/users/mocks"
	"github.com/vvkh/social-network/internal/server"
)

func TestHandleGet(t *testing.T) {
	log := zap.NewNop().Sugar()
	ctrl := gomock.NewController(t)
	users := mocks.NewMockUseCase(ctrl)
	s := server.New(log, ":80", "../../../templates", users, nil, nil)
	request := httptest.NewRequest("GET", "/login/", nil)
	responseWriter := httptest.NewRecorder()
	s.Handle(responseWriter, request)
	response := responseWriter.Result()
	require.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	page := string(body)
	require.Contains(t, page, `<input id="username" name="username" type="text" placeholder="John Doe" required />`)
	require.Contains(t, page, `<input id="password" name="password" type="password" required />`)
	require.Contains(t, page, `<input type="submit" value="login" />`)
}

func TestHandlePost(t *testing.T) {
	type loginArgs struct {
		username string
		password string
	}

	tests := []struct {
		name              string
		form              string
		mockWantIn        *loginArgs
		mockResponse      string
		mockResponseError error
		wantStatus        int
		wantHeaders       map[string]string
		wantBody          string
	}{
		{
			name:              "error_shown_if_login_failed_dut_to_empty_credentials",
			mockWantIn:        &loginArgs{},
			mockResponseError: users.EmptyCredentials,
			wantStatus:        http.StatusBadRequest,
			wantBody:          "both password and username are required",
		},
		{
			name: "error_shown_if_login_failed_dut_to_bad_credentials",
			form: "username=foo&password=bar",
			mockWantIn: &loginArgs{
				username: "foo",
				password: "bar",
			},
			mockResponseError: users.AuthenticationFailed,
			wantStatus:        http.StatusForbidden,
			wantBody:          "authentication failed",
		},
		{
			name: "error_shown_if_login_failed_due_to_other_error",
			form: "username=foo&password=bar",
			mockWantIn: &loginArgs{
				username: "foo",
				password: "bar",
			},
			mockResponseError: errors.New("some error"),
			wantStatus:        http.StatusInternalServerError,
			wantBody:          "server failed",
		},
		{
			name: "redirects_to_index_if_successful",
			form: "username=foo&password=bar",
			mockWantIn: &loginArgs{
				username: "foo",
				password: "bar",
			},
			mockResponse: "token",
			wantStatus:   http.StatusFound,
			wantHeaders: map[string]string{
				"Location":   "/",
				"Set-Cookie": "token=token; Path=/; HttpOnly",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log := zap.NewNop().Sugar()
			ctrl := gomock.NewController(t)
			usersUC := mocks.NewMockUseCase(ctrl)
			usersUC.EXPECT().Login(gomock.Any(), test.mockWantIn.username, test.mockWantIn.password).Return(test.mockResponse, test.mockResponseError)

			s := server.New(log, ":80", "../../../templates", usersUC, nil, nil)
			form := strings.NewReader(test.form)
			request := httptest.NewRequest("POST", "/login/", form)
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			responseWriter := httptest.NewRecorder()
			s.Handle(responseWriter, request)
			response := responseWriter.Result()
			require.Equal(t, test.wantStatus, response.StatusCode)

			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			page := string(body)
			require.Contains(t, page, test.wantBody)
			for header, value := range test.wantHeaders {
				require.Equal(t, value, response.Header.Get(header))
			}
		})
	}
}
