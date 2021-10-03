package chats

import (
	"github.com/vvkh/social-network/internal/domain/chats/entity"
	"github.com/vvkh/social-network/internal/navbar"
)

type Contex struct {
	Navbar *navbar.Context
	Chats  []ChatDto
}

type ChatDto struct {
	ID                  uint64
	Title               string
	UnreadMessagesCount int64
}

func dtoFromModels(chats []entity.Chat) []ChatDto {
	dtos := make([]ChatDto, 0, len(chats))
	for _, chat := range chats {
		dtos = append(dtos, dtoFromModel(chat))
	}
	return dtos
}

func dtoFromModel(chat entity.Chat) ChatDto {
	return ChatDto{
		ID:                  chat.ID,
		Title:               chat.Title,
		UnreadMessagesCount: chat.UnreadMessagesCount,
	}
}
