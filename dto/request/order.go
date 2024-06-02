package request

import (
	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type CreateOrder struct {
	Payment Payment `json:"payment" binding:"required"`
	Order   []Order `json:"order" binding:"required"`
}

type Order struct {
	PharmacyId  uint   `json:"pharmacy_id" binding:"gt=0,gte=1,required"`
	CartItemsId []uint `json:"cart_items" binding:"gt=0,gte=1,required"`
	ShipmentId  uint   `json:"shipment_id" binding:"gt=0,gte=1, required" `
}

func (req *CreateOrder) OrderDTO() []entity.Order {
	orders := make([]entity.Order, 0)
	payment := NewPayment(req.Payment)

	for _, order := range req.Order {
		pharmacy := entity.Pharmacy{ID: order.PharmacyId}
		cartItems := make([]*entity.CartItem, 0)
		for _, id := range order.CartItemsId {
			cartItems = append(cartItems, &entity.CartItem{ID: id})
		}
		shipment := entity.ShipmentMethod{ID: order.ShipmentId}
		orders = append(orders,
			entity.Order{Pharmacy: &pharmacy,
				PharmacyId:     pharmacy.ID,
				Cart:           cartItems,
				Payment:        &payment,
				ShipmentMethod: shipment,
			})
	}

	return orders
}
