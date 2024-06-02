package request

import "Alice-Seahat-Healthcare/seahat-be/entity"

type AddCartItemRequest struct {
	PharmacyDrugID int   `json:"pharmacy_drug_id" binding:"required,gte=1"`
	Quantity       int   `json:"quantity" binding:"required,gte=1"`
	IsPrescripted  *bool `json:"is_prescripted" binding:"required,boolean"`
}

type UpdateQtyItemRequest struct {
	Quantity int `json:"quantity" binding:"required,gte=1"`
}

type DeleteCartItemRequest struct {
	IDs []int `json:"ids" binding:"gt=0,dive,required,gte=1"`
}

type CartItem struct {
	Id uint `json:"cart_items_id" binding:"gt=0,gte=1,required"`
}

func (req AddCartItemRequest) CartItem() entity.CartItem {
	return entity.CartItem{
		PharmacyDrugID: uint(req.PharmacyDrugID),
		Quantity:       uint(req.Quantity),
		IsPrescripted:  *req.IsPrescripted,
	}
}

func (req UpdateQtyItemRequest) CartItem(cartItemID int) entity.CartItem {
	return entity.CartItem{
		ID:       uint(cartItemID),
		Quantity: uint(req.Quantity),
	}
}

func (req DeleteCartItemRequest) Uint() []uint {
	ids := make([]uint, 0)

	for _, id := range req.IDs {
		ids = append(ids, uint(id))
	}

	return ids
}
