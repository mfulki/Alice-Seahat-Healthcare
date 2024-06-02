package response

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type SubdistrictDto struct {
	ID        uint      `json:"id"`
	CityID    uint      `json:"city_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func NewSubdistrictDto(sd entity.Subdistrict) SubdistrictDto {
	return SubdistrictDto{
		ID:        sd.ID,
		CityID:    sd.CityID,
		Name:      sd.Name,
		CreatedAt: sd.CreatedAt,
	}
}

func NewMultipleSubdistrictDto(sds []entity.Subdistrict) []SubdistrictDto {
	dtos := make([]SubdistrictDto, 0)

	for _, sd := range sds {
		dtos = append(dtos, NewSubdistrictDto(sd))
	}

	return dtos
}
