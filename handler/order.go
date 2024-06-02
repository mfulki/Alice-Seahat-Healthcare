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

type OrderHandler struct {
	orderUsecase usecase.OrderUsecase
}

func NewOrderHandler(orderUsecase usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{
		orderUsecase: orderUsecase,
	}
}

func (h *OrderHandler) CreateOrder(ctx *gin.Context) {
	req := new(request.CreateOrder)
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	ordersReq := req.OrderDTO()
	orders, err := h.orderUsecase.CreateOrder(ctx, ordersReq)
	if err != nil {
		ctx.Error(err)
		return
	}

	var orderPayment []*entity.Order
	payment := orders[0].Payment
	for i := 0; i < len(orders); i++ {
		var orderAddres entity.Order
		orders[i].Payment = nil
		orderAddres = orders[i]
		orderPayment = append(orderPayment, &orderAddres)
	}
	payment.Orders = orderPayment
	res := response.NewPaymentDto(payment)

	ctx.JSON(http.StatusCreated, response.Body{
		Message: constant.OrderCreatedSuccessfully,
		Data:    res,
	})
}

func (h *OrderHandler) UpdateConfirmOrder(ctx *gin.Context) {
	orderReq := entity.Order{}
	id := ctx.Param("id")
	orderId, err := strconv.Atoi(id)
	if err != nil || orderId < 1 {
		ctx.Error(apperror.InvalidParam)
		return
	}

	orderReq.Id = uint(orderId)
	order, err := h.orderUsecase.UpdateConfirmOrder(ctx, orderReq)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.OrderConfirmedMsg,
		Data:    response.NewOrderDto(*order),
	})

}

func (h *OrderHandler) OrderProceed(ctx *gin.Context) {
	orderReq := entity.Order{}
	id := ctx.Param("id")
	orderId, err := strconv.Atoi(id)
	if err != nil || orderId < 1 {
		ctx.Error(apperror.InvalidParam)
		return
	}

	orderReq.Id = uint(orderId)

	_, err = h.orderUsecase.OrderProceed(ctx, orderReq)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.OrderConfirmed,
	})

}
func (h *OrderHandler) GetOrderByPharmacyManagerId(ctx *gin.Context) {
	orders, err := h.orderUsecase.GetAllOrderByPharmacyManagerId(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	res := make([]*response.OrderGetDTO, 0)
	for _, order := range orders {
		res = append(res, response.NewOrderGetDto(*order))
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    res,
	})
}
func (h *OrderHandler) OrderSent(ctx *gin.Context) {
	orderReq := entity.Order{}
	id := ctx.Param("id")
	orderId, err := strconv.Atoi(id)
	if err != nil || orderId < 1 {
		ctx.Error(apperror.InvalidParam)
		return
	}

	orderReq.Id = uint(orderId)
	err = h.orderUsecase.OrderSent(ctx, orderReq)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.OrderSend,
	})
}
func (h *OrderHandler) OrderCancelByPM(ctx *gin.Context) {
	orderReq := entity.Order{}
	id := ctx.Param("id")
	orderId, err := strconv.Atoi(id)
	if err != nil || orderId < 1 {
		ctx.Error(apperror.InvalidParam)
		return
	}

	orderReq.Id = uint(orderId)
	err = h.orderUsecase.OrderCancelByPM(ctx, orderReq)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.CancelOrderMsg,
	})
}
