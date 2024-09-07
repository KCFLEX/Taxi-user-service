package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/KCFLEX/Taxi-user-service/errorpac"
	"github.com/KCFLEX/Taxi-user-service/internal/handlers/models"
	"github.com/KCFLEX/Taxi-user-service/internal/repository/entity"

	"golang.org/x/crypto/bcrypt"
)

type Token interface {
	GenerateToken(ctx context.Context, userID string, duration time.Duration) (string, error)
	ValidateToken(ctx context.Context, tokenString string) error

	ParseToken(ctx context.Context, tokenStr string) (string, error)
}

type Repository interface {
	UserExists(ctx context.Context, user entity.User) (bool, error)
	CreateUser(ctx context.Context, user entity.User) error
	UserPhoneExists(ctx context.Context, user entity.User) (entity.User, error)
	//redis methods below
	StoreTokenInRedis(ctx context.Context, userID string, token string, expiration time.Duration) error //store token in redis
	ValidateTokenRedis(ctx context.Context, token string, userID string) error
	// redis method for jwt token validation
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
	fmt.Println(userIdStr)
	// generate token and return token
	tokenStr, err := srv.token.GenerateToken(ctx, userIdStr, 24*time.Minute)

	if err != nil {
		return "", &errorpac.CustomErr{
			OriginalErr: err,
			SpecificErr: errorpac.ErrTokenGenFail,
		}
	}

	// store token in redis
	err = srv.repo.StoreTokenInRedis(ctx, userIdStr, tokenStr, 2*time.Hour)
	if err != nil {
		return "", &errorpac.CustomErr{
			OriginalErr: err,
			SpecificErr: errorpac.ErrFailToStoreToken,
		}
	}

	return tokenStr, nil

}

func (srv *Service) VerifyToken(ctx context.Context, token string) (string, error) {
	userId, err := srv.token.ParseToken(ctx, token)
	if err != nil {
		return "", err
	}
	err = srv.token.ValidateToken(ctx, token)
	if err != nil {
		return "", err
	}
	return userId, nil
}

// Next create protected routes for authenticated user access only

func (srv *Service) CheckTokenInRedis(ctx context.Context, token string) error {
	userId, err := srv.token.ParseToken(ctx, token)
	if err != nil {
		return err
	}
	return srv.repo.ValidateTokenRedis(ctx, token, userId)
}
