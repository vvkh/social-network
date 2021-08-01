package repository

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/profiles/entity"
	"github.com/vvkh/social-network/internal/domain/profiles/repository/dto"
)

func (r *repo) GetByName(ctx context.Context, firstNamePrefix string, lastNamePrefix string) ([]entity.Profile, error) {
	query := `SELECT * FROM profiles WHERE first_name LIKE ? AND last_name LIKE ?`

	var profileDtos []dto.Profile
	if err := r.db.SelectContext(ctx, &profileDtos, query, firstNamePrefix+"%", lastNamePrefix+"%"); err != nil {
		return nil, err
	}
	return dto.ToProfiles(profileDtos), nil
}
