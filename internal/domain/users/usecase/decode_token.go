package usecase

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/users/entity"
)

func (u *usecase) DecodeToken(ctx context.Context, token string) (entity.AccessToken, error) {
	return entity.Parse(token, u.jwtSecret)
}
