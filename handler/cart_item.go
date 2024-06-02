package handler

import (
	"net/http"
	"strconv"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/dto/request"
	"Alice-Seahat-Healthcare/seahat-be/dto/response"
	"Alice-Seahat-Healthcare/seahat-be/usecase"

	"github.com/gin-gonic/gin"
)

type CartItemHandler struct {
	cartItemUsecase usecase.CartItemUsecase
}

func NewCartItemHandler(cartItemUsecase usecase.CartItemUsecase) *CartItemHandler {
	return &CartItemHandler{
		cartItemUsecase: cartItemUsecase,
	}
}

func (h *CartItemHandler) GetAllCartItem(ctx *gin.Context) {
	items, err := h.cartItemUsecase.GetAllCartItem(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    response.NewCartItemsWithSummary(items),
	})
}

func (h *CartItemHandler) CreateCartItem(ctx *gin.Context) {
	body := new(request.AddCartItemRequest)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	item, err := h.cartItemUsecase.AddCartItem(ctx, body.CartItem())
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response.Body{
		Message: constant.DataCreatedMsg,
		Data:    response.NewCartItemDto(*item),
	})
}

func (h *CartItemHandler) UpdateQtyCartItem(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id < 1 {
		ctx.Error(apperror.InvalidParam)
		return
	}

	body := new(request.UpdateQtyItemRequest)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	err = h.cartItemUsecase.UpdateCartItem(ctx, body.CartItem(id))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataEditMsg,
	})
}

func (h *CartItemHandler) DeleteManyCartItem(ctx *gin.Context) {
	body := new(request.DeleteCartItemRequest)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	err := h.cartItemUsecase.DeleteBulkCartItem(ctx, body.Uint())
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataDeletedMsg,
	})
}
