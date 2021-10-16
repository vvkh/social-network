package chats

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/chats/entity"
)

type Repository interface {
	GetOrCreateChat(ctx context.Context, chatTitle string, oneProfileID uint64, otherProfileID uint64) (uint64, error)
	ListChats(ctx context.Context, profileID uint64) ([]entity.Chat, error)
}
