package repository

import (
	"context"

	"github.com/vvkh/social-network/internal/profiles/repository/dto"

	"github.com/vvkh/social-network/internal/profiles/entity"
)

func (r *repo) GetByID(ctx context.Context, id uint64) (entity.Profile, error) {
	query := `SELECT * FROM profiles WHERE id = ?`

	var profileDto dto.Profile
	err := r.db.GetContext(ctx, &profileDto, query, id)
	if err != nil {
		return entity.Profile{}, nil
	}
	return dto.ToProfile(profileDto), nil
}
