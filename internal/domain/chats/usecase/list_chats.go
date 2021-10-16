package usecase

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/chats/entity"
)

func (u *uc) ListChats(ctx context.Context, profileID uint64) ([]entity.Chat, error) {
	return u.repo.ListChats(ctx, profileID)
}
