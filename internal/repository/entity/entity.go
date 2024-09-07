package entity

import "time"

type User struct {
	ID        int
	Name      string
	Phone     string
	Email     string
	Password  string
	Rating    float64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
