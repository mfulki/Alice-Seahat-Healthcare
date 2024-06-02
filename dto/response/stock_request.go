package response

import "Alice-Seahat-Healthcare/seahat-be/entity"

type StockRequest struct {
	Id               uint                `json:"stock_request_id"`
	SenderPharmacy   *GetPharmacy        `json:"sender_pharmacy"`
	ReceiverPharmacy *GetPharmacy        `json:"receiver_pharmacy"`
	StockRequestDrug []*StockRequestDrug `json:"stock_request_drug,omitempty"`
	Status           *string             `json:"status"`
	CreatedAt        *string             `json:"created_at"`
	UpdatedAt        *string             `json:"updated_at"`
	DeletedAt        *string             `json:"deleted_at"`
}

type DrugWithPharmacyDrugDto struct {
	DrugDto
	SenderPharmacyDrug   PharmacyDrugDto `json:"sender_pharmacy_drug"`
	ReceiverPharmacyDrug PharmacyDrugDto `json:"receiver_pharmacy_drug"`
}

func NewStockRequestResponseDto(resp *entity.StockRequest) *StockRequest {
	stockRequestDrug := NewStockRequestDrug(resp.StockRequestDrug)
	var createdAt, deletedAt, updatedAt *string
	if resp.CreatedAt != nil {
		created := resp.CreatedAt.Time.String()
		createdAt = &created
	}
	if resp.DeletedAt != nil {
		deleted := resp.DeletedAt.Time.String()
		deletedAt = &deleted
	}
	receiverPharmacy := NewGetPharmacies(resp.ReceiverPharmacy)
	senderPharmacy := NewGetPharmacies(resp.SenderPharmacy)
	return &StockRequest{
		Id:               resp.Id,
		SenderPharmacy:   senderPharmacy,
		ReceiverPharmacy: receiverPharmacy,
		Status:           &resp.Status,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		DeletedAt:        deletedAt,
		StockRequestDrug: stockRequestDrug,
	}
}

func NewDrugWithPharmacyDrugDto(dpd entity.DrugWithPharmacyDrug) (DrugWithPharmacyDrugDto, error) {
	dto := DrugWithPharmacyDrugDto{}

	senderPharmacyDrug, err := NewPharmacyDrugDto(dpd.SenderPharmacyDrug)
	if err != nil {
		return dto, nil
	}

	receiverPharmacyDrug, err := NewPharmacyDrugDto(dpd.ReceiverPharmacyDrug)
	if err != nil {
		return dto, nil
	}

	dto.DrugDto = *NewDrugDto(dpd.Drug)
	dto.SenderPharmacyDrug = *senderPharmacyDrug
	dto.ReceiverPharmacyDrug = *receiverPharmacyDrug

	return dto, nil
}

func NewMultipleDrugWithPharmacyDrugDto(dpds []entity.DrugWithPharmacyDrug) ([]DrugWithPharmacyDrugDto, error) {
	dtos := make([]DrugWithPharmacyDrugDto, 0)

	for _, dpd := range dpds {
		dto, err := NewDrugWithPharmacyDrugDto(dpd)
		if err != nil {
			return nil, err
		}

		dtos = append(dtos, dto)
	}

	return dtos, nil
}
