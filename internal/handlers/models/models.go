package models

import (
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	phoneRegex = regexp.MustCompile(`^\+?\d{1,3}[-.\s]?\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}$`)
)

type errors map[string]string

type SignUPInfo struct {
	Name     string `json:"name"`
	PhoneNO  string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Errors   map[string]string
}

func (s *SignUPInfo) Validate() bool {
	s.Errors = make(errors)
	if !emailRegex.MatchString(s.Email) {
		s.Errors["email"] = "please enter valid email"
	}
	if !phoneRegex.MatchString(s.PhoneNO) {
		s.Errors["phone"] = "please insert valid phone number"
	}

	return len(s.Errors) == 0
}

func (s *SignUPInfo) Required() errors {
	s.Errors = make(errors)

	if strings.TrimSpace(s.Name) == "" {
		s.Errors["name"] = "Name is required"
	}

	// Check if PhoneNO is empty
	if strings.TrimSpace(s.PhoneNO) == "" {
		s.Errors["phone"] = "Phone number is required"
	}

	// Check if Email is empty
	if strings.TrimSpace(s.Email) == "" {
		s.Errors["email"] = "Email is required"
	}

	// Check if Password is empty
	if strings.TrimSpace(s.Password) == "" {
		s.Errors["password"] = "Password is required"
	}

	return s.Errors

}
