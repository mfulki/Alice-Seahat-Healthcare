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

type PartnerHandler struct {
	partnerUsecase usecase.PartnerUsecase
}

func NewPartnerHandler(partnerUsecase usecase.PartnerUsecase) *PartnerHandler {
	return &PartnerHandler{
		partnerUsecase: partnerUsecase,
	}
}

func (h *PartnerHandler) GetAll(ctx *gin.Context) {
	collection := request.GetCollectionQuery(ctx)
	partners, err := h.partnerUsecase.GetAll(ctx, &collection)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       response.NewMultiplePartnerDto(partners),
		Pagination: response.NewPaginationDto(collection),
	})
}

func (h *PartnerHandler) CreatePartner(ctx *gin.Context) {
	body := new(request.CreatePartner)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	partner, err := h.partnerUsecase.CreatePartner(ctx, body.Partner())
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response.Body{
		Message: constant.DataCreatedMsg,
		Data:    response.NewPartnerDto(*partner),
	})
}

func (h *PartnerHandler) GetPartnerByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id < 1 {
		ctx.Error(apperror.InvalidParam)
		return
	}

	p, err := h.partnerUsecase.GetPartnerByID(ctx, uint(id))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    response.NewPartnerDto(*p),
	})
}

func (h *PartnerHandler) UpdatePartnerByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id < 1 {
		ctx.Error(apperror.InvalidParam)
		return
	}

	body := new(request.UpdatePartner)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	partner := body.Partner()
	partner.ID = uint(id)

	if err := h.partnerUsecase.UpdatePartnerByID(ctx, partner); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response.Body{
		Message: constant.DataEditMsg,
	})
}
