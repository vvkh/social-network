package repository

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/profiles/entity"
	"github.com/vvkh/social-network/internal/domain/profiles/repository/dto"
)

func (r *repo) ListProfiles(ctx context.Context, firstNamePrefix string, lastNamePrefix string, limit int) ([]entity.Profile, bool, error) {
	var profileDtos []dto.Profile
	var err error

	// TODO: rewrite with squirrel
	// TODO: add limit
	if firstNamePrefix == "" && lastNamePrefix == "" {
		query := `SELECT * FROM profiles`
		err = r.db.SelectContext(ctx, &profileDtos, query)
	} else if firstNamePrefix == "" {
		query := `SELECT * FROM profiles WHERE last_name LIKE ?`
		err = r.db.SelectContext(ctx, &profileDtos, query, lastNamePrefix+"%")
	} else {
		query := `SELECT * FROM profiles WHERE first_name LIKE ? AND last_name LIKE ?`
		err = r.db.SelectContext(ctx, &profileDtos, query, firstNamePrefix+"%", lastNamePrefix+"%")
	}
	if err != nil {
		return nil, false, err
	}
	return dto.ToProfiles(profileDtos), false, nil
}
