package friends_requests

import "github.com/vvkh/social-network/internal/domain/profiles/entity"

type Contex struct {
	Self            ProfileDto
	PendingRequests []ProfileDto
}

type ProfileDto struct {
	ID        uint64
	UserID    uint64
	FirstName string
	LastName  string
}

func dtoFromModel(profile entity.Profile) ProfileDto {
	return ProfileDto{
		ID:        profile.ID,
		UserID:    profile.UserID,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
	}
}

func dtoFromModels(profiles []entity.Profile) []ProfileDto {
	dtos := make([]ProfileDto, 0, len(profiles))
	for _, model := range profiles {
		dtos = append(dtos, dtoFromModel(model))
	}
	return dtos
}
