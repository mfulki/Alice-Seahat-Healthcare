package response

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type PharmacyDrugDto struct {
	ID        uint         `json:"id"`
	DrugID    uint         `json:"drug_id,omitempty"`
	Drug      *DrugDto     `json:"drug,omitempty"`
	Pharmacy  *PharmacyDto `json:"pharmacy"`
	Category  *CategoryDto `json:"category"`
	Stock     uint         `json:"stock"`
	Price     uint         `json:"price"`
	IsActive  bool         `json:"is_active"`
	CreatedAt time.Time    `json:"created_at,omitempty"`
}

type GetPharmacyDrugAndDrugDto struct {
	ID         uint         `json:"id,omitempty"`
	DrugID     uint         `json:"drug_id,omitempty"`
	CategoryID uint         `json:"category_id,omitempty"`
	PharmacyID uint         `json:"pharmacy_id,omitempty"`
	Drug       DrugDto      `json:"drug,omitempty"`
	Category   *CategoryDto `json:"category,omitempty"`
	Stock      uint         `json:"stock,omitempty"`
	Price      uint         `json:"price,omitempty"`
	IsActive   *bool        `json:"is_active,omitempty"`
	CreatedAt  time.Time    `json:"created_at,omitempty"`
}

type CreatePharmacyDrugDto struct {
	ID        uint      `json:"id"`
	DrugID    uint      `json:"drug_id"`
	Pharmacy  uint      `json:"pharmacy_id"`
	Stock     uint      `json:"stock"`
	Price     uint      `json:"price"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func NewPharmacyDrugDto(pharmacyDrug entity.PharmacyDrug) (*PharmacyDrugDto, error) {
	var drugDto *DrugDto
	if pharmacyDrug.Drug.Name != "" {
		drugDto = NewDrugDto(pharmacyDrug.Drug)
	}

	var categoryDto *CategoryDto
	if pharmacyDrug.Category.Name != "" {
		dto := NewCategoryDto(pharmacyDrug.Category)
		categoryDto = &dto
	}

	var pharmacyDto *PharmacyDto
	if pharmacyDrug.Pharmacy.Name != "" {
		pharmacy, err := NewPharmacyDto(pharmacyDrug.Pharmacy)
		if err != nil {
			return nil, err
		}

		pharmacyDto = pharmacy
	}

	return &PharmacyDrugDto{
		ID:        pharmacyDrug.ID,
		DrugID:    pharmacyDrug.DrugID,
		Drug:      drugDto,
		Pharmacy:  pharmacyDto,
		Stock:     pharmacyDrug.Stock,
		Price:     pharmacyDrug.Price,
		IsActive:  pharmacyDrug.IsActive,
		CreatedAt: pharmacyDrug.CreatedAt,
		Category:  categoryDto,
	}, nil
}

func NewGetPharmacyDrugJoinDrugsDto(pharmacyDrug entity.PharmacyDrug) (*GetPharmacyDrugAndDrugDto, error) {
	drug := NewDrugDto(pharmacyDrug.Drug)

	var category *CategoryDto
	if pharmacyDrug.Category.Name != "" {
		dto := NewCategoryDto(pharmacyDrug.Category)
		category = &dto
	}

	return &GetPharmacyDrugAndDrugDto{
		ID:         pharmacyDrug.ID,
		DrugID:     pharmacyDrug.DrugID,
		CategoryID: pharmacyDrug.CategoryID,
		PharmacyID: pharmacyDrug.PharmacyID,
		Drug:       *drug,
		Category:   category,
		Stock:      pharmacyDrug.Stock,
		Price:      pharmacyDrug.Price,
		IsActive:   &pharmacyDrug.IsActive,
		CreatedAt:  pharmacyDrug.CreatedAt,
	}, nil
}

func NewCreatePharmacyDrugDto(pharmacyDrug entity.PharmacyDrug) *CreatePharmacyDrugDto {
	return &CreatePharmacyDrugDto{
		ID:        pharmacyDrug.ID,
		DrugID:    pharmacyDrug.DrugID,
		Pharmacy:  pharmacyDrug.PharmacyID,
		Stock:     pharmacyDrug.Stock,
		Price:     pharmacyDrug.Price,
		IsActive:  pharmacyDrug.IsActive,
		CreatedAt: pharmacyDrug.CreatedAt,
	}
}
