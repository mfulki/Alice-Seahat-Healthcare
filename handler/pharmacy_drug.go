package handler

import (
	"net/http"
	"strconv"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/dto/request"
	"Alice-Seahat-Healthcare/seahat-be/dto/response"
	"Alice-Seahat-Healthcare/seahat-be/usecase"
	"Alice-Seahat-Healthcare/seahat-be/utils"

	"github.com/gin-gonic/gin"
)

type PharmacyDrugHandler struct {
	pharmacyDrugUsecase usecase.PharmacyDrugUsecase
}

func NewPharmacyDrugHandler(pharmacyDrugUsecase usecase.PharmacyDrugUsecase) *PharmacyDrugHandler {
	return &PharmacyDrugHandler{
		pharmacyDrugUsecase: pharmacyDrugUsecase,
	}
}

func (h *PharmacyDrugHandler) GetAllWithinRadius(ctx *gin.Context) {
	collection := request.GetCollectionQuery(ctx)
	query := new(request.PharmaryDrugsQuery)
	if err := ctx.ShouldBindQuery(query); err != nil {
		ctx.Error(err)
		return
	}

	addr := query.Address()
	user := utils.GetUserFromJwt(ctx.Request)
	if user != nil {
		addr.UserID = user.ID
	}

	pharmacyDrugs, err := h.pharmacyDrugUsecase.GetAllWithinRadius(ctx, addr, &collection)
	if err != nil {
		ctx.Error(err)
		return
	}

	res := make([]*response.PharmacyDrugDto, 0)
	for _, pharmacyDrug := range pharmacyDrugs {
		parsed, err := response.NewPharmacyDrugDto(*pharmacyDrug)
		if err != nil {
			ctx.Error(err)
			return
		}
		res = append(res, parsed)
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       res,
		Pagination: response.NewPaginationDto(collection),
	})
}

func (h *PharmacyDrugHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.Error(err)
		return
	}

	pharmacyDrug, err := h.pharmacyDrugUsecase.GetByID(ctx, uint(parsedId))
	if err != nil {
		ctx.Error(err)
		return
	}

	res, err := response.NewPharmacyDrugDto(*pharmacyDrug)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    res,
	})
}

func (h *PharmacyDrugHandler) GetAllPharmacyDrugs(ctx *gin.Context) {
	id := ctx.Param("id")
	collection := request.GetCollectionQuery(ctx)

	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.Error(err)
		return
	}

	pharmacyDrugs, err := h.pharmacyDrugUsecase.GetAllPharmacyDrug(ctx, &collection, uint(parsedId))
	if err != nil {
		ctx.Error(err)
		return
	}

	pharmacyDrugsDTO := []response.GetPharmacyDrugAndDrugDto{}

	for i := 0; i < len(pharmacyDrugs); i++ {
		pharmacyDrug, _ := response.NewGetPharmacyDrugJoinDrugsDto(pharmacyDrugs[i])
		pharmacyDrugsDTO = append(pharmacyDrugsDTO, *pharmacyDrug)
	}

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       pharmacyDrugsDTO,
		Pagination: response.NewPaginationDto(collection),
	})
}

func (h *PharmacyDrugHandler) CreatePharmacyDrug(ctx *gin.Context) {
	body := new(request.AddPharmacyDrug)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	pharmacyDrug, err := h.pharmacyDrugUsecase.CreatePharmacyDrug(ctx, body.PharmacyDrug())
	if err != nil {
		ctx.Error(err)
		return
	}

	pharmacyDrugDto := response.NewCreatePharmacyDrugDto(*pharmacyDrug)

	ctx.JSON(http.StatusCreated, response.Body{
		Message: constant.DataCreatedMsg,
		Data:    pharmacyDrugDto,
	})
}

func (h *PharmacyDrugHandler) UpdatePharmacyDrug(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ctx.Error(apperror.InvalidParam)
		return
	}

	body := new(request.AddPharmacyDrug)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	err = h.pharmacyDrugUsecase.UpdatePharmacyDrug(ctx, body.PharmacyDrug(), uint(parsedId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response.Body{
		Message: constant.DataEditMsg,
	})
}

func (h *PharmacyDrugHandler) GetNearestPharmacyByDrugID(ctx *gin.Context) {
	drugID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || drugID < 1 {
		ctx.Error(apperror.InvalidParam)
		return
	}

	query := new(request.PharmaryDrugsQuery)
	if err := ctx.ShouldBindQuery(query); err != nil {
		ctx.Error(err)
		return
	}

	addr := query.Address()
	user := utils.GetUserFromJwt(ctx.Request)
	if user != nil {
		addr.UserID = user.ID
	}

	collection := request.GetCollectionQuery(ctx)
	pds, err := h.pharmacyDrugUsecase.GetNearestPharmacies(ctx, uint(drugID), addr, &collection)
	if err != nil {
		ctx.Error(err)
		return
	}

	res := make([]*response.PharmacyDrugDto, 0)
	for _, pharmacyDrug := range pds {
		parsed, err := response.NewPharmacyDrugDto(pharmacyDrug)
		if err != nil {
			ctx.Error(err)
			return
		}
		res = append(res, parsed)
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       res,
		Pagination: response.NewPaginationDto(collection),
	})
}

func (h *PharmacyDrugHandler) GetDrugOfTheDay(ctx *gin.Context) {
	pharmacyDrugs, err := h.pharmacyDrugUsecase.GetAllDrugsOfTheDay(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	res := make([]*response.PharmacyDrugDto, 0)
	for _, pharmacyDrug := range pharmacyDrugs {
		parsed, err := response.NewPharmacyDrugDto(pharmacyDrug)
		if err != nil {
			ctx.Error(err)
			return
		}

		res = append(res, parsed)
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    res,
	})
}
