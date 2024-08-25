package services

import (
	"context"

	"github.com/KCFLEX/Taxi-user-service/errorpac"
	"github.com/KCFLEX/Taxi-user-service/internal/handlers/models"
	"github.com/KCFLEX/Taxi-user-service/internal/repository/entity"
)

type Repository interface {
	UserExists(ctx context.Context, user entity.User) (bool, error)
	CreateUser(ctx context.Context, user entity.User) error
}

// remember you were doing transformation transforming models.UserInfo to entity.User
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

	newUser := entity.User{
		Name:     User.Name,
		Phone:    User.PhoneNO,
		Email:    User.Email,
		Password: User.Password,
	}

	check, err := srv.repo.UserExists(ctx, newUser)
	if err != nil {
		return err
	}

	if check {
		return errorpac.ErrUserAlreadyExists
	}

	err = srv.repo.CreateUser(ctx, newUser)
	if err != nil {
		return &errorpac.CustomErr{
			SpecificErr: errorpac.ErrCreateUserFail,
			OriginalErr: err,
		}
	}

	return nil
}
