package response

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type PharmacyManagerDto struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func NewPharmacyManagerDto(pm entity.PharmacyManager) PharmacyManagerDto {
	return PharmacyManagerDto{
		ID:        pm.ID,
		Name:      pm.Name,
		Email:     pm.Email,
		CreatedAt: pm.CreatedAt,
	}
}

func NewMultiplePharmacyManagerDto(pms []entity.PharmacyManager) []PharmacyManagerDto {
	dtos := make([]PharmacyManagerDto, 0)

	for _, pm := range pms {
		dtos = append(dtos, NewPharmacyManagerDto(pm))
	}

	return dtos
}
