package dto

import "github.com/vvkh/social-network/internal/domain/profiles/entity"

type Profile struct {
	ID        uint64 `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Age       uint8  `db:"age"`
	Sex       string `db:"sex"`
	About     string `db:"about"`
	Location  string `db:"location"`
	UserID    uint64 `db:"user_id"`
}

func FromProfile(model entity.Profile) Profile {
	return Profile{
		ID:        model.ID,
		UserID:    model.UserID,
		FirstName: model.FirstName,
		LastName:  model.LastName,
		Age:       model.Age,
		Sex:       string(model.Sex),
		About:     model.About,
		Location:  model.Location,
	}
}

func ToProfile(dto Profile) entity.Profile {
	return entity.Profile{
		ID:        dto.ID,
		UserID:    dto.UserID,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Age:       dto.Age,
		Sex:       entity.Sex(dto.Sex),
		About:     dto.About,
		Location:  dto.Location,
	}
}

func ToProfiles(dtos []Profile) []entity.Profile {
	profiles := make([]entity.Profile, 0, len(dtos))
	for _, dto := range dtos {
		profiles = append(profiles, ToProfile(dto))
	}
	return profiles
}
