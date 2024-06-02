package request

import "Alice-Seahat-Healthcare/seahat-be/entity"

type Payment struct {
	Method    string `json:"payment_method" binding:"required,min=4" `
	AddressId uint   `json:"address_id" binding:"gt=0,gte=1,required"`
}

type PaymentProof struct {
	Proof string `json:"payment_proof" binding:"required,min=10"`
}
type AdminActionPayment struct {
	Id uint `uri:"id"`
}

func (req *PaymentProof) UpdatePaymentProof() entity.Payment {
	return entity.Payment{
		Proof: &req.Proof,
	}

}

func NewPayment(req Payment) entity.Payment {
	address := entity.Address{ID: req.AddressId}
	return entity.Payment{
		Method:  req.Method,
		Address: &address,
	}
}
func (req *AdminActionPayment) UpdateAction() entity.Payment {
	return entity.Payment{
		Id: req.Id,
	}
}
