package response

import "Alice-Seahat-Healthcare/seahat-be/entity"

type StockRequestDrug struct {
	Id             uint        `json:"stock_request_drug_id"`
	StockRequestId uint        `json:"stock_request_id" `
	Drug           GetDrugName `json:"drug" `
	Quantity       int         `json:"quantity" `
	CreatedAt      *string     `json:"created_at"`
	UpdatedAt      *string     `json:"updated_at"`
	DeletedAt      *string     `json:"deleted_at"`
}

func NewStockRequestDrug(resp []*entity.StockRequestDrug) []*StockRequestDrug {
	stockRequestDrugs := make([]*StockRequestDrug, 0)

	for _, res := range resp {
		var createdAt, deletedAt, updatedAt *string
		if res.CreatedAt != nil {
			created := res.CreatedAt.Time.String()
			createdAt = &created
		}
		if res.DeletedAt != nil {
			deleted := res.DeletedAt.Time.String()
			deletedAt = &deleted
		}
		if res.UpdatedAt != nil {
			updated := res.UpdatedAt.Time.String()
			updatedAt = &updated
		}
		drug := NewGetDrug(res.Drug)
		stockRequestDrug := &StockRequestDrug{
			Id:             res.Id,
			StockRequestId: res.StockRequestId,
			Drug:           *drug,
			Quantity:       res.Quantity,
			CreatedAt:      createdAt,
			UpdatedAt:      updatedAt,
			DeletedAt:      deletedAt,
		}
		stockRequestDrugs = append(stockRequestDrugs, stockRequestDrug)
	}
	return stockRequestDrugs

}
