package entity

import "time"

type CartItem struct {
	ID             uint
	UserID         uint
	PharmacyDrugID uint
	Quantity       uint
	IsPrescripted  bool
	Price          uint
	TotalPrice     uint
	PharmacyDrug   PharmacyDrug
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}
