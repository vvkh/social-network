package repository

import (
	"context"
)

func (r *repo) StopFriendship(ctx context.Context, userID uint64, otherUserID uint64) error {
	query := `
	DELETE FROM friendship
    WHERE requested_from = ? and requested_to = ? OR requested_from = ? and requested_to = ?
	AND state = "accepted"`
	_, err := r.db.ExecContext(ctx, query, userID, otherUserID, otherUserID, userID)
	return err
}
