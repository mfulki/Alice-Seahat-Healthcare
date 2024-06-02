package response

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type AddressDto struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	SubdistrictID uint      `json:"subdistrict_id"`
	Address       string    `json:"address"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	IsMain        bool      `json:"is_main"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

func NewAddressDto(addr entity.Address) AddressDto {
	return AddressDto{
		ID:            addr.ID,
		UserID:        addr.UserID,
		SubdistrictID: addr.SubdistrictID,
		Address:       addr.Address,
		Latitude:      addr.Latitude,
		Longitude:     addr.Longitude,
		IsMain:        addr.IsMain,
		IsActive:      addr.IsActive,
		CreatedAt:     addr.CreatedAt,
	}
}

func NewMultipleAddressDto(addrs []entity.Address) []AddressDto {
	dtos := make([]AddressDto, 0)

	for _, addr := range addrs {
		dto := NewAddressDto(addr)

		if addr.IsMain {
			dtos = append([]AddressDto{dto}, dtos...)
			continue
		}

		dtos = append(dtos, dto)
	}

	return dtos
}
