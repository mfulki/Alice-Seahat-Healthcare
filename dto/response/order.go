package response

import (
	"database/sql"

	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type OrderDTO struct {
	Id             uint              `json:"order_id"`
	Payment        *PaymentDTO       `json:"payment,omitempty"`
	PharmacyId     uint              `json:"pharmacy_id"`
	OrderNumber    string            `json:"order_number"`
	TotalPrice     int               `json:"total_price"`
	FinishedAt     *sql.NullTime     `json:"finished_at,omitempty"`
	Status         string            `json:"status"`
	ShipmentMethod ShipmentMethodDto `json:"shipment_method"`
	Detail         []*OrderDetailDTO `json:"order_detail,omitempty"`
}
type OrderGetDTO struct {
	Id             uint                 `json:"order_id"`
	Payment        *PaymentDTO          `json:"payment,omitempty"`
	PharmacyId     uint                 `json:"pharmacy_id"`
	Pharmacy       *GetPharmacy         `json:"pharmacy"`
	OrderNumber    string               `json:"order_number"`
	TotalPrice     int                  `json:"total_price"`
	FinishedAt     *string              `json:"finished_at"`
	Status         string               `json:"status"`
	ShipmentMethod ShipmentMethodDto    `json:"shipment_method"`
	Detail         []*GetOrderDetailDTO `json:"order_detail,omitempty"`
}

func NewOrderDto(order entity.Order) *OrderDTO {
	payment := NewPaymentDto(order.Payment)
	orderDetails := make([]*OrderDetailDTO, 0)
	for _, orderDetail := range order.Detail {
		orderDetails = append(orderDetails, NewOrderDetailDto(*orderDetail))
	}
	shipment := NewShipmentMethodDto(order.ShipmentMethod)

	return &OrderDTO{
		Id:             order.Id,
		Payment:        payment,
		PharmacyId:     order.PharmacyId,
		OrderNumber:    order.OrderNumber,
		TotalPrice:     order.TotalPrice,
		FinishedAt:     order.FinishedAt,
		Status:         order.Status,
		ShipmentMethod: shipment,
		Detail:         orderDetails,
	}
}
func NewOrderGetDto(order entity.Order) *OrderGetDTO {
	payment := NewPaymentDto(order.Payment)
	orderDetails := make([]*GetOrderDetailDTO, 0)
	for _, orderDetail := range order.Detail {
		orderDetails = append(orderDetails, NewGetOrderDetailDto(*orderDetail))
	}
	pharmacy := NewGetPharmacies(*order.Pharmacy)
	shipment := NewShipmentMethodDto(order.ShipmentMethod)
	var finishedAt *string
	if order.FinishedAt != nil {
		finished := order.FinishedAt.Time.String()
		finishedAt = &finished
	}
	return &OrderGetDTO{
		Id:             order.Id,
		Payment:        payment,
		PharmacyId:     order.PharmacyId,
		Pharmacy:       pharmacy,
		OrderNumber:    order.OrderNumber,
		TotalPrice:     order.TotalPrice,
		FinishedAt:     finishedAt,
		Status:         order.Status,
		ShipmentMethod: shipment,
		Detail:         orderDetails,
	}
}
