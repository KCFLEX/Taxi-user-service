package services

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/KCFLEX/Taxi-user-service/errorpac"
	"github.com/KCFLEX/Taxi-user-service/internal/handlers/models"
	"github.com/KCFLEX/Taxi-user-service/internal/repository/entity"
	"github.com/KCFLEX/Taxi-user-service/internal/service/mockrepo"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestSignUP(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newRepoMock := mockrepo.NewMockRepository(ctrl)
	newTokenMock := mockrepo.NewMockToken(ctrl)
	srv := New(newRepoMock, newTokenMock)
	user := entity.User{
		Name:     "vincent",
		Phone:    "+1-245-213-2222",
		Email:    "test@gamil.com",
		Password: "kfidrofjkeoirfjkl",
	}
	// Succesfull (user does not exist)
	newRepoMock.EXPECT().UserExists(gomock.Any(), gomock.AssignableToTypeOf(entity.User{})).DoAndReturn(func(ctx context.Context, u entity.User) (bool, error) {
		if u.Name != user.Name || u.Phone != user.Phone || u.Email != user.Email {
			t.Errorf("expected user details to match, got %v", u)
		}
		return false, nil
	})
	newRepoMock.EXPECT().CreateUser(gomock.Any(), gomock.AssignableToTypeOf(entity.User{})).DoAndReturn(func(ctx context.Context, u entity.User) error {
		if u.Name != user.Name || u.Phone != user.Phone || u.Email != user.Email {
			t.Errorf("expected user details to match, got %v", u)
		}
		return nil
	})

	sameUser := models.UserInfo{
		Name:     "vincent",
		PhoneNO:  "+1-245-213-2222",
		Email:    "test@gamil.com",
		Password: "kfidrofjkeoirfjkl",
	}
	err := srv.SignUP(context.Background(), sameUser)

	if err != nil {
		t.Errorf("didn't expect any error but got error: %v", err)
	}
	// user exists
	newRepoMock.EXPECT().UserExists(gomock.Any(), gomock.AssignableToTypeOf(entity.User{})).DoAndReturn(func(ctx context.Context, u entity.User) (bool, error) {
		if u.Name != user.Name || u.Phone != user.Phone || u.Email != user.Email {
			t.Errorf("expected user details to match, got %v", u)
		}
		return true, nil
	})
	expectedErr := errorpac.ErrUserAlreadyExists

	goterr := srv.SignUP(context.Background(), sameUser)

	assert.EqualError(t, goterr, expectedErr.Error())

	//if creating user failed
	newRepoMock.EXPECT().UserExists(gomock.Any(), gomock.AssignableToTypeOf(entity.User{})).DoAndReturn(func(ctx context.Context, u entity.User) (bool, error) {
		if u.Name != user.Name || u.Phone != user.Phone || u.Email != user.Email {
			t.Errorf("expected user details to match, got %v", u)
		}
		return false, nil
	})
	newRepoMock.EXPECT().CreateUser(gomock.Any(), gomock.AssignableToTypeOf(entity.User{})).DoAndReturn(func(ctx context.Context, u entity.User) error {
		if u.Name != user.Name || u.Phone != user.Phone || u.Email != user.Email {
			t.Errorf("expected user details to match, got %v", u)
		}
		return errors.New("database error")
	})
	expectedErr = &errorpac.CustomErr{
		SpecificErr: errorpac.ErrCreateUserFail,
		OriginalErr: errors.New("database error"),
	}
	goterr = srv.SignUP(context.Background(), sameUser)
	t.Log(goterr.Error())
	if expectedErr.Error() != goterr.Error() {
		t.Errorf("expected: %v but got: %v", expectedErr.Error(), goterr.Error())
	}

}

