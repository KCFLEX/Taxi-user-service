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

type Driver struct {
	ID     int
	Status string // Driver's status ("free", "busy")
}

type TaxiOrder struct {
	UserID    int
	TaxiType  string
	From      string
	To        string
	Status    string
	DriverID  string
	CreatedAt time.Time
}
