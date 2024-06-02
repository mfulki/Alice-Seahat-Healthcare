package response

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type PartnerDto struct {
	ID                uint               `json:"id"`
	PharmacyManagerID uint               `json:"pharmacy_manager_id"`
	Name              string             `json:"name"`
	Logo              string             `json:"logo"`
	IsActive          bool               `json:"is_active"`
	CreatedAt         time.Time          `json:"created_at"`
	PharmacyManager   PharmacyManagerDto `json:"pharmacy_manager,omitempty"`
}

func NewPartnerDto(p entity.Partner) PartnerDto {
	return PartnerDto{
		ID:                p.ID,
		PharmacyManagerID: p.PharmacyManagerID,
		Name:              p.Name,
		Logo:              p.Logo,
		IsActive:          p.IsActive,
		CreatedAt:         p.CreatedAt,
		PharmacyManager:   NewPharmacyManagerDto(p.PharmacyManager),
	}
}

func NewMultiplePartnerDto(ps []entity.Partner) []PartnerDto {
	dtos := make([]PartnerDto, 0)

	for _, p := range ps {
		dtos = append(dtos, NewPartnerDto(p))
	}

	return dtos
}
