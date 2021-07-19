package repository

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/profiles/entity"
	"github.com/vvkh/social-network/internal/domain/profiles/repository/dto"
)

func (r *repo) GetByUserID(ctx context.Context, id uint64) ([]entity.Profile, error) {
	query := `SELECT * FROM profiles WHERE user_id = ?`
	var profilesDto []dto.Profile
	err := r.db.SelectContext(ctx, &profilesDto, query, id)
	if err != nil {
		return nil, err
	}
	profiles := make([]entity.Profile, 0, len(profilesDto))
	for _, profileDto := range profilesDto {
		profiles = append(profiles, dto.ToProfile(profileDto))
	}
	return profiles, nil
}
