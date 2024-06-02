package handler

import (
	"net/http"

	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/dto/request"
	"Alice-Seahat-Healthcare/seahat-be/dto/response"
	"Alice-Seahat-Healthcare/seahat-be/usecase"

	"github.com/gin-gonic/gin"
)

type ManufacturerHandler struct {
	manufacturerUsecase usecase.ManufacturerUsecase
}

func NewManufacturerHandler(manufacturerUsecase usecase.ManufacturerUsecase) *ManufacturerHandler {
	return &ManufacturerHandler{
		manufacturerUsecase: manufacturerUsecase,
	}
}

func (h *ManufacturerHandler) GetAllManufacturers(ctx *gin.Context) {
	collection := request.GetCollectionQuery(ctx)
	manufacturers, err := h.manufacturerUsecase.GetAll(ctx, &collection)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       response.NewMultipleManufacturerDto(manufacturers),
		Pagination: response.NewPaginationDto(collection),
	})
}
