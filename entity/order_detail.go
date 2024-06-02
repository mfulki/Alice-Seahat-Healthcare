package entity

type OrderDetail struct {
	Id             uint
	OrderId        uint
	PharmacyDrugId uint
	PharmacyDrug   PharmacyDrug
	Quantity       uint
	Price          uint
}
