package response

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type ProvinceDto struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func NewProvinceDto(c entity.Province) ProvinceDto {
	return ProvinceDto{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
	}
}

func NewMultipleProvinceDto(cs []entity.Province) []ProvinceDto {
	dtos := make([]ProvinceDto, 0)

	for _, c := range cs {
		dtos = append(dtos, NewProvinceDto(c))
	}

	return dtos
}
