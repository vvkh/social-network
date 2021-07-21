package profiles

import "github.com/vvkh/social-network/internal/domain/profiles/entity"

type Dto struct {
	FirstName string
	LastName  string
	ID        uint64
}

type Context struct {
	UserID   uint64
	Profiles []Dto
}

func dtoFromModel(profile entity.Profile) Dto {
	return Dto{
		ID:        profile.ID,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
	}
}
