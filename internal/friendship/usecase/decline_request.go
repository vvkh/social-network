package usecase

import (
	"context"
)

func (u *usecase) DeclineRequest(ctx context.Context, userIDFrom uint64, userIDTo uint64) error {
	return u.repo.DeclineRequest(ctx, userIDFrom, userIDTo)
}
