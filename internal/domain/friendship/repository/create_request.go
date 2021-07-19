package repository

import (
	"context"
)

func (r *repo) CreateRequest(ctx context.Context, profileIDFrom uint64, profileIDTo uint64) error {
	query := `
	INSERT INTO friendship (requested_from, requested_to)
	VALUES (?, ?)`
	_, err := r.db.ExecContext(ctx, query, profileIDFrom, profileIDTo)
	return err
}
