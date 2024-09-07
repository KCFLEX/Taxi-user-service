package models

import (
	"regexp"
	"strings"

	"github.com/KCFLEX/Taxi-user-service/errorpac"
	"golang.org/x/crypto/bcrypt"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	phoneRegex = regexp.MustCompile(`^\+?\d{1,3}[-.\s]?\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}$`)
)

type UserInfo struct {
	Name     string `json:"name"`
	PhoneNO  string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *UserInfo) Validate() error {
	if !emailRegex.MatchString(s.Email) {
		return errorpac.ErrInvalidEmail
	}
	if !phoneRegex.MatchString(s.PhoneNO) {
		return errorpac.ErrInvaiidPhone
	}
	return nil
}

func (s *UserInfo) Required() error {

	if strings.TrimSpace(s.Name) == "" {
		return errorpac.ErrNameRequired
	}

	// Check if PhoneNO is empty
	if strings.TrimSpace(s.PhoneNO) == "" {
		return errorpac.ErrInvaiidPhone
	}

	// Check if Email is empty
	if strings.TrimSpace(s.Email) == "" {
		return errorpac.ErrInvalidEmail
	}

	// Check if Password is empty
	if strings.TrimSpace(s.Password) == "" {
		return errorpac.ErrPasswordRequired
	}

	return nil

}

func HashPass(password string) ([]byte, error) {
	passwordEncode, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return []byte{}, &errorpac.CustomErr{
			OriginalErr: err,
			SpecificErr: errorpac.ErrPassHashFail,
		}
	}

	return passwordEncode, nil
}
