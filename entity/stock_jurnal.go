package entity

import (
	"database/sql"
)

type StockJurnal struct {
	Id          uint
	DrugId      uint
	DrugName    string
	PharmacyId  uint
	Quantity    int
	Description string
	CreatedAt   *sql.NullTime
	UpdatedAt   *sql.NullTime
	DeletedAt   *sql.NullTime
}
