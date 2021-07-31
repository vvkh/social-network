package profiles_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vvkh/social-network/internal/domain/profiles/entity"
	"github.com/vvkh/social-network/internal/domain/profiles/mocks"
	"github.com/vvkh/social-network/internal/middlewares"
	"github.com/vvkh/social-network/internal/server"
)

func TestHandle(t *testing.T) {

	tests := []struct {
		name         string
		mockResponse []entity.Profile
		wantBody     []string
	}{
		{
			name: "all_profiles_are_displayed",
			mockResponse: []entity.Profile{
				{
					ID:        2,
					FirstName: "John",
					LastName:  "Doe",
				},
				{
					ID:        3,
					FirstName: "Topsy",
					LastName:  "Cret",
				},
			},
			wantBody: []string{
				`<a href="/profiles/2/">John Doe</a>`,
				`<a href="/profiles/3/">Topsy Cret</a>`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log := zap.NewNop().Sugar()
			ctrl := gomock.NewController(t)
			profilesUseCase := mocks.NewMockUseCase(ctrl)
			profilesUseCase.EXPECT().ListProfiles(gomock.Any()).Return(test.mockResponse, nil)
			s := server.New(log, ":80", "../../../templates", nil, profilesUseCase, nil)
			request := httptest.NewRequest("GET", "/profiles/", nil)

			self := entity.Profile{
				ID:        1,
				UserID:    2,
				FirstName: "John",
				LastName:  "Doe",
			}
			request = request.WithContext(middlewares.AddProfileToCtx(request.Context(), self))
			responseWriter := httptest.NewRecorder()
			s.Handle(responseWriter, request)
			response := responseWriter.Result()
			require.Equal(t, http.StatusOK, response.StatusCode)
			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			for _, bodyPart := range test.wantBody {
				require.Contains(t, string(body), bodyPart)
			}
		})
	}
}
