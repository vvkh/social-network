package chat

import (
	"fmt"

	"github.com/vvkh/social-network/internal/domain/chats/entity"
	profileEntity "github.com/vvkh/social-network/internal/domain/profiles/entity"
	"github.com/vvkh/social-network/internal/navbar"
)

type Contex struct {
	Navbar   *navbar.Context
	Chat     ChatDto
	Messages []MessageDto
}

type ChatDto struct {
	ID    uint64
	Title string
}

type MessageDto struct {
	Content string
	SentAt  string
	Author  string
}

func messagesDtoFromModel(messages []entity.Message, profiles []profileEntity.Profile) []MessageDto {
	authorById := make(map[uint64]profileEntity.Profile, len(profiles))

	for _, profile := range profiles {
		authorById[profile.ID] = profile
	}

	dtos := make([]MessageDto, 0, len(messages))
	for _, message := range messages {
		dtos = append(dtos, messageDtoFromModel(message, authorById[message.AuthorProfileID]))
	}
	return dtos
}

func messageDtoFromModel(message entity.Message, author profileEntity.Profile) MessageDto {
	return MessageDto{
		Content: message.Content,
		SentAt:  message.SentAt.Format("2006 01 02 15:04"),
		Author:  fmt.Sprintf("%s %s", author.FirstName, author.LastName),
	}
}

func chatDroFromModel(chat entity.Chat) ChatDto {
	return ChatDto{
		ID:    chat.ID,
		Title: chat.Title,
	}
}
