package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/KCFLEX/Taxi-user-service/errorpac"
	"github.com/KCFLEX/Taxi-user-service/internal/config"
	"github.com/KCFLEX/Taxi-user-service/internal/repository/entity"
	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
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
	return &Repository{
		db: db,
	}, nil
}

func (repo *Repository) Close() error {
	return repo.db.Close()
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
