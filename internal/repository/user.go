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

// method to check if user email exists
// if it does return password
func (repo *Repository) UserPhoneExists(ctx context.Context, user entity.User) (entity.User, error) {
	query := `SELECT password, id FROM users WHERE phone = $1`

	var id int
	var password string

	err := repo.db.QueryRowContext(ctx, query, user.Phone).Scan(&password, &id)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, &errorpac.CustomErr{
				OriginalErr: err,
				SpecificErr: errorpac.ErrUserDoesNotExist,
			}
		}
	}

	return entity.User{
		ID:       id,
		Password: password,
	}, nil
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
