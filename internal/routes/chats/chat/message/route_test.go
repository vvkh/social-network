package message_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vvkh/social-network/internal/domain/chats/mocks"
	"github.com/vvkh/social-network/internal/domain/profiles/entity"
	"github.com/vvkh/social-network/internal/middlewares"
	"github.com/vvkh/social-network/internal/server"
)

func TestSendMessage(t *testing.T) {
	tests := []struct {
		name            string
		url             string
		self            entity.Profile
		form            string
		mockWantChatID  uint64
		mockWantMessage string
	}{
		{
			name: "send message",
			url:  "/chats/1/message/",
			self: entity.Profile{
				ID: 1,
			},
			form:            "message=123",
			mockWantChatID:  1,
			mockWantMessage: "123",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			chatsMock := mocks.NewMockUseCase(ctrl)
			chatsMock.EXPECT().SendMessage(gomock.Any(), test.self.ID, test.mockWantChatID, test.mockWantMessage).Return(nil)

			log, _ := zap.NewDevelopment()
			s := server.New(log.Sugar(), ":80", "../../../../../templates", nil, nil, nil, chatsMock)

			form := strings.NewReader(test.form)
			request := httptest.NewRequest("POST", test.url, form)
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			request = request.WithContext(middlewares.AddProfileToCtx(request.Context(), test.self))
			responseWriter := httptest.NewRecorder()

			s.Handle(responseWriter, request)

			response := responseWriter.Result()
			require.Equal(t, http.StatusFound, response.StatusCode)
		})
	}
}
