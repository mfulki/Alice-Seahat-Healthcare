package entity

import (
	"time"
)

type Subdistrict struct {
	ID        uint
	CityID    uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
