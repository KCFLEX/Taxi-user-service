package services

import (
	"context"
	"strconv"
	"time"

	"github.com/KCFLEX/Taxi-user-service/errorpac"
	"github.com/KCFLEX/Taxi-user-service/internal/handlers/models"
	"github.com/KCFLEX/Taxi-user-service/internal/repository/entity"
	"golang.org/x/crypto/bcrypt"
)

type Token interface {
	GenerateToken(ctx context.Context, userID string, duration time.Duration) (string, error)
}

type Repository interface {
	UserExists(ctx context.Context, user entity.User) (bool, error)
	CreateUser(ctx context.Context, user entity.User) error
	UserPhoneExists(ctx context.Context, user entity.User) (entity.User, error)
}

// remember you were doing transformation transforming models.UserInfo to entity.User
type Service struct {
	repo  Repository
	token Token
}

func New(repo Repository, token Token) *Service {
	return &Service{
		repo:  repo,
		token: token,
	}
}

func (srv *Service) SignUP(ctx context.Context, User models.UserInfo) error {
	// next we do sessions and tokens
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

func (srv *Service) SignIN(ctx context.Context, user models.UserInfo) (string, error) {

	checkUser := entity.User{
		Phone: user.PhoneNO,
	}

	userId, err := srv.repo.UserPhoneExists(ctx, checkUser)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userId.Password), []byte(user.Password))
	if err != nil {
		return "", errorpac.ErrPasswordInvalid
	}

	userIdStr := strconv.Itoa(userId.ID)

	// generate token and return token
	tokenStr, err := srv.token.GenerateToken(ctx, userIdStr, 24*time.Hour)

	if err != nil {
		return "", &errorpac.CustomErr{
			OriginalErr: err,
			SpecificErr: errorpac.ErrTokenGenFail,
		}
	}

	return tokenStr, nil

}
