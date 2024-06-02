package entity

import (
	"time"
)

type City struct {
	ID         uint
	ProvinceID uint
	Name       string
	Type       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
