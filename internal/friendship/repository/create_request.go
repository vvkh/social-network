package repository

import (
	"context"
)

func (r *repo) CreateRequest(ctx context.Context, userIDFrom uint64, userIDTo uint64) error {
	query := `
	INSERT INTO friendship (requested_from, requested_to)
	VALUES (?, ?)`
	_, err := r.db.ExecContext(ctx, query, userIDFrom, userIDTo)
	return err
}
