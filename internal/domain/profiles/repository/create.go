package repository

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/profiles/entity"
	"github.com/vvkh/social-network/internal/domain/profiles/repository/dto"
)

func (r *repo) CreateProfile(ctx context.Context, profile entity.Profile) (uint64, error) {
	query := `
	INSERT INTO profiles (first_name, last_name, age, sex, about, location, user_id)
	VALUES (:first_name, :last_name, :age, :sex, :about, :location, :user_id)`

	result, err := r.db.NamedExecContext(ctx, query, dto.FromProfile(profile))
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return uint64(id), err
}
