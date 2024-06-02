package response

import "Alice-Seahat-Healthcare/seahat-be/entity"

type OrderDetailDTO struct {
	Id             uint `json:"order_detail_id"`
	OrderId        uint `json:"order_id"`
	PharmacyDrugId uint `json:"pharmacy_drug_id"`
	Quantity       uint `json:"quantity"`
	Price          uint `json:"price"`
}
type GetOrderDetailDTO struct {
	Id             uint                      `json:"order_detail_id"`
	OrderId        uint                      `json:"order_id,omitempty"`
	PharmacyDrugId uint                      `json:"pharmacy_drug_id,omitempty"`
	PharmacyDrug   GetPharmacyDrugAndDrugDto `json:"pharmacy_drug,omitempty"`
	Quantity       uint                      `json:"quantity"`
	Price          uint                      `json:"price"`
}

func NewOrderDetailDto(orderDetail entity.OrderDetail) *OrderDetailDTO {
	return &OrderDetailDTO{
		Id:             orderDetail.Id,
		OrderId:        orderDetail.OrderId,
		PharmacyDrugId: orderDetail.PharmacyDrugId,
		Quantity:       orderDetail.Quantity,
		Price:          orderDetail.Price,
	}
}

func NewGetOrderDetailDto(orderDetail entity.OrderDetail) *GetOrderDetailDTO {
	pharmacyDrug, _ := NewGetPharmacyDrugJoinDrugsDto(orderDetail.PharmacyDrug)
	return &GetOrderDetailDTO{
		Id:             orderDetail.Id,
		OrderId:        orderDetail.OrderId,
		PharmacyDrugId: orderDetail.PharmacyDrugId,
		Quantity:       orderDetail.Quantity,
		Price:          orderDetail.Price,
		PharmacyDrug:   *pharmacyDrug,
	}
}
