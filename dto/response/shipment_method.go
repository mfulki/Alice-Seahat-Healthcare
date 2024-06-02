package response

import (
	"database/sql"

	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type ShipmentMethodDto struct {
	ID          uint          `json:"id,omitempty"`
	Name        string        `json:"name"`
	CourierName string        `json:"courier_name,omitempty"`
	Price       *uint         `json:"price"`
	Duration    uint          `json:"duration,omitempty"`
	CreatedAt   *sql.NullTime `json:"created_at,omitempty"`
}

func NewShipmentMethodDto(sm entity.ShipmentMethod) ShipmentMethodDto {
	return ShipmentMethodDto{
		ID:          sm.ID,
		Name:        sm.Name,
		CourierName: sm.CourierName,
		Price:       sm.Price,
		Duration:    sm.Duration,
		CreatedAt:   sm.CreatedAt,
	}
}

func NewMultipleShipmentMethodDto(sms []*entity.ShipmentMethod) []ShipmentMethodDto {
	dtos := make([]ShipmentMethodDto, 0)

	for _, sm := range sms {
		dtos = append(dtos, NewShipmentMethodDto(*sm))
	}

	return dtos
}
