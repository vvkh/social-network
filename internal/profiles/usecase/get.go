package usecase

import (
	"fmt"

	"github.com/vvkh/social-network/internal/profiles/entity"
)

func (u *usecase) GetByID(id uint64) (entity.Profile, error) {
	profile, err := u.repository.GetByID(id)
	if err != nil {
		return entity.Profile{}, fmt.Errorf("repository.GetByID error: %w", err)
	}
	return profile, nil
}
