package request

import "Alice-Seahat-Healthcare/seahat-be/entity"

type StockRequest struct {
	SenderPharmacyId   uint                `json:"sender_pharmacy_id" binding:"gt=0,gte=1,required"`
	ReceiverPharmacyId uint                `json:"receiver_pharmacy_id" binding:"gt=0,gte=1,required"`
	StockRequestDrug   []*StockRequestDrug `json:"stock_request_drug" binding:"gt=0,dive,required"`
}
type ActionRequest struct {
	StockRequestId uint `uri:"id" binding:"required,gt=0,gte=1"`
}

type DrugStockRequest struct {
	SenderPharmacyID   uint `json:"sender_pharmacy_id" binding:"required,gte=1"`
	ReceiverPharmacyID uint `json:"receiver_pharmacy_id" binding:"required,gte=1"`
}

func (req *ActionRequest) UpdateStockMutationDTO() entity.StockRequest {
	return entity.StockRequest{
		Id: req.StockRequestId,
	}
}

func (req *StockRequest) StockRequestDto() entity.StockRequest {
	stockRequestDrug := NewStockRequestDrug(req.StockRequestDrug)
	senderPharmacy := entity.Pharmacy{ID: req.SenderPharmacyId}
	receiverPharmacy := entity.Pharmacy{ID: req.ReceiverPharmacyId}
	return entity.StockRequest{
		SenderPharmacy:   senderPharmacy,
		ReceiverPharmacy: receiverPharmacy,
		StockRequestDrug: stockRequestDrug,
	}
}
