package response

import "Alice-Seahat-Healthcare/seahat-be/entity"

type StockJournalResponse struct {
	Id          uint    `json:"stock_journal_id"`
	DrugId      uint    `json:"drug_id"`
	DrugName    string  `json:"drug_name"`
	PharmacyId  uint    `json:"pharmacy_id"`
	Quantity    int     `json:"quantity"`
	Description string  `json:"description"`
	CreatedAt   *string `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
	DeletedAt   *string `json:"deleted_at"`
}

func NewStockJournalResponse(resp entity.StockJurnal) *StockJournalResponse {
	var createdAt, deletedAt, updatedAt *string
	if resp.CreatedAt != nil {
		created := resp.CreatedAt.Time.String()
		createdAt = &created
	}
	if resp.DeletedAt != nil {
		deleted := resp.DeletedAt.Time.String()
		deletedAt = &deleted
	}
	if resp.UpdatedAt != nil {
		updated := resp.UpdatedAt.Time.String()
		updatedAt = &updated
	}
	return &StockJournalResponse{
		Id:          resp.Id,
		DrugId:      resp.DrugId,
		DrugName:    resp.DrugName,
		PharmacyId:  resp.PharmacyId,
		Quantity:    resp.Quantity,
		Description: resp.Description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		DeletedAt:   deletedAt,
	}
}
