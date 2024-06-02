package entity

import "time"

type Prescription struct {
	ID             uint
	TelemedicineID uint
	DrugID         uint
	Quantity       uint
	Notes          string
	Drug           Drug
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}
