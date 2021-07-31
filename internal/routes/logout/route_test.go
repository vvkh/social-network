package logout_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vvkh/social-network/internal/server"
)

func TestHandle(t *testing.T) {
	tests := []struct {
		name        string
		wantHeaders map[string]string
	}{
		{
			name: "logout_clears_cookie_and_redirects_to_index",
			wantHeaders: map[string]string{
				"Location":   "/",
				"Set-Cookie": "token=; Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT; HttpOnly",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			log := zap.NewNop().Sugar()
			s := server.New(log, ":80", "../../../templates", nil, nil, nil)
			request := httptest.NewRequest("GET", "/logout/", nil)
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
