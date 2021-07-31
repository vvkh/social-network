package register_test

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

	"github.com/vvkh/social-network/internal/domain/users/mocks"
	"github.com/vvkh/social-network/internal/server"
)

func TestHandleGet(t *testing.T) {
	log := zap.NewNop().Sugar()
	ctrl := gomock.NewController(t)
	usersUseCase := mocks.NewMockUseCase(ctrl)
	s := server.New(log, ":80", "../../../templates", usersUseCase, nil, nil)
	request := httptest.NewRequest("GET", "/register/", nil)
	responseWriter := httptest.NewRecorder()
	s.Handle(responseWriter, request)
	response := responseWriter.Result()
	require.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	page := string(body)
	require.Contains(t, page, `<input id="username" name="username" type="text" placeholder="johndoe" required/>`)
	require.Contains(t, page, `<input id="password" name="password" type="password" required/>`)
}

func TestHandlePost(t *testing.T) {
	type registerArgs struct {
		username  string
		password  string
		firstName string
		lastName  string
		age       uint8
		sex       string
		location  string
		about     string
	}

	tests := []struct {
		name              string
		form              string
		mockWantIn        *registerArgs
		mockResponseError error
		wantStatus        int
		wantHeaders       map[string]string
		wantBody          string
	}{
		//{
		//	name:              "error_shown_if_register_failed_dut_to_empty_credentials",
		//	mockWantIn:        &registerArgs{},
		//	mockResponseError: users.EmptyCredentials,
		//	wantStatus:        http.StatusBadRequest,
		//	wantBody:          "both password and username are required",
		//},
		//{
		//	name: "error_shown_if_login_failed_dut_to_bad_credentials",
		//	form: "username=foo&password=bar",
		//	mockWantIn: &registerArgs{
		//		username: "foo",
		//		password: "bar",
		//	},
		//	mockResponseError: users.AuthenticationFailed,
		//	wantStatus:        http.StatusForbidden,
		//	wantBody:          "authentication failed",
		//},
		//{
		//	name: "error_shown_if_login_failed_due_to_other_error",
		//	form: "username=foo&password=bar",
		//	mockWantIn: &registerArgs{
		//		username: "foo",
		//		password: "bar",
		//	},
		//	mockResponseError: errors.New("some error"),
		//	wantStatus:        http.StatusInternalServerError,
		//	wantBody:          "server failed",
		//},
		{
			name: "redirects_to_index_if_successful",
			form: "username=foo&password=bar&first_name=John&last_name=Doe&age=30&location=USA&sex=male&about=something",
			mockWantIn: &registerArgs{
				username:  "foo",
				password:  "bar",
				firstName: "John",
				lastName:  "Doe",
				age:       30,
				location:  "USA",
				sex:       "male",
				about:     "something",
			},
			wantStatus: http.StatusFound,
			wantHeaders: map[string]string{
				"Location": "/login/",
			},
		},
		{
			name:       "shows_error_if_bad_data",
			form:       "username=foo&password=bar&first_name=John&last_name=Doe&age=foo&location=USA&sex=male&about=something",
			wantStatus: http.StatusBadRequest,
			wantBody:   "age must be a positive number",
		},
		{
			name: "shows_error_if_usecase_failed",
			form: "username=foo&password=bar&first_name=John&last_name=Doe&age=30&location=USA&sex=male&about=something",
			mockWantIn: &registerArgs{
				username:  "foo",
				password:  "bar",
				firstName: "John",
				lastName:  "Doe",
				age:       30,
				location:  "USA",
				sex:       "male",
				about:     "something",
			},
			mockResponseError: errors.New("some error"),
			wantStatus:        http.StatusInternalServerError,
			wantBody:          "user creation failed",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log := zap.NewNop().Sugar()
			ctrl := gomock.NewController(t)
			usersUC := mocks.NewMockUseCase(ctrl)
			if test.mockWantIn != nil {
				usersUC.EXPECT().CreateUser(
					gomock.Any(),
					test.mockWantIn.username,
					test.mockWantIn.password,
					test.mockWantIn.firstName,
					test.mockWantIn.lastName,
					test.mockWantIn.age,
					test.mockWantIn.location,
					test.mockWantIn.sex,
					test.mockWantIn.about,
				).Return(uint64(0), uint64(0), test.mockResponseError)
			}
			s := server.New(log, ":80", "../../../templates", usersUC, nil, nil)
			form := strings.NewReader(test.form)
			request := httptest.NewRequest("POST", "/register/", form)
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
