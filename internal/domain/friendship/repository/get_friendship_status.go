package repository

import (
	"context"
	"database/sql"

	"github.com/vvkh/social-network/internal/domain/friendship/entity"
	"github.com/vvkh/social-network/internal/domain/friendship/repository/dto"
)

func (r *repo) GetFriendshipStatus(ctx context.Context, one uint64, other uint64) (entity.FriendshipStatus, error) {
	query := `
	SELECT requested_from, requested_to, state FROM friendship
	WHERE requested_from IN (?, ?) and requested_to IN (?, ?)
	`
	var friendship dto.Friendship
	err := r.db.GetContext(ctx, &friendship, query, one, other, one, other)
	if err == sql.ErrNoRows {
		return entity.FriendshipStatus{State: entity.StateNone}, nil
	}
	if err != nil {
		return entity.FriendshipStatus{}, err
	}
	return dto.ToModel(friendship), nil
}
