package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/vvkh/social-network/internal/profiles/entity"
)

func Test_repo_CreateProfile(t *testing.T) {

	tests := []struct {
		name    string
		profile entity.Profile
	}{
		{
			name: "create_profile",
			profile: entity.Profile{
				FirstName: "John",
				LastName:  "Doe",
				Age:       14,
				Sex:       "male",
				About:     "im John",
				Location:  "USA",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewDefault()
			if err != nil {
				t.Errorf("should not be error while creating repository, but got = %v", err)
				return
			}

			ctx := context.Background()
			gotID, err := r.CreateProfile(ctx, tt.profile)
			if err != nil {
				t.Errorf("should not be error while creating profile, but got = %v", err)
				return
			}

			gotProfile, err := r.GetByID(ctx, gotID)
			if err != nil {
				t.Errorf("should not be error while getting profile")
				return
			}
			wantProfile := tt.profile
			wantProfile.ID = gotID

			if !reflect.DeepEqual(gotProfile, wantProfile) {
				t.Errorf("want profile = %v, gotID = %v", tt.profile, gotProfile)
			}
		})
	}
}
