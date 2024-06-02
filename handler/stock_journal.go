package handler

import (
	"net/http"

	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/dto/request"
	"Alice-Seahat-Healthcare/seahat-be/dto/response"
	"Alice-Seahat-Healthcare/seahat-be/usecase"

	"github.com/gin-gonic/gin"
)

type StockJournalHandler struct {
	stockJournalUsecase usecase.StockJournalUsecase
}

func NewStockJournalHandler(stockJournalUsecase usecase.StockJournalUsecase) *StockJournalHandler {
	return &StockJournalHandler{
		stockJournalUsecase: stockJournalUsecase,
	}
}

func (h *StockJournalHandler) GetAllStockJournalByPharmacyId(ctx *gin.Context) {
	body := new(request.StockJournalReq)
	if err := ctx.ShouldBindUri(body); err != nil {
		ctx.Error(err)
		return
	}

	stockJournal := body.NewStockJournalReq()
	collection := request.GetCollectionQuery(ctx)
	stockJournals, err := h.stockJournalUsecase.GetAllStockJournalBYPharmacyId(ctx, stockJournal, &collection)
	if err != nil {
		ctx.Error(err)
		return
	}
	res := make([]*response.StockJournalResponse, 0)
	for _, stockJournal := range stockJournals {
		res = append(res, response.NewStockJournalResponse(*stockJournal))
	}
	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       res,
		Pagination: response.NewPaginationDto(collection),
	})

}
