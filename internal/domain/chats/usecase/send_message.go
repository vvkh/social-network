package usecase

import "context"

func (u *uc) SendMessage(ctx context.Context, authorProfileID uint64, chatID uint64, content string) error {
	return nil
}
