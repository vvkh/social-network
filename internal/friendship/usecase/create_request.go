package usecase

import (
	"context"
)

func (u *usecase) CreateRequest(ctx context.Context, userIDFrom uint64, userIDTo uint64) error {
	return u.repo.CreateRequest(ctx, userIDFrom, userIDTo)
}
