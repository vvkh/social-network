package usecase

import (
	"github.com/vvkh/social-network/internal/domain/chats"
)

type uc struct {
	repo chats.Repository
}

func New(repo chats.Repository) *uc {
	return &uc{repo: repo}
}
