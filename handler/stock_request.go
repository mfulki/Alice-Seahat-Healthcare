package handler

import (
	"net/http"

	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/dto/request"
	"Alice-Seahat-Healthcare/seahat-be/dto/response"
	"Alice-Seahat-Healthcare/seahat-be/usecase"

	"github.com/gin-gonic/gin"
)

type StockRequestHandler struct {
	stockRequestUsecase usecase.StockRequestUsecase
}

func NewStockRequestHandler(stockRequestUsecase usecase.StockRequestUsecase) *StockRequestHandler {
	return &StockRequestHandler{
		stockRequestUsecase: stockRequestUsecase,
	}
}

func (h *StockRequestHandler) StockMutationManualRequest(ctx *gin.Context) {
	req := new(request.StockRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	stockReq := req.StockRequestDto()
	_, err := h.stockRequestUsecase.ManualStockRequest(ctx, stockReq)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response.Body{
		Message: constant.StockRequestCreated,
	})
}

func (h *StockRequestHandler) StockMutationManualApprove(ctx *gin.Context) {
	req := new(request.ActionRequest)
	if err := ctx.ShouldBindUri(req); err != nil {
		ctx.Error(err)
		return
	}
	body := req.UpdateStockMutationDTO()
	stockReq, err := h.stockRequestUsecase.UpdateStockMutationApprove(ctx, body)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, response.Body{
		Message: constant.StockRequestCreated,
		Data:    response.NewStockRequestResponseDto(stockReq),
	})
}
func (h *StockRequestHandler) StockMutationManualCancel(ctx *gin.Context) {
	req := new(request.ActionRequest)
	if err := ctx.ShouldBindUri(req); err != nil {
		ctx.Error(err)
		return
	}
	body := req.UpdateStockMutationDTO()

	stockReq, err := h.stockRequestUsecase.CancelStockMutation(ctx, body)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, response.Body{
		Message: constant.StockRequestCreated,
		Data:    response.NewStockRequestResponseDto(stockReq),
	})
}

func (h *StockRequestHandler) GetAllStockRequest(ctx *gin.Context) {
	collection := request.GetCollectionQuery(ctx)
	stockRequest, err := h.stockRequestUsecase.GetAllStockRequest(ctx, &collection)
	if err != nil {
		ctx.Error(err)
		return
	}

	res := make([]*response.StockRequest, 0)
	for _, stockRequest := range stockRequest {
		res = append(res, response.NewStockRequestResponseDto(stockRequest))
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       res,
		Pagination: response.NewPaginationDto(collection),
	})

}

func (h *StockRequestHandler) GetAvailableDrugFromSenderAndReceiverPharmacy(ctx *gin.Context) {
	body := new(request.DrugStockRequest)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	collection := request.GetCollectionQuery(ctx)
	drugs, err := h.stockRequestUsecase.GetDrugsWithSenderAndReceiverPharmacy(ctx, body.SenderPharmacyID, body.ReceiverPharmacyID, &collection)
	if err != nil {
		ctx.Error(err)
		return
	}

	dtos, err := response.NewMultipleDrugWithPharmacyDrugDto(drugs)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       dtos,
		Pagination: response.NewPaginationDto(collection),
	})
}
