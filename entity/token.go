package entity

import "time"

type Token struct {
	ID        uint
	UserID    *uint
	DoctorID  *uint
	Type      string
	Token     string
	ExpiredAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
