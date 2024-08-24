package errorpac

import (
	"errors"
	"fmt"
)

var (
	ErrPassHashFail     = errors.New("error generating password hash")
	ErrCreateUserFail   = errors.New("error creating user")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrInvaiidPhone     = errors.New("invalid phone format")
	ErrNameRequired     = errors.New("name is required")
	ErrPasswordRequired = errors.New("password is required")
	ErrDuplicateEmail   = errors.New("email already exists")
)

type CustomErr struct {
	SpecificErr error
	OriginalErr error
}

func (c *CustomErr) Error() string {
	return fmt.Sprintf("%v, %v", c.SpecificErr, c.OriginalErr)
}

func (c *CustomErr) Is(target error) bool {
	return c.SpecificErr == target
}
