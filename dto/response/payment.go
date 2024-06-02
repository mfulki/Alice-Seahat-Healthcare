package response

import (
	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type PaymentDTO struct {
	Id              uint        `json:"payment_id"`
	UserId          uint        `json:"user_id"`
	UserName        string      `json:"user_name"`
	Method          string      `json:"payment_method"`
	Proof           *string     `json:"payment_proof"`
	FullUserAddress string      `json:"full_user_address,omitempty"`
	TotalPrice      int         `json:"total_price"`
	Number          string      `json:"payment_number"`
	Status          string      `json:"payment_status"`
	Orders          []*OrderDTO `json:"orders,omitempty"`
}
type GetPaymentDTO struct {
	Id              uint           `json:"payment_id"`
	UserId          uint           `json:"user_id"`
	UserName        string         `json:"user_name"`
	Method          string         `json:"payment_method"`
	Proof           *string        `json:"payment_proof"`
	FullUserAddress string         `json:"full_user_address"`
	TotalPrice      int            `json:"total_price"`
	Number          string         `json:"payment_number"`
	Status          string         `json:"payment_status"`
	ExpiredAt       *string        `json:"expired_at"`
	CreatedAt       *string        `json:"created_at"`
	DeletedAt       *string        `json:"deleted_at"`
	Orders          []*OrderGetDTO `json:"orders"`
}

func NewPaymentDto(payment *entity.Payment) *PaymentDTO {
	if payment == nil {
		return nil
	}

	orders := make([]*OrderDTO, 0)
	for _, order := range payment.Orders {
		orders = append(orders, NewOrderDto(*order))
	}

	return &PaymentDTO{
		Id:              payment.Id,
		UserId:          payment.UserId,
		UserName:        payment.UserName,
		Method:          payment.Method,
		Proof:           payment.Proof,
		FullUserAddress: payment.FullUserAddress,
		TotalPrice:      payment.TotalPrice,
		Number:          payment.Number,
		Status:          payment.Status,
		Orders:          orders,
	}
}
func NewGetPaymentDto(payment *entity.Payment) *GetPaymentDTO {
	if payment == nil {
		return nil
	}

	orders := make([]*OrderGetDTO, 0)
	for _, order := range payment.Orders {
		orders = append(orders, NewOrderGetDto(*order))
	}
	var exp, del, create *string
	if payment.ExpiredAt != nil {
		expired := payment.ExpiredAt.Time.String()
		exp = &expired
	}
	if payment.DeletedAt != nil {
		delete := payment.DeletedAt.Time.String()
		del = &delete
	}
	if payment.CreatedAt != nil {
		created := payment.CreatedAt.Time.String()
		create = &created
	}

	return &GetPaymentDTO{
		Id:              payment.Id,
		UserId:          payment.UserId,
		UserName:        payment.UserName,
		Method:          payment.Method,
		Proof:           payment.Proof,
		FullUserAddress: payment.FullUserAddress,
		TotalPrice:      payment.TotalPrice,
		Number:          payment.Number,
		Status:          payment.Status,
		Orders:          orders,
		ExpiredAt:       exp,
		CreatedAt:       create,
		DeletedAt:       del,
	}
}
