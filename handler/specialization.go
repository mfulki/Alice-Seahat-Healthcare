package handler

import (
	"net/http"

	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/dto/response"
	"Alice-Seahat-Healthcare/seahat-be/usecase"

	"github.com/gin-gonic/gin"
)

type SpecializationHandler struct {
	specializationUsecase usecase.SpecializationUsecase
}

func NewSpecializationHandler(specializationUsecase usecase.SpecializationUsecase) *SpecializationHandler {
	return &SpecializationHandler{
		specializationUsecase: specializationUsecase,
	}
}

func (h *SpecializationHandler) GetAll(ctx *gin.Context) {
	szs, err := h.specializationUsecase.GetAll(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    response.NewMultipleSpecializationDto(szs),
	})
}
