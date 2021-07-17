package repository

import (
	"context"
)

func (r *repo) AcceptRequest(ctx context.Context, userIDFrom uint64, userIDTo uint64) error {
	query := `
	UPDATE friendship SET state = "accepted"
    WHERE requested_from = ? AND requested_to = ? AND state = "created"
    `
	_, err := r.db.ExecContext(ctx, query, userIDFrom, userIDTo)
	return err
}
