package response

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type DrugDto struct {
	Id             uint             `json:"id,omitempty"`
	ManufacturerId uint             `json:"manufacturer_id,omitempty"`
	Manufacturer   *ManufacturerDto `json:"manufacturer,omitempty"`
	Name           string           `json:"name"`
	GenericName    string           `json:"generic_name"`
	Composition    string           `json:"composition"`
	Description    string           `json:"description"`
	Classification string           `json:"classification"`
	Form           string           `json:"form"`
	UnitInPack     uint             `json:"unit_in_pack"`
	SellingUnit    string           `json:"selling_unit"`
	Weight         uint             `json:"weight"`
	Height         uint             `json:"height"`
	Length         uint             `json:"length"`
	Width          uint             `json:"width" `
	ImageURL       string           `json:"image_url"`
	IsActive       bool             `json:"is_active"`
	CreatedAt      time.Time        `json:"created_at"`
}
type GetDrugName struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

func NewDrugDto(drug entity.Drug) *DrugDto {
	manufacturer := NewManufacturerDto(drug.Manufacturer)

	return &DrugDto{
		drug.ID,
		drug.Manufacturer.ID,
		manufacturer,
		drug.Name,
		drug.GenericName,
		drug.Composition,
		drug.Description,
		drug.Classification,
		drug.Form,
		drug.UnitInPack,
		drug.SellingUnit,
		drug.Weight,
		drug.Height,
		drug.Length,
		drug.Width,
		drug.ImageURL,
		drug.IsActive,
		drug.CreatedAt,
	}
}

func NewGetDrug(drug entity.Drug) *GetDrugName {
	return &GetDrugName{
		Id:   drug.ID,
		Name: drug.Name,
	}
}
