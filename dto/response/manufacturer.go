package response

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type ManufacturerDto struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func NewManufacturerDto(manufacturer entity.Manufacturer) *ManufacturerDto {
	if manufacturer.Name == "" {
		return nil
	}

	return &ManufacturerDto{
		ID:        manufacturer.ID,
		Name:      manufacturer.Name,
		CreatedAt: manufacturer.CreatedAt,
	}
}

func NewMultipleManufacturerDto(mfs []entity.Manufacturer) []ManufacturerDto {
	dtos := make([]ManufacturerDto, 0)

	for _, mf := range mfs {
		dto := NewManufacturerDto(mf)
		if dto == nil {
			continue
		}

		dtos = append(dtos, *dto)
	}

	return dtos
}
