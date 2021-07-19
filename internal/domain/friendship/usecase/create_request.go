package usecase

import (
	"context"
)

func (u *usecase) CreateRequest(ctx context.Context, profileIDFrom uint64, profileIDTo uint64) error {
	return u.repo.CreateRequest(ctx, profileIDFrom, profileIDTo)
}
