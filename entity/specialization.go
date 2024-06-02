package entity

import "time"

type Specialization struct {
	ID         uint
	Name       string
	CreatedAt  time.Time
	UpdatedAt time.Time
	DeletedAt  *time.Time
}
