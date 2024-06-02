package request

import (
	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type DrugDTO struct {
	ManufacturerId uint   `json:"manufacturer_id" binding:"required,gte=1"`
	Name           string `json:"name" binding:"required,gte=2"`
	GenericName    string `json:"generic_name" binding:"required,gte=2"`
	Composition    string `json:"composition" binding:"required,gte=8"`
	Description    string `json:"description" binding:"required,gte=15"`
	Classification string `json:"classification" binding:"required,oneof='obat bebas' 'obat keras' 'obat bebas terbatas' 'non obat'"`
	Form           string `json:"form" binding:"required,oneof=pill capsule"`
	UnitInPack     uint   `json:"unit_in_pack" binding:"required,gte=1"`
	SellingUnit    string `json:"selling_unit" binding:"required,oneof=botol strip pack box tube piece pouch sachet kaleng pot tablet unit kapsul paket suppositoria bag pen vial lembar ovula ampul test"`
	Weight         uint   `json:"weight" binding:"required,gte=1"`
	Height         uint   `json:"height" binding:"required,gte=1"`
	Length         uint   `json:"length" binding:"required,gte=1"`
	Width          uint   `json:"width"  binding:"required,gte=1"`
	ImageURL       string `json:"image_url" binding:"required,url"`
	IsActive       *bool  `json:"is_active" binding:"required,boolean"`
}

func (req *DrugDTO) ConvertIntoEntityDrugs() *entity.Drug {
	manufacturer := entity.Manufacturer{ID: req.ManufacturerId}
	return &entity.Drug{
		Manufacturer:   manufacturer,
		Name:           req.Name,
		GenericName:    req.GenericName,
		Composition:    req.Composition,
		Description:    req.Description,
		Classification: req.Classification,
		Form:           req.Form,
		UnitInPack:     req.UnitInPack,
		SellingUnit:    req.SellingUnit,
		Weight:         req.Weight,
		Height:         req.Height,
		Length:         req.Length,
		Width:          req.Width,
		ImageURL:       req.ImageURL,
		IsActive:       *req.IsActive,
	}

}

func (req *DrugDTO) ConvertMultipleReqIntoEntityDrugs(reqUri GetIdUri) *entity.Drug {
	manufacturer := entity.Manufacturer{ID: req.ManufacturerId}
	return &entity.Drug{
		ID:             reqUri.DrugId,
		Manufacturer:   manufacturer,
		Name:           req.Name,
		GenericName:    req.GenericName,
		Composition:    req.Composition,
		Description:    req.Description,
		Classification: req.Classification,
		Form:           req.Form,
		UnitInPack:     req.UnitInPack,
		SellingUnit:    req.SellingUnit,
		Weight:         req.Weight,
		Height:         req.Height,
		Length:         req.Length,
		Width:          req.Width,
		ImageURL:       req.ImageURL,
		IsActive:       *req.IsActive,
	}
}

type GetIdUri struct {
	DrugId uint `uri:"id" binding:"required"`
}

func (req *GetIdUri) Drugs() entity.Drug {
	return entity.Drug{
		ID: req.DrugId,
	}
}
