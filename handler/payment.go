package handler

import (
	"net/http"
	"strconv"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/dto/request"
	"Alice-Seahat-Healthcare/seahat-be/dto/response"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/usecase"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentUsecase usecase.PaymentUsecase
}

func NewPaymentHandler(paymentUsecase usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{
		paymentUsecase: paymentUsecase,
	}
}

func (h *PaymentHandler) UpdatePaymentProof(ctx *gin.Context) {
	req := new(request.PaymentProof)
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	paymentReq := req.UpdatePaymentProof()

	id := ctx.Param("id")
	paymentId, err := strconv.Atoi(id)
	if err != nil || paymentId < 1 {
		ctx.Error(apperror.InvalidParam)
		return
	}

	paymentReq.Id = uint(paymentId)

	if err != nil {
		ctx.Error(err)
		return
	}
	orders, err := h.paymentUsecase.UpdatePaymentProof(ctx, paymentReq)
	if err != nil {
		ctx.Error(err)
		return
	}

	res := make([]*response.OrderDTO, 0)
	for _, order := range orders {
		res = append(res, response.NewOrderDto(*order))
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.UpdatePaymentProofMsg,
		Data:    res,
	})

}

func (h *PaymentHandler) UserCancelPayment(ctx *gin.Context) {
	paymentReq := entity.Payment{}
	id := ctx.Param("id")
	paymentId, err := strconv.Atoi(id)
	if err != nil || paymentId < 1 {
		ctx.Error(apperror.InvalidIdParams)
		return
	}

	paymentReq.Id = uint(paymentId)

	orders, err := h.paymentUsecase.UserCancelPayment(ctx, paymentReq)
	if err != nil {
		ctx.Error(err)
		return
	}

	res := make([]*response.OrderDTO, 0)
	for _, order := range orders {
		res = append(res, response.NewOrderDto(*order))
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.CancelPaymentMsg,
		Data:    res,
	})

}
func (h *PaymentHandler) AdminCancelPayment(ctx *gin.Context) {
	paymentReq := entity.Payment{}
	id := ctx.Param("id")
	paymentId, err := strconv.Atoi(id)
	if err != nil || paymentId < 1 {
		ctx.Error(apperror.InvalidIdParams)
		return
	}

	paymentReq.Id = uint(paymentId)

	orders, err := h.paymentUsecase.AdminCancelPayment(ctx, paymentReq)
	if err != nil {
		ctx.Error(err)
		return
	}

	res := make([]*response.OrderDTO, 0)
	for _, order := range orders {
		res = append(res, response.NewOrderDto(*order))
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.CancelPaymentMsg,
		Data:    res,
	})

}
func (h *PaymentHandler) AdminRejectPayment(ctx *gin.Context) {
	paymentReq := entity.Payment{}
	id := ctx.Param("id")
	paymentId, err := strconv.Atoi(id)
	if err != nil || paymentId < 1 {
		ctx.Error(apperror.InvalidIdParams)
		return
	}

	paymentReq.Id = uint(paymentId)

	err = h.paymentUsecase.AdminRejectPayment(ctx, paymentReq)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.CancelPaymentMsg,
	})

}
func (h *PaymentHandler) PaymentConfirmation(ctx *gin.Context) {
	req := new(request.AdminActionPayment)
	if err := ctx.ShouldBindUri(req); err != nil {
		ctx.Error(err)
		return
	}
	body := req.UpdateAction()
	_, err := h.paymentUsecase.PaymentConfirmation(ctx, body)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.PaymentConfirmed,
	})

}

func (h *PaymentHandler) GetAllPaymentToConfirm(ctx *gin.Context) {
	collection := request.GetCollectionQuery(ctx)
	payments, err := h.paymentUsecase.GetAllPaymentToConfirm(ctx, &collection)
	if err != nil {
		ctx.Error(err)
		return
	}
	res := make([]*response.PaymentDTO, 0)
	for _, payment := range payments {
		res = append(res, response.NewPaymentDto(payment))
	}
	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       res,
		Pagination: response.NewPaginationDto(collection),
	})
}

func (h *PaymentHandler) GetAllPaymentByUserId(ctx *gin.Context) {
	collection := request.GetCollectionQuery(ctx)
	payments, err := h.paymentUsecase.GetAllPaymentByUserId(ctx, &collection)
	if err != nil {
		ctx.Error(err)
		return
	}
	res := make([]*response.GetPaymentDTO, 0)
	for _, payment := range payments {
		res = append(res, response.NewGetPaymentDto(payment))
	}
	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       res,
		Pagination: response.NewPaginationDto(collection),
	})
}
