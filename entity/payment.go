package entity

import (
	"database/sql"
	"time"
)

type Payment struct {
	Id              uint
	UserId          uint
	UserName        string
	Method          string
	Proof           *string
	FullUserAddress string
	Address         *Address
	TotalPrice      int
	Number          string
	Status          string
	ExpiredAt       *sql.NullTime
	Orders          []*Order
	CreatedAt       *sql.NullTime
	UpdatedAt       time.Time
	DeletedAt       *sql.NullTime
}
