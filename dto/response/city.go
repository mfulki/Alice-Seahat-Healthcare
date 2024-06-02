package response

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type CityDto struct {
	ID         uint      `json:"id"`
	ProvinceID uint      `json:"province_id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"created_at"`
}

func NewCityDto(c entity.City) CityDto {
	return CityDto{
		ID:         c.ID,
		ProvinceID: c.ProvinceID,
		Name:       c.Name,
		Type:       c.Type,
		CreatedAt:  c.CreatedAt,
	}
}

func NewMultipleCityDto(cs []entity.City) []CityDto {
	dtos := make([]CityDto, 0)

	for _, c := range cs {
		dtos = append(dtos, NewCityDto(c))
	}

	return dtos
}
