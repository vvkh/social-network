package repository

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/friendship/repository/dto"
)

func (r *repo) ListFriends(ctx context.Context, profileID uint64) ([]uint64, error) {
	query := `
	SELECT requested_from, requested_to FROM friendship
	WHERE state = "accepted" and (requested_from = ? or requested_to = ?)
    `
	var relations []dto.Friendship
	err := r.db.SelectContext(ctx, &relations, query, profileID, profileID)
	if err != nil {
		return nil, err
	}

	allExceptRequestedID := make([]uint64, 0, len(relations))
	for _, relation := range relations {
		if relation.RequestedTo != profileID {
			allExceptRequestedID = append(allExceptRequestedID, relation.RequestedTo)
		}
		if relation.RequestedFrom != profileID {
			allExceptRequestedID = append(allExceptRequestedID, relation.RequestedFrom)
		}
	}
	return allExceptRequestedID, nil
}
