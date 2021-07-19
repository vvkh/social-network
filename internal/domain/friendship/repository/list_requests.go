package repository

import (
	"context"
)

func (r *repo) ListPendingRequests(ctx context.Context, profileID uint64) ([]uint64, error) {
	query := `
	SELECT requested_from FROM friendship WHERE requested_to = ? AND state = "created"
    `
	var ids []uint64
	err := r.db.SelectContext(ctx, &ids, query, profileID)
	if err != nil {
		return nil, err
	}
	return ids, nil
}
