package entity

import "database/sql"

type StockRequest struct {
	Id                 uint
	SenderPharmacy     Pharmacy
	ReceiverPharmacy   Pharmacy
	Status             string
	StockRequestDrug   []*StockRequestDrug
	CreatedAt          *sql.NullTime
	UpdatedAt          *sql.NullTime
	DeletedAt          *sql.NullTime
}

type DrugWithPharmacyDrug struct {
	Drug
	SenderPharmacyDrug   PharmacyDrug
	ReceiverPharmacyDrug PharmacyDrug
}
