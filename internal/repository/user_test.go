package repository

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KCFLEX/Taxi-user-service/errorpac"
	"github.com/KCFLEX/Taxi-user-service/internal/repository/entity"
	"github.com/go-redis/redismock/v9"
	"github.com/lib/pq"
)

func TestUserExists(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock: ", err)
	}
	defer db.Close()

	//	redisClient, mockClient := redismock.NewClientMock()

	repo := &Repository{
		db: db,
	}
	user := entity.User{
		Email: "test@example.com",
		Phone: "+1234567890",
	}

	// when user exists
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 OR phone = $2)`
	mock.ExpectQuery(query).WithArgs("test@example.com", "+1234567890").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	expectedBool := true

	got, err := repo.UserExists(context.Background(), user)
	if err != nil {
		t.Errorf("expected no error but got error: %v", err)
	}

	if got != expectedBool {
		t.Errorf("User should exist but the result shows false: %v", got)
	}

	// when user does not exist
	mock.ExpectQuery(query).WithArgs("test@example.com", "+1234567890").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
	expectedBool = false

	got, err = repo.UserExists(context.Background(), user)
	if err != nil {
		t.Errorf("expected no error but got error: %v", err)
	}

	if got != expectedBool {
		t.Errorf("User should not exist but the result shows %v meaning the user exists", got)
	}

	// when database errors occurs
	mock.ExpectQuery(query).WithArgs("test@example.com", "+1234567890").WillReturnError(sqlmock.ErrCancelled)
	expectedBool = false
	expectedErr := &errorpac.CustomErr{
		OriginalErr: sqlmock.ErrCancelled,
		SpecificErr: errors.New("failed to check if user exists"),
	}
	got, err = repo.UserExists(context.Background(), user)
	if expectedErr.Error() != err.Error() {
		t.Errorf("expected: %v but got: %v", expectedErr.Error(), err.Error())
	}

	if got != expectedBool {
		t.Errorf("expected: %v but got: %v", expectedBool, got)
	}

}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatal("error init mock: ", err)
	}
	defer db.Close()
	repo := &Repository{
		db: db,
	}
	user := entity.User{
		Name:     "vincent",
		Phone:    "+1-245-213-2222",
		Email:    "test@gamil.com",
		Password: "kfidrofjkeoirfjkl",
	}
	// successfully created user without errors
	query := `INSERT INTO users (name, phone, email, password) VALUES ($1, $2, $3, $4)`
	mock.ExpectExec(query).WithArgs("vincent", "+1-245-213-2222", "test@gamil.com", "kfidrofjkeoirfjkl").WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateUser(context.Background(), user)

	if err != nil {
		t.Errorf("creation of user failed got err: %v", err)
	}

	// database error: duplicate email (unique constraint violation)
	mock.ExpectExec(query).WithArgs("vincent", "+1-245-213-2222", "test@gamil.com", "kfidrofjkeoirfjkl").WillReturnError(&pq.Error{Code: "23505"})
	expectedErr := &errorpac.CustomErr{
		OriginalErr: &pq.Error{Code: "23505"},
		SpecificErr: errorpac.ErrDuplicateEmail,
	}
	err = repo.CreateUser(context.Background(), user)

	if expectedErr.Error() != err.Error() {
		t.Errorf("expected: %v but got: %v", expectedErr.Error(), err.Error())
	}

	// database error: failed to create user

	mock.ExpectExec(query).WithArgs("vincent", "+1-245-213-2222", "test@gamil.com", "kfidrofjkeoirfjkl").WillReturnError(errors.New("user creation failed"))
	expectedErr = &errorpac.CustomErr{
		OriginalErr: errors.New("user creation failed"),
		SpecificErr: errorpac.ErrCreateUserFail,
	}
	err = repo.CreateUser(context.Background(), user)

	if expectedErr.Error() != err.Error() {
		t.Errorf("expected: %v but got: %v", expectedErr.Error(), err.Error())
	}
}

func TestStoreTokenInRedis(t *testing.T) {
	redisClient, mock := redismock.NewClientMock()

	repo := &Repository{
		redisDB: redisClient,
	}

	// Step 2: Define test variables
	ctx := context.Background()
	userID := "user123"
	token := "someToken"
	expiration := time.Hour * 2

	// succesfull token storage
	mock.ExpectSet("auth:user123", "someToken", 2*time.Hour).SetVal("OK")

	err := repo.StoreTokenInRedis(ctx, userID, token, expiration)

	if err != nil {
		t.Errorf("Token insertion failed got err: %v", err)
	}

	// Ensure all Redis expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
	// when token storing fails
	mock.ExpectSet("auth:user123", "someToken", 2*time.Hour).SetErr(errors.New("failed to store token in redis"))
	expectedErr := "failed to store token in redis"
	err = repo.StoreTokenInRedis(ctx, userID, token, expiration)

	if expectedErr != err.Error() {
		t.Errorf("expected: %v but got: %v", expectedErr, err.Error())
	}

}

func TestValidateTokenRedis(t *testing.T) {
	redisClient, mock := redismock.NewClientMock()

	repo := &Repository{
		redisDB: redisClient,
	}
	userID := "user123"
	ctx := context.Background()
	validToken := "validToken"
	invalidToken := "invalidToken"
	// successfull token validation
	mock.ExpectGet("auth:user123").SetVal(validToken)

	err := repo.ValidateTokenRedis(ctx, validToken, userID)

	if err != nil {
		t.Errorf("token validation failed got err: %v", err)
	}

	// invalid token
	mock.ExpectGet("auth:user123").SetVal(invalidToken)
	expectedErr := errorpac.ErrInvaiidToken

	err = repo.ValidateTokenRedis(ctx, validToken, userID)
	if expectedErr != err {
		t.Errorf("expected: %v but got: %v", expectedErr, err)
	}

	//Token not found in redis
	mock.ExpectGet("auth:user123").RedisNil()
	expectedErr = errors.New("token not found in Redis")
	err = repo.ValidateTokenRedis(ctx, validToken, userID)
	if expectedErr.Error() != err.Error() {
		t.Errorf("expected: %v but got: %v", expectedErr, err)
	}
}
