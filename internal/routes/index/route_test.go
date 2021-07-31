package index_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vvkh/social-network/internal/domain/profiles/entity"
	"github.com/vvkh/social-network/internal/middlewares"
	"github.com/vvkh/social-network/internal/server"
)

func TestHandle(t *testing.T) {
	tests := []struct {
		name        string
		profile     entity.Profile
		wantHeaders map[string]string
	}{
		{
			name: "index_redirects_to_profile_id",
			profile: entity.Profile{
				ID: 1,
			},
			wantHeaders: map[string]string{
				"Location": "/profiles/1/",
			},
		},
		{
			name: "index_redirects_to_profile_id_2",
			profile: entity.Profile{
				ID: 2,
			},
			wantHeaders: map[string]string{
				"Location": "/profiles/2/",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log := zap.NewNop().Sugar()
			s := server.New(log, ":80", "../../../templates", nil, nil, nil)
			request := httptest.NewRequest("GET", "/", nil)
			request = request.WithContext(middlewares.AddProfileToCtx(request.Context(), test.profile))
			responseWriter := httptest.NewRecorder()
			s.Handle(responseWriter, request)
			response := responseWriter.Result()
			require.Equal(t, http.StatusFound, response.StatusCode)

			for header, value := range test.wantHeaders {
				require.Equal(t, response.Header.Get(header), value)
			}
		})
	}
}
