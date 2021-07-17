package usecase

import (
	"context"
)

func (u *usecase) AcceptRequest(ctx context.Context, profileIDFrom uint64, profileIDTo uint64) error {
	return u.repo.AcceptRequest(ctx, profileIDFrom, profileIDTo)
}
