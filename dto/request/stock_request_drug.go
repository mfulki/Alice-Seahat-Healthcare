package request

import "Alice-Seahat-Healthcare/seahat-be/entity"

type StockRequestDrug struct {
	Id             uint `json:"stock_request_drug_id"`
	StockRequestId uint `json:"stock_request_id" `
	DrugId         uint `json:"drug_id" binding:"gt=0,gte=1,required"`
	Quantity       int  `json:"quantity" binding:"gt=0,gte=1,required"`
}

func NewStockRequestDrug(reqs []*StockRequestDrug) []*entity.StockRequestDrug {
	stockRequestDrugs := make([]*entity.StockRequestDrug, 0)
	for _, req := range reqs {
		stockRequestDrug := &entity.StockRequestDrug{StockRequestId: req.StockRequestId, DrugId: req.DrugId, Quantity: req.Quantity}
		stockRequestDrugs = append(stockRequestDrugs, stockRequestDrug)
	}
	return stockRequestDrugs

}
