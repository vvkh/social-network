package usecase

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/vvkh/social-network/internal/profiles/entity"
	"github.com/vvkh/social-network/internal/profiles/mocks"
)

func Test_usecase_CreateProfile(t *testing.T) {
	type args struct {
		userID    uint64
		firstName string
		lastName  string
		age       uint8
		location  string
		sex       string
		about     string
	}
	tests := []struct {
		name                   string
		args                   args
		repoMockCreateResponse uint64
		want                   entity.Profile
		wantErr                bool
	}{
		{
			name: "empty_name",
			args: args{
				userID:   1,
				lastName: "Doe",
				age:      16,
				sex:      "male",
			},
			wantErr: true,
		},
		{
			name: "empty_last_name",
			args: args{
				userID:    1,
				firstName: "John",
				age:       16,
				sex:       "male",
			},
			wantErr: true,
		},
		{
			name: "empty_age",
			args: args{
				userID:    1,
				firstName: "John",
				lastName:  "Doe",
				sex:       "male",
			},
			wantErr: true,
		},
		{
			name: "empty_sex",
			args: args{
				userID:    1,
				firstName: "John",
				lastName:  "Doe",
				age:       16,
			},
			wantErr: true,
		},
		{
			name: "invalid_sex",
			args: args{
				userID:    1,
				firstName: "John",
				lastName:  "Doe",
				age:       16,
				sex:       "123",
			},
			wantErr: true,
		},
		{
			name: "valid_female",
			args: args{
				userID:    1,
				firstName: "Johna",
				lastName:  "Doe",
				age:       16,
				sex:       "female",
			},
			repoMockCreateResponse: 42,
			wantErr:                false,
			want: entity.Profile{
				UserID:    1,
				ID:        42,
				FirstName: "Johna",
				LastName:  "Doe",
				Age:       16,
				Sex:       "female",
				About:     "",
				Location:  "",
			},
		},
		{
			name: "valid_male",
			args: args{
				userID:    1,
				firstName: "John",
				lastName:  "Doe",
				age:       16,
				sex:       "male",
			},
			repoMockCreateResponse: 42,
			wantErr:                false,
			want: entity.Profile{
				UserID:    1,
				ID:        42,
				FirstName: "John",
				LastName:  "Doe",
				Age:       16,
				Sex:       "male",
				About:     "",
				Location:  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repoMock := mocks.NewMockRepository(ctrl)
			u := New(repoMock)

			ctx := context.Background()
			if tt.repoMockCreateResponse != 0 {
				repoMock.EXPECT().CreateProfile(ctx, gomock.Any()).Return(tt.repoMockCreateResponse, nil)
			}
			got, err := u.CreateProfile(ctx, tt.args.userID, tt.args.firstName, tt.args.lastName, tt.args.age, tt.args.location, tt.args.sex, tt.args.about)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("CreateProfile() got = %v, want %v", got, tt.want)
			}

			repoMock.EXPECT().GetByID(ctx, got.ID).Return([]entity.Profile{got}, nil)
			fetchedByID, err := u.GetByID(ctx, got.ID)
			if err != nil {
				t.Fatalf("unexpected error while fetching profile: %v", err)
			}
			if !reflect.DeepEqual([]entity.Profile{got}, fetchedByID) {
				t.Fatalf("expected to get by id = %v, got = %v", got, fetchedByID)
			}
		})
	}
}
