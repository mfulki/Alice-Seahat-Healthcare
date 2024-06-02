package entity

type AdminReportByDrug struct {
	DrugID           uint
	DrugName         string
	ManufacturerName string
	Classification   string
	Form             string
	UnitInPack       uint
	SellingUnit      string
	ImageURL         string
	TotalQuantiy     uint
	TotalPrice       uint
}

type AdminReportByCategory struct {
	CategoryID   uint
	CategoryName string
	TotalQuantiy uint
	TotalPrice   uint
}
