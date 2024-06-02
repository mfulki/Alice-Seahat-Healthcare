package entity

import "time"

type Partner struct {
	ID                uint
	PharmacyManagerID uint
	PharmacyManager   PharmacyManager
	Name              string
	Logo              string
	IsActive          bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}
