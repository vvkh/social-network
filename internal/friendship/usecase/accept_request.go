package usecase

import (
	"context"
)

func (u *usecase) AcceptRequest(ctx context.Context, userIDFrom uint64, userIDTo uint64) error {
	return u.repo.AcceptRequest(ctx, userIDFrom, userIDTo)
}
