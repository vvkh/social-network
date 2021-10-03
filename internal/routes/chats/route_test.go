package chats_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	chatsEntity "github.com/vvkh/social-network/internal/domain/chats/entity"
	"github.com/vvkh/social-network/internal/domain/chats/mocks"
	friendShipMocks "github.com/vvkh/social-network/internal/domain/friendship/mocks"
	"github.com/vvkh/social-network/internal/domain/profiles/entity"
	"github.com/vvkh/social-network/internal/middlewares"
	"github.com/vvkh/social-network/internal/server"
)

func TestChatsPage(t *testing.T) {
	tests := []struct {
		name          string
		self          entity.Profile
		chatsResponse []chatsEntity.Chat
		wantBodyParts []string
	}{
		{
			name:          "page title is 'Chats'",
			wantBodyParts: []string{"<h1>Chats</h1>"},
		},
		{
			name: "chats are displayed",
			chatsResponse: []chatsEntity.Chat{
				{
					ID:                  1,
					Title:               "Some chat",
					UnreadMessagesCount: 1,
				},
				{
					ID:                  2,
					Title:               "Some chat #2",
					UnreadMessagesCount: 0,
				},
			},
			wantBodyParts: []string{
				`<a href="/chats/1/">Some chat (1)</a>`,
				`<a href="/chats/2/">Some chat #2</a>`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			chatsMock := mocks.NewMockUseCase(ctrl)
			chatsMock.EXPECT().ListChats(gomock.Any(), test.self.ID).Return(test.chatsResponse, nil)
			chatsMock.EXPECT().GetUnreadMessagesCount(gomock.Any(), test.self.ID).Return(int64(0), nil)
			friendshipMock := friendShipMocks.NewMockUseCase(ctrl)
			friendshipMock.EXPECT().ListPendingRequests(gomock.Any(), test.self.ID).Return(nil, nil)
			log, _ := zap.NewDevelopment()
			s := server.New(log.Sugar(), ":80", "../../../templates", nil, nil, friendshipMock, chatsMock)

			request := httptest.NewRequest("GET", "/chats/", nil)
			request = request.WithContext(middlewares.AddProfileToCtx(request.Context(), test.self))
			responseWriter := httptest.NewRecorder()

			s.Handle(responseWriter, request)

			response := responseWriter.Result()
			require.Equal(t, http.StatusOK, response.StatusCode)

			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			for _, wantPart := range test.wantBodyParts {
				require.Contains(t, string(body), wantPart)
			}
		})
	}
}
