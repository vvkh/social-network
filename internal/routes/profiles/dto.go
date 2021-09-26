package profiles

import (
	"github.com/vvkh/social-network/internal/domain/profiles/entity"
	"github.com/vvkh/social-network/internal/navbar"
)

type Dto struct {
	FirstName string
	LastName  string
	ID        uint64
	UserID    uint64
}

type Filters struct {
	FirstName string
	LastName  string
}

type ShowMoreLimit struct {
	NextLimit int
}
type Context struct {
	Navbar          *navbar.Context
	Self            Dto
	Profiles        []Dto
	Filters         Filters
	DisplayShowMore *ShowMoreLimit
}

func dtoFromModel(profile entity.Profile) Dto {
	return Dto{
		UserID:    profile.UserID,
		ID:        profile.ID,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
	}
}

func dtoFromModels(profiles []entity.Profile) []Dto {
	dtos := make([]Dto, 0, len(profiles))
	for _, profile := range profiles {
		dtos = append(dtos, dtoFromModel(profile))
	}
	return dtos
}
