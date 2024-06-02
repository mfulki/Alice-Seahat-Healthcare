package request

import "Alice-Seahat-Healthcare/seahat-be/entity"

type StockJournalReq struct {
	PharmacyId uint `uri:"id" binding:"required,gt=0,gte=1"`
}

func (req *StockJournalReq) NewStockJournalReq() entity.StockJurnal {
	return entity.StockJurnal{
		PharmacyId: req.PharmacyId,
	}

}
