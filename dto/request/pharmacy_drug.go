package request

import (
	"strconv"

	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type AddPharmacyDrug struct {
	DrugID     uint  `json:"drug_id" binding:"required"`
	PharmacyID uint  `json:"pharmacy_id" binding:"required"`
	CategoryID uint  `json:"category_id" binding:"required"`
	Stock      uint  `json:"stock" binding:"required"`
	Price      uint  `json:"price" binding:"required"`
	IsActive   *bool `json:"is_active" binding:"required"`
}

type PharmaryDrugsQuery struct {
	Latitude  string `form:"lat"`
	Longitude string `form:"long"`
}

func (req *AddPharmacyDrug) PharmacyDrug() entity.PharmacyDrug {
	return entity.PharmacyDrug{
		DrugID:     req.DrugID,
		PharmacyID: req.PharmacyID,
		CategoryID: req.CategoryID,
		Stock:      req.Stock,
		Price:      req.Price,
		IsActive:   *req.IsActive,
	}
}

func (q *PharmaryDrugsQuery) Address() entity.Address {
	parsedLongitude, _ := strconv.ParseFloat(q.Longitude, 64)
	parsedLatitude, _ := strconv.ParseFloat(q.Latitude, 64)

	return entity.Address{
		Latitude:  parsedLatitude,
		Longitude: parsedLongitude,
	}
}
