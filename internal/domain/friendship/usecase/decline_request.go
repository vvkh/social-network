package usecase

import (
	"context"
)

func (u *usecase) DeclineRequest(ctx context.Context, profileIDFrom uint64, profileIDTo uint64) error {
	return u.repo.DeclineRequest(ctx, profileIDFrom, profileIDTo)
}
