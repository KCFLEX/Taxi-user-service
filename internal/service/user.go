package services

import (
	"context"

	"github.com/KCFLEX/Taxi-user-service/errorpac"
	"github.com/KCFLEX/Taxi-user-service/internal/handlers/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user models.UserInfo) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (srv *Service) SignUP(ctx context.Context, User models.UserInfo) error {
	// next we do sessions and tokens
	// check if email exists already
	err := User.Required()
	if err != nil {
		return err
	}

	err = User.Validate()
	if err != nil {
		return err
	}
	passwordEncode, err := models.HashPass(User.Password)
	if err != nil {
		return err
	}

	User.Password = string(passwordEncode)

	err = srv.repo.CreateUser(ctx, User)
	if err != nil {
		return &errorpac.CustomErr{
			SpecificErr: errorpac.ErrCreateUserFail,
			OriginalErr: err,
		}
	}

	return nil
}
