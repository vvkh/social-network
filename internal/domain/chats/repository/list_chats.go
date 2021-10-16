package repository

import (
	"context"

	"github.com/vvkh/social-network/internal/domain/chats/entity"
)

var (
	query = `SELECT id, title FROM chats JOIN chat_members ON chats.id = chat_members.chat_id WHERE member_profile_id = ?`
)

type dto struct {
	ID    uint64
	Title string
}

func (r *repo) ListChats(ctx context.Context, profileID uint64) ([]entity.Chat, error) {
	var chatDtos []dto

	err := r.db.SelectContext(ctx, &chatDtos, query, profileID)
	if err != nil {
		return nil, err
	}
	return dtoToModels(chatDtos), nil
}

func dtoToModels(dtos []dto) []entity.Chat {
	models := make([]entity.Chat, 0, len(dtos))
	for _, dto := range dtos {
		models = append(models, entity.Chat{
			Title: dto.Title,
			ID:    dto.ID,
		})
	}
	return models
}
