package repository

import (
	"context"
)

func (r *repo) AcceptRequest(ctx context.Context, profileIDFrom uint64, profileIDTo uint64) error {
	query := `
	UPDATE friendship SET state = "accepted"
    WHERE requested_from = ? AND requested_to = ? AND state in ("created", "declined") -- declined requests can be still accepted later
    `
	_, err := r.db.ExecContext(ctx, query, profileIDFrom, profileIDTo)
	return err
}
