package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var (
	getChatQuery = `
SELECT id FROM chats JOIN chat_members
ON chat_members.chat_id = chats.id
WHERE member_profile_id IN (?)
GROUP BY chats.id
HAVING count(1) = ?
`
	createChatQuery     = `INSERT INTO chats (title) VALUES (?)`
	addChatMembersQuery = `INSERT INTO chat_members (chat_id, member_profile_id) VALUES (?, ?), (?, ?)`
)

func (r *repo) GetOrCreateChat(ctx context.Context, chatTitle string, oneProfileID uint64, otherProfileID uint64) (uint64, error) {
	profiles := []interface{}{oneProfileID, otherProfileID}
	getChatSql, args, err := sqlx.In(getChatQuery, profiles, len(profiles))
	if err != nil {
		return 0, err
	}
	var chatID uint64
	err = r.db.GetContext(ctx, &chatID, getChatSql, args...)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	if err != sql.ErrNoRows {
		return chatID, nil
	}

	txn, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer txn.Rollback()

	result, err := txn.ExecContext(ctx, createChatQuery, chatTitle)
	if err != nil {
		return 0, err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	chatID = uint64(lastInsertID)

	_, err = txn.ExecContext(ctx, addChatMembersQuery, chatID, oneProfileID, chatID, otherProfileID)
	if err != nil {
		return 0, err
	}

	err = txn.Commit()
	return chatID, err
}
