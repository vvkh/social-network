package users

import (
	"errors"
)

var (
	AuthenticationFailed = errors.New("unable to authenticate user with provided credentials")
	EmptyCredentials     = errors.New("got empty values for required credentials")
)
