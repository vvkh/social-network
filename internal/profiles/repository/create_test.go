package repository

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/joho/godotenv"

	"github.com/vvkh/social-network/internal/profiles/entity"
)

func Test_repo_CreateProfile(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Errorf("should not be error while parsing dotenv but got: %v", err)
		return
	}
	if os.Getenv("SKIP_DB_TEST") == "1" {
		t.SkipNow()
	}
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
			require.NoError(t, err)

			ctx := context.Background()
			gotID, err := r.CreateProfile(ctx, tt.profile)
			require.NoError(t, err)

			gotProfile, err := r.GetByID(ctx, gotID)
			require.NoError(t, err)

			wantProfile := tt.profile
			wantProfile.ID = gotID

			if !reflect.DeepEqual(gotProfile, []entity.Profile{wantProfile}) {
				t.Errorf("want profile = %v, gotID = %v", tt.profile, gotProfile)
			}
		})
	}
}
