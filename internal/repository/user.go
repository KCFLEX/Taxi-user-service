package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/KCFLEX/Taxi-user-service/errorpac"
	"github.com/KCFLEX/Taxi-user-service/internal/config"
	"github.com/KCFLEX/Taxi-user-service/internal/repository/entity"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	db      *sql.DB
	redisDB *redis.Client
}

func RedisConn(conn string) (*redis.Client, error) {
	db, err := redis.ParseURL(conn)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(db)

	// Test the Redis connection
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.New("failed to connect to redis ")
	}

	return client, nil
}

func DbConnect(conn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func New(config config.Config) (*Repository, error) {
	db, err := DbConnect(config.DbConn)
	if err != nil {
		return nil, err
	}

	redisDB, err := RedisConn(config.RedisConn)
	if err != nil {
		return nil, err
	}
	return &Repository{
		db:      db,
		redisDB: redisDB,
	}, nil

}

func (repo *Repository) Close() error {
	return repo.db.Close()
}

func (repo *Repository) CloseRedis() error {
	return repo.redisDB.Close()
}

// method to check if user exists already
func (repo *Repository) UserExists(ctx context.Context, user entity.User) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 OR phone = $2)`

	var exists bool

	err := repo.db.QueryRowContext(ctx, query, user.Email, user.Phone).Scan(&exists)
	if err != nil {
		return false, &errorpac.CustomErr{
			OriginalErr: err,
			SpecificErr: errors.New("failed to check if user exists"),
		}
	}

	return exists, nil
}

// method to add user to db
func (repo *Repository) CreateUser(ctx context.Context, user entity.User) error {
	query := `INSERT INTO users (name, phone, email, password) VALUES ($1, $2, $3, $4)`

	_, err := repo.db.ExecContext(ctx, query, user.Name, user.Phone, user.Email, user.Password)
	if err != nil {
		// checking if the error is due to unqiue constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" { // Unique violation error code
			return &errorpac.CustomErr{
				OriginalErr: err,
				SpecificErr: errorpac.ErrDuplicateEmail,
			}
		}
		return &errorpac.CustomErr{
			OriginalErr: err,
			SpecificErr: errorpac.ErrCreateUserFail,
		}
	}
	return nil
}

// method to check if user phone number exists
// if it does return password and id
func (repo *Repository) UserPhoneExists(ctx context.Context, user entity.User) (entity.User, error) {
	query := `SELECT password, id, deleted_at FROM users WHERE phone = $1`

	var UserExists entity.User

	err := repo.db.QueryRowContext(ctx, query, user.Phone).Scan(&UserExists.Password, &UserExists.ID, &UserExists.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, &errorpac.CustomErr{
				OriginalErr: err,
				SpecificErr: errorpac.ErrUserDoesNotExist,
			}
		}
	}

	return UserExists, nil
}

func (repo *Repository) StoreTokenInRedis(ctx context.Context, userID string, token string, expiration time.Duration) error {
	err := repo.redisDB.Set(ctx, "auth:"+userID, token, expiration).Err()

	if err != nil {
		return errors.New("failed to store token in redis")
	}

	return nil
}

func (repo *Repository) ValidateTokenRedis(ctx context.Context, token string, userID string) error {
	storedToken, err := repo.redisDB.Get(ctx, "auth:"+userID).Result()
	if err != nil {
		if err == redis.Nil {
			return errors.New("token not found in Redis")
		}
		return err
	}

	if storedToken != token {
		return errorpac.ErrInvaiidToken
	}
	return nil
}

func (repo *Repository) GetProfileByID(ctx context.Context, id int) (entity.User, error) {
	query := `SELECT name, phone, email, rating  FROM users WHERE id = $1`

	var userProfile entity.User

	err := repo.db.QueryRowContext(ctx, query, id).Scan(&userProfile.Name, &userProfile.Phone, &userProfile.Email, userProfile.Rating)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, &errorpac.CustomErr{
				OriginalErr: err,
				SpecificErr: errorpac.ErrUserDoesNotExist,
			}
		}
	}

	return userProfile, nil
}

func (repo *Repository) DeleteProfileByID(ctx context.Context, id int) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = $1`

	result, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &errorpac.CustomErr{
			OriginalErr: err,
			SpecificErr: errorpac.ErrUserDoesNotExist,
		}
	}

	return nil

}

func (repo *Repository) UpdateProfileByID(ctx context.Context, updateInfo entity.User) error {
	query := `UPDATE users SET name = $1, email = $2, phone = $3 WHERE id = $4`

	result, err := repo.db.ExecContext(ctx, query, updateInfo.Name, updateInfo.Email, updateInfo.Phone, updateInfo.ID)
	if err != nil {
		return &errorpac.CustomErr{
			OriginalErr: err,
			SpecificErr: errors.New("failed to update query"),
		}
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &errorpac.CustomErr{
			OriginalErr: err,
			SpecificErr: errorpac.ErrUserDoesNotExist,
		}
	}

	return nil

}
