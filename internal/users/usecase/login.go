package usecase

import (
	"context"
)

func (u *usecase) Login(ctx context.Context, username string, password string) (uint64, error) {
	return u.repo.Login(ctx, username, password)
}
