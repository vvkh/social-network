package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAccessToken(t *testing.T) {
	tests := []struct {
		name      string
		userID    uint64
		profileID uint64
		expiresAt time.Time
		key       string
	}{
		{
			name:      "encode_decode_token",
			userID:    99,
			profileID: 999,
			expiresAt: time.Now().UTC().Add(time.Minute).Truncate(time.Second),
			key:       "secret",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := AccessToken{
				UserID:    tt.userID,
				ProfileID: tt.profileID,
				ExpiresAt: tt.expiresAt,
			}
			rawToken, err := token.ToString(tt.key)
			require.NoError(t, err)

			parsedToken, err := Parse(rawToken, tt.key)
			require.NoError(t, err)
			require.Equal(t, token, parsedToken)
		})
	}
}
