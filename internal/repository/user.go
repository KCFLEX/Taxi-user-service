package repository

import (
	"context"
	"database/sql"

	"github.com/KCFLEX/Taxi-user-service/errorpac"
	"github.com/KCFLEX/Taxi-user-service/internal/config"
	"github.com/KCFLEX/Taxi-user-service/internal/handlers/models"
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

// method to add user to db
func (repo *Repository) CreateUser(ctx context.Context, user models.UserInfo) error {
	query := `INSERT INTO users (name, phone, email, password) VALUES ($1, $2, $3, $4)`

	_, err := repo.db.ExecContext(ctx, query, user.Name, user.PhoneNO, user.Email, user.Password)
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
