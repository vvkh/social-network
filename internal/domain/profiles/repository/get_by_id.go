package repository

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/vvkh/social-network/internal/domain/profiles/entity"
	"github.com/vvkh/social-network/internal/domain/profiles/repository/dto"
)

func (r *repo) GetByID(ctx context.Context, id ...uint64) ([]entity.Profile, error) {
	query, args, err := sqlx.In(`SELECT * FROM profiles WHERE id in (?)`, id)
	if err != nil {
		return nil, err
	}

	var profilesDto []dto.Profile
	err = r.db.SelectContext(ctx, &profilesDto, query, args...)
	if err != nil {
		return nil, err
	}

	profiles := make([]entity.Profile, 0, len(profilesDto))
	for _, profileDto := range profilesDto {
		profiles = append(profiles, dto.ToProfile(profileDto))
	}
	return profiles, nil
}
