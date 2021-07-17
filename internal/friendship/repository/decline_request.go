package repository

import (
	"context"
)

func (r *repo) DeclineRequest(ctx context.Context, userIDFrom uint64, userIDTo uint64) error {
	query := `
	UPDATE friendship SET state = "declined" 
	WHERE requested_from = ? and requested_to = ? and state = "created"
    `
	_, err := r.db.ExecContext(ctx, query, userIDFrom, userIDTo)
	return err
}
