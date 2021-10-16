package chats

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/chats/entity"
)

//go:generate mockgen -destination=mocks/usecase.go -package=mocks -source=usecase.go

type UseCase interface {
	GetUnreadMessagesCount(ctx context.Context, profileID uint64) (int64, error)
	GetOrCreateChat(ctx context.Context, oneProfileID uint64, otherProfileID int64) (int64, error)
	ListChats(ctx context.Context, profileID uint64) ([]entity.Chat, error)
	ListChatMessages(ctx context.Context, profileID uint64, chatID uint64) (entity.Chat, []entity.Message, error)
	SendMessage(ctx context.Context, authorProfileID uint64, chatID uint64, content string) error
}
