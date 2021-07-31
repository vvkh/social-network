package login_test

import (
	"io"
	"net/http"
	"net/http/httptest"
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
	require.Contains(t, page, `<input id="username" name="username" type="text" placeholder="John Doe" />`)
	require.Contains(t, page, `<input id="password" name="password" type="password" />`)
	require.Contains(t, page, `<input id="password" name="password" type="password" />`)
}
