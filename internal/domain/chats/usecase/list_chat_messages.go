package usecase

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/chats/entity"
)

func (u *uc) ListChatMessages(ctx context.Context, profileID uint64, chatID uint64) (entity.Chat, []entity.Message, error) {
	return entity.Chat{}, nil, nil
}
