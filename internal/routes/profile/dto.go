package profile

import (
	"github.com/vvkh/social-network/internal/domain/profiles/entity"
	"github.com/vvkh/social-network/internal/navbar"
)

type Context struct {
	Navbar                          *navbar.Context
	Self                            ProfileDto
	Profile                         ProfileDto
	AreFriends                      bool
	IsWaitingFriendshipApproval     bool
	HasNotConfirmedFriendship       bool
	FriendshipRequestDeclined       bool
	FriendshipRequestDeclinedBySelf bool
}

type ProfileDto struct {
	ID        uint64
	UserID    uint64
	FirstName string
	LastName  string
	Age       uint8
	Sex       string
	About     string
	Location  string
}

func dtoFromModel(profile entity.Profile) ProfileDto {
	return ProfileDto{
		ID:        profile.ID,
		UserID:    profile.UserID,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Age:       profile.Age,
		Sex:       string(profile.Sex),
		About:     profile.About,
		Location:  profile.Location,
	}
}
