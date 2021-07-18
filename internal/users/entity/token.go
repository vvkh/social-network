package entity

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type AccessToken struct {
	UserID    uint64
	ProfileID uint64
	ExpiresAt time.Time
}

func (t AccessToken) ToString(key string) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     t.UserID,
		"profile": t.ProfileID,
		"exp":     t.ExpiresAt.Unix(),
	}).SignedString([]byte(key))
	return token, err
}

func Parse(raw string, key string) (AccessToken, error) {
	parsed, err := jwt.Parse(raw, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		return AccessToken{}, err
	}
	if claims, ok := parsed.Claims.(jwt.MapClaims); ok && parsed.Valid {
		parsed, ok := parseFromClaims(claims)
		if !ok {
			return AccessToken{}, errors.New("Uknown error while reading claims")
		}
		return parsed, nil
	}
	return AccessToken{}, errors.New("Unknown error while parsing")
}

func parseFromClaims(claims jwt.MapClaims) (AccessToken, bool) {
	userID, ok := claims["sub"]
	if !ok {
		return AccessToken{}, false
	}
	userIDConverted, ok := userID.(float64)
	if !ok {
		return AccessToken{}, false
	}
	profile, ok := claims["profile"]
	if !ok {
		return AccessToken{}, false
	}
	profileConverted, ok := profile.(float64)
	if !ok {
		return AccessToken{}, false
	}
	expiresAt, ok := claims["exp"]
	if !ok {
		return AccessToken{}, false
	}
	expiredAtConverted, ok := expiresAt.(float64)
	if !ok {
		return AccessToken{}, false
	}
	return AccessToken{
		UserID:    uint64(userIDConverted),
		ProfileID: uint64(profileConverted),
		ExpiresAt: time.Unix(int64(expiredAtConverted), 0).UTC(),
	}, true
}
