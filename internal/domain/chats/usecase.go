package chats

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/chats/entity"
)

//go:generate mockgen -destination=mocks/usecase.go -package=mocks -source=usecase.go

type UseCase interface {
	GetUnreadMessagesCount(ctx context.Context, profileID uint64) (int64, error)
	GetOrCreateChat(ctx context.Context, oneProfileID uint64, otherProfileID int64) (int64, error)
	ListChatMessages(ctx context.Context, chatID uint64) ([]entity.Message, error)
	SendMessage(ctx context.Context, chatID uint64, authorProfileID int64, message string) error
}
