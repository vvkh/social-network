package chats

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/chats/entity"
)

//go:generate mockgen -destination=mocks/usecase.go -package=mocks -source=usecase.go

type UseCase interface {
	GetUnreadMessagesCount(ctx context.Context, profileID int64) (int64, error)
	GetOrCreateChat(ctx context.Context, oneProfileID int64, otherProfileID int64) (int64, error)
	ListChatMessages(ctx context.Context, chatID int64) ([]entity.Message, error)
	SendMessage(ctx context.Context, chatID int64, authorProfileID int64, message string) error
}
