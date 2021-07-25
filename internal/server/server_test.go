package server

import (
	"errors"
	"fmt"
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
	}
)

func TestAuthProtectedRoutesRedirectToLogin(t *testing.T) {
	for _, route := range authProtectedRoutes {
		t.Run(route.method+" "+route.url, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			usersUseCase := mocks.NewMockUseCase(ctrl)
			profilesUseCase := profilesMock.NewMockUseCase(ctrl)
			friendshipUseCase := friendshipMock.NewMockUseCase(ctrl)
			s := New(":80", "../../templates", usersUseCase, profilesUseCase, friendshipUseCase)

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
					ID:     3, // different profile id
					UserID: 1,
				},
			},
			wantResetCookie: true,
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
				profilesUseCase.EXPECT().GetByUserID(gomock.Any(), test.mocksToken.UserID).Return(test.getProfileMockResponse, test.getProfileMockErr)
				s := New(":80", "../../templates", usersUseCase, profilesUseCase, nil)

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
