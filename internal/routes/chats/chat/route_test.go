package chat_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	chatsEntity "github.com/vvkh/social-network/internal/domain/chats/entity"
	"github.com/vvkh/social-network/internal/domain/chats/mocks"
	friendShipMocks "github.com/vvkh/social-network/internal/domain/friendship/mocks"
	"github.com/vvkh/social-network/internal/domain/profiles/entity"
	profilesMocks "github.com/vvkh/social-network/internal/domain/profiles/mocks"
	"github.com/vvkh/social-network/internal/middlewares"
	"github.com/vvkh/social-network/internal/server"
)

func TestChatPage(t *testing.T) {
	tests := []struct {
		name                   string
		self                   entity.Profile
		chat                   chatsEntity.Chat
		chatMessages           []chatsEntity.Message
		wantGetProfilesRequest []interface{}
		getProfilesResponse    []entity.Profile
		wantBodyParts          []string
	}{
		{
			name: "page title is char title",
			chat: chatsEntity.Chat{
				Title: "My chat",
				ID:    1,
			},
			wantBodyParts: []string{"<h1>My chat</h1>"},
		},
		{
			name: "messages is displayed",
			chat: chatsEntity.Chat{
				Title: "My chat",
				ID:    1,
			},
			chatMessages: []chatsEntity.Message{
				{
					Content:         "Hello world!",
					AuthorProfileID: 1,
					SentAt:          time.Date(2021, 10, 3, 23, 58, 0, 0, time.Local),
				},
				{
					Content:         "Hi there!",
					AuthorProfileID: 2,
					SentAt:          time.Date(2021, 10, 3, 23, 59, 0, 0, time.Local),
				},
			},
			wantGetProfilesRequest: []interface{}{
				uint64(1),
				uint64(2),
			},
			getProfilesResponse: []entity.Profile{
				{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
				},
				{
					ID:        2,
					FirstName: "Topsy",
					LastName:  "Cret",
				},
			},
			wantBodyParts: []string{
				"2021 10 03 23:58 <b>John Doe</b>: Hello world!",
				"2021 10 03 23:59 <b>Topsy Cret</b>: Hi there!",
			},
		},
		{
			name: "contains submit message form",
			chat: chatsEntity.Chat{
				Title: "My chat",
				ID:    1,
			},
			chatMessages:           []chatsEntity.Message{},
			wantGetProfilesRequest: []interface{}{},
			wantBodyParts: []string{
				`<form method="POST" action="/chats/1/message/">`,
				`<textarea id="message" name="message"></textarea>`,
				`<input type="submit" value="send message">`,
				`</form>`,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			chatsMock := mocks.NewMockUseCase(ctrl)
			chatsMock.EXPECT().ListChatMessages(gomock.Any(), test.self.ID, test.chat.ID).Return(test.chat, test.chatMessages, nil)
			chatsMock.EXPECT().GetUnreadMessagesCount(gomock.Any(), test.self.ID).Return(int64(0), nil)
			friendshipMock := friendShipMocks.NewMockUseCase(ctrl)
			friendshipMock.EXPECT().ListPendingRequests(gomock.Any(), test.self.ID).Return(nil, nil)
			profilesMock := profilesMocks.NewMockUseCase(ctrl)
			if len(test.chatMessages) > 0 {
				profilesMock.EXPECT().GetByID(gomock.Any(), test.wantGetProfilesRequest...).Return(test.getProfilesResponse, nil)
			}
			log, _ := zap.NewDevelopment()
			s := server.New(log.Sugar(), ":80", "../../../../templates", nil, profilesMock, friendshipMock, chatsMock)

			request := httptest.NewRequest("GET", "/chats/1/", nil)
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
