package entity

import "database/sql"

type StockRequestDrug struct {
	Id             uint
	StockRequestId uint
	DrugId         uint
	Drug           Drug
	Quantity       int
	CreatedAt      *sql.NullTime
	UpdatedAt      *sql.NullTime
	DeletedAt      *sql.NullTime
}
