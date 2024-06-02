package entity

import (
	"database/sql"
	"time"
)

type ShipmentMethod struct {
	ID          uint
	Name        string
	CourierName string
	Price       *uint
	Duration    uint
	CreatedAt   *sql.NullTime
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
