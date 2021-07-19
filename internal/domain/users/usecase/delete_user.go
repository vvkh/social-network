package usecase

import "context"

func (u *usecase) DeleteUser(ctx context.Context, userID uint64) error {
	return u.repo.DeleteUser(ctx, userID)
}
