package response

import (
	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type AdminDrugReportDTO struct {
	DrugID           uint   `json:"drug_id"`
	DrugName         string `json:"drug_name"`
	ManufacturerName string `json:"manufacturer_name"`
	Classification   string `json:"classification"`
	Form             string `json:"form"`
	UnitInPack       uint   `json:"unit_in_pack"`
	SellingUnit      string `json:"selling_unit"`
	ImageURL         string `json:"image_url"`
	TotalQuantiy     uint   `json:"total_quantity"`
	TotalPrice       uint   `json:"total_price"`
}

type AdminCategoryReportDTO struct {
	CategoryID   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
	TotalQuantiy uint   `json:"total_quantity"`
	TotalPrice   uint   `json:"total_price"`
}

func NewAdminDrugReportDTO(p entity.AdminReportByDrug) AdminDrugReportDTO {
	return AdminDrugReportDTO{
		DrugID:           p.DrugID,
		DrugName:         p.DrugName,
		ManufacturerName: p.ManufacturerName,
		Classification:   p.Classification,
		Form:             p.Form,
		UnitInPack:       p.UnitInPack,
		SellingUnit:      p.SellingUnit,
		ImageURL:         p.ImageURL,
		TotalQuantiy:     p.TotalQuantiy,
		TotalPrice:       p.TotalPrice,
	}
}

func NewAdminCategoryReportDTO(p entity.AdminReportByCategory) AdminCategoryReportDTO {
	return AdminCategoryReportDTO{
		CategoryID:   p.CategoryID,
		CategoryName: p.CategoryName,
		TotalQuantiy: p.TotalQuantiy,
		TotalPrice:   p.TotalPrice,
	}
}
