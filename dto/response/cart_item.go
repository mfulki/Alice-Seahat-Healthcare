package response

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type CartItemDto struct {
	ID             uint      `json:"id"`
	UserID         uint      `json:"user_id"`
	PharmacyDrugID uint      `json:"pharmacy_drug_id"`
	Quantity       uint      `json:"quantity"`
	IsPrescripted  bool      `json:"is_prescripted"`
	Price          uint      `json:"price,omitempty"`
	Stock          *uint     `json:"stock,omitempty"`
	Drug           *DrugDto  `json:"drug"`
	CreatedAt      time.Time `json:"created_at"`
}

type CartItemsWithPharmacyDto struct {
	PharmacyID   uint          `json:"pharmacy_id"`
	PharmacyName string        `json:"pharmacy_name"`
	Items        []CartItemDto `json:"items"`
}

type CartItemsWithPharmacyDtos []CartItemsWithPharmacyDto

func NewCartItemDto(item entity.CartItem) CartItemDto {
	return CartItemDto{
		ID:             item.ID,
		UserID:         item.UserID,
		PharmacyDrugID: item.PharmacyDrugID,
		Quantity:       item.Quantity,
		IsPrescripted:  item.IsPrescripted,
		CreatedAt:      item.CreatedAt,
	}
}

func NewCartItemsWithSummary(items []entity.CartItem) CartItemsWithPharmacyDtos {
	total := uint(0)
	dtos := make(CartItemsWithPharmacyDtos, 0)

	for _, item := range items {
		total += item.Quantity * uint(item.PharmacyDrug.Price)
		cartItemDto := CartItemDto{
			ID:             item.ID,
			UserID:         item.UserID,
			PharmacyDrugID: item.PharmacyDrugID,
			Quantity:       item.Quantity,
			Price:          uint(item.PharmacyDrug.Price),
			Stock:          &item.PharmacyDrug.Stock,
			IsPrescripted:  item.IsPrescripted,
			Drug:           NewDrugDto(item.PharmacyDrug.Drug),
			CreatedAt:      item.CreatedAt,
		}

		dtoIndex := dtos.ContainPharmacyID(item.PharmacyDrug.Pharmacy.ID)
		if dtoIndex == -1 {
			dtos = append(dtos, CartItemsWithPharmacyDto{
				PharmacyID:   item.PharmacyDrug.Pharmacy.ID,
				PharmacyName: item.PharmacyDrug.Pharmacy.Name,
				Items:        []CartItemDto{cartItemDto},
			})

			continue
		}

		dtos[dtoIndex].Items = append(dtos[dtoIndex].Items, cartItemDto)
	}

	return dtos
}

func (dtos CartItemsWithPharmacyDtos) ContainPharmacyID(id uint) int {
	for index, dto := range dtos {
		if dto.PharmacyID == id {
			return index
		}
	}

	return -1
}
