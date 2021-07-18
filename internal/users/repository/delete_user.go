package repository

import "context"

func (r *repo) DeleteUser(ctx context.Context, userID uint64) error {
	query := `
	DELETE FROM users WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
