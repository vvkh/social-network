package usecase

import (
	"fmt"

	"github.com/vvkh/social-network/internal/profiles/entity"
)

func (u *usecase) CreateProfile(firstName string, lastName string, age uint8, location string, sex string, about string) (entity.Profile, error) {
	if firstName == "" {
		return entity.Profile{}, fmt.Errorf("first name cant be empty")
	}
	if lastName == "" {
		return entity.Profile{}, fmt.Errorf("last name cant be empty")
	}
	if age == 0 {
		return entity.Profile{}, fmt.Errorf("age cant be empty")
	}
	if !entity.ValidSex[entity.Sex(sex)] {
		return entity.Profile{}, fmt.Errorf("invalid sex %s", sex)
	}
	profile := entity.Profile{
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
		Sex:       entity.Sex(sex),
		About:     about,
		Location:  location,
	}
	id, err := u.repository.CreateProfile(profile)
	if err != nil {
		return entity.Profile{}, fmt.Errorf("repository.CreateProfile error: %w", err)
	}
	profile.ID = id
	return profile, nil
}
