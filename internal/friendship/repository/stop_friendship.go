package repository

import (
	"context"
)

func (r *repo) StopFriendship(ctx context.Context, profileID uint64, otherProfileID uint64) error {
	query := `
	DELETE FROM friendship
    WHERE requested_from = ? and requested_to = ? OR requested_from = ? and requested_to = ?
	AND state = "accepted"`
	_, err := r.db.ExecContext(ctx, query, profileID, otherProfileID, otherProfileID, profileID)
	return err
}
