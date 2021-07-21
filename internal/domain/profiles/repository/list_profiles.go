package repository

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/profiles/entity"
	"github.com/vvkh/social-network/internal/domain/profiles/repository/dto"
)

func (r *repo) ListProfiles(ctx context.Context) ([]entity.Profile, error) {
	query := `SELECT * FROM profiles`
	var profileDtos []dto.Profile
	if err := r.db.SelectContext(ctx, &profileDtos, query); err != nil {
		return nil, err
	}
	return dto.ToProfiles(profileDtos), nil
}
