package profiles

import "github.com/vvkh/social-network/internal/domain/profiles/entity"

type Dto struct {
	FirstName string
	LastName  string
	ID        uint64
	UserID    uint64
}

type Context struct {
	Self     Dto
	Profiles []Dto
}

func dtoFromModel(profile entity.Profile) Dto {
	return Dto{
		UserID:    profile.UserID,
		ID:        profile.ID,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
	}
}
