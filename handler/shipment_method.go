package handler

import (
	"net/http"

	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/dto/response"
	"Alice-Seahat-Healthcare/seahat-be/usecase"

	"github.com/gin-gonic/gin"
)

type ShipmentMethodHandler struct {
	shipmentMethodUsecase usecase.ShipmentMethodUsecase
}

func NewShipmentMethodHandler(shipmentMethodUsecase usecase.ShipmentMethodUsecase) *ShipmentMethodHandler {
	return &ShipmentMethodHandler{
		shipmentMethodUsecase: shipmentMethodUsecase,
	}
}

func (h *ShipmentMethodHandler) GetAllShipmentMethods(ctx *gin.Context) {
	shipmentMethods, err := h.shipmentMethodUsecase.GetAllShipmentMethod(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    response.NewMultipleShipmentMethodDto(shipmentMethods),
	})
}
