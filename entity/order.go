package entity

import (
	"database/sql"
)

type Order struct {
	Id             uint
	Payment        *Payment
	PharmacyId     uint
	Pharmacy       *Pharmacy
	OrderNumber    string
	TotalPrice     int
	FinishedAt     *sql.NullTime
	Status         string
	ShipmentMethod ShipmentMethod
	Cart           []*CartItem
	Detail         []*OrderDetail
	CreatedAt      *sql.NullTime
}
