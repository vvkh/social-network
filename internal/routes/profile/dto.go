package profile

import "github.com/vvkh/social-network/internal/domain/profiles/entity"

type Dto struct {
	UserID    uint64
	FirstName string
	LastName  string
	Age       uint8
	Sex       string
	About     string
	Location  string
}

func dtoFromModel(profile entity.Profile) Dto {
	return Dto{
		UserID:    profile.UserID,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Age:       profile.Age,
		Sex:       string(profile.Sex),
		About:     profile.About,
		Location:  profile.Location,
	}
}
