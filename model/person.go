package model

import "time"

type LoginShowPerson struct {
	CustomerID  int
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	Address     string
	Password    string `json:"-"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
