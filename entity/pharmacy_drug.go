package entity

import "time"

type PharmacyDrug struct {
	ID         uint
	DrugID     uint
	PharmacyID uint
	CategoryID uint
	Stock      uint
	Price      uint
	IsActive   bool
	Drug       Drug
	Pharmacy   Pharmacy
	Category   Category
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
