package usecase

import "context"

func (u *uc) GetUnreadMessagesCount(ctx context.Context, profileID uint64) (int64, error) {
	return 0, nil
}
