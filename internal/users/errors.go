package users

import (
	"errors"
)

var (
	AuthenticationFailed = errors.New("unable to authenticate user with provided credentials")
)