func TestSignIN(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	newRepoMock := mockrepo.NewMockRepository(ctrl)
	newTokenMock := mockrepo.NewMockToken(ctrl)

	srv := New(newRepoMock, newTokenMock)

	// Hash the password for the user returned by the mock
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("kfidrofjkeoirfjkl"), bcrypt.DefaultCost)

	user := entity.User{
		ID:       0,
		Phone:    "+1-245-213-2222",
		Password: string(hashedPassword), // Hashed password for comparison
	}

	userIDstr := strconv.Itoa(user.ID)

	checkUser := entity.User{
		Phone: "+1-245-213-2222",
	}
	// successful sigIN
	// Mock the calls with correct expectations
	newRepoMock.EXPECT().UserPhoneExists(gomock.Any(), gomock.Eq(checkUser)).Return(user, nil)
	newTokenMock.EXPECT().GenerateToken(gomock.Any(), userIDstr, gomock.Any()).Return("token", nil)
	newRepoMock.EXPECT().StoreTokenInRedis(gomock.Any(), userIDstr, "token", gomock.Any()).Return(nil)

	// Input for SignIN
	sameUser := models.UserInfo{
		Name:     "vincent",
		PhoneNO:  "+1-245-213-2222",
		Email:    "test@gamil.com",
		Password: "kfidrofjkeoirfjkl", // Plain text password for input
	}

	// Execute the SignIN function
	_, err := srv.SignIN(context.Background(), sameUser)

	// Check for unexpected errors
	if err != nil {
		t.Errorf("didn't expect any error but got error: %v", err)
	}

	//when user does not exist
	newRepoMock.EXPECT().UserPhoneExists(gomock.Any(), gomock.Eq(checkUser)).Return(entity.User{}, errorpac.ErrUserDoesNotExist)
	expectedErr := errorpac.ErrUserDoesNotExist
	_, err = srv.SignIN(context.Background(), sameUser)

	if expectedErr.Error() != err.Error() {
		t.Errorf("expected: %v but got: %v", expectedErr.Error(), err)
	}

	//when user exists but token generation fails
	newRepoMock.EXPECT().UserPhoneExists(gomock.Any(), gomock.Eq(checkUser)).Return(user, nil)
	newTokenMock.EXPECT().GenerateToken(gomock.Any(), userIDstr, gomock.Any()).Return("", errors.New("token generation error"))

	expectedErr = &errorpac.CustomErr{
		SpecificErr: errorpac.ErrTokenGenFail,
		OriginalErr: errors.New("token generation error"),
	}
	_, err = srv.SignIN(context.Background(), sameUser)

	if expectedErr.Error() != err.Error() {
		t.Errorf("expected: %v but got: %v", expectedErr.Error(), err.Error())
	}

	//when storing the generated token fails
	newRepoMock.EXPECT().UserPhoneExists(gomock.Any(), gomock.Eq(checkUser)).Return(user, nil)
	newTokenMock.EXPECT().GenerateToken(gomock.Any(), userIDstr, gomock.Any()).Return("token", nil)
	newRepoMock.EXPECT().StoreTokenInRedis(gomock.Any(), userIDstr, "token", gomock.Any()).Return(errors.New("failed to store toke in redis"))

	expectedErr = &errorpac.CustomErr{
		SpecificErr: errorpac.ErrFailToStoreToken,
		OriginalErr: errors.New("failed to store toke in redis"),
	}

	_, err = srv.SignIN(context.Background(), sameUser)

	if expectedErr.Error() != err.Error() {
		t.Errorf("expected: %v but got: %v", expectedErr.Error(), err.Error())
	}

}

func TestVerifyToken(t *testing.T) {
	ctrl := gomock.NewController(t)

	newRepoMock := mockrepo.NewMockRepository(ctrl)
	newTokenMock := mockrepo.NewMockToken(ctrl)

	srv := New(newRepoMock, newTokenMock)

	// successfull token verification
	newTokenMock.EXPECT().ParseToken(gomock.Any(), "token").Return("2", nil)
	newTokenMock.EXPECT().ValidateToken(gomock.Any(), "token").Return(nil)

	userID, err := srv.VerifyToken(context.Background(), "token")
	t.Log(userID)
	expectedUserID := "2"
	if expectedUserID != userID {
		t.Errorf("expected: %v but got: %v", expectedUserID, userID)
	}

	if err != nil {
		t.Errorf("expected no error but got err: %v", err)
	}
	// when the parsing of token fails
	newTokenMock.EXPECT().ParseToken(gomock.Any(), "token").Return("", errorpac.ErrTokenParsingFail)

	userID, err = srv.VerifyToken(context.Background(), "token")
	t.Log(userID)
	expectedUserID = ""
	expectedErr := errorpac.ErrTokenParsingFail
	if expectedUserID != userID {
		t.Errorf("expected: %v but got: %v", expectedUserID, userID)
	}

	if expectedErr.Error() != err.Error() {
		t.Errorf("expected: %v but got: %v", expectedErr.Error(), err.Error())
	}
	// when token vaildation fails
	newTokenMock.EXPECT().ParseToken(gomock.Any(), "token").Return("2", nil)
	newTokenMock.EXPECT().ValidateToken(gomock.Any(), "token").Return(errorpac.ErrInvaiidToken)

	_, err = srv.VerifyToken(context.Background(), "token")
	expectedErr = errorpac.ErrInvaiidToken
	if expectedErr.Error() != err.Error() {
		t.Errorf("expected: %v but got: %v", expectedErr.Error(), err.Error())
	}
}

func TestCheckTokenInRedis(t *testing.T) {
	ctrl := gomock.NewController(t)

	newRepoMock := mockrepo.NewMockRepository(ctrl)
	newTokenMock := mockrepo.NewMockToken(ctrl)

	srv := New(newRepoMock, newTokenMock)

	// successful token check
	newTokenMock.EXPECT().ParseToken(gomock.Any(), "token").Return("2", nil)
	newRepoMock.EXPECT().ValidateTokenRedis(gomock.Any(), "token", "2").Return(nil)

	err := srv.CheckTokenInRedis(context.Background(), "token")

	if err != nil {
		t.Errorf("expected no error but got err: %v", err)
	}

	// when the parsing of token fails
	newTokenMock.EXPECT().ParseToken(gomock.Any(), "token").Return("", errorpac.ErrTokenParsingFail)

	err = srv.CheckTokenInRedis(context.Background(), "token")

	expectedErr := errorpac.ErrTokenParsingFail
	if expectedErr.Error() != err.Error() {
		t.Errorf("expected: %v but got: %v", expectedErr.Error(), err.Error())
	}

	// when token validation fails
	newTokenMock.EXPECT().ParseToken(gomock.Any(), "token").Return("2", nil)
	newRepoMock.EXPECT().ValidateTokenRedis(gomock.Any(), "token", "2").Return(errorpac.ErrInvaiidToken)

	err = srv.CheckTokenInRedis(context.Background(), "token")

	expectedErr = errorpac.ErrInvaiidToken
	if expectedErr.Error() != err.Error() {
		t.Errorf("expected: %v but got: %v", expectedErr.Error(), err.Error())
	}

}
