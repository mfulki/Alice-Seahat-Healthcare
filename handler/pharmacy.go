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

type PharmacyHandler struct {
	pharmacyUsecase usecase.PharmacyUsecase
}

func NewPharmacyHandler(pharmacyUsecase usecase.PharmacyUsecase) *PharmacyHandler {
	return &PharmacyHandler{
		pharmacyUsecase: pharmacyUsecase,
	}
}

func (h *PharmacyHandler) AddPharmacy(ctx *gin.Context) {
	body := new(request.AddPharmacy)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	pharmacy, err := h.pharmacyUsecase.CreatePharmacy(ctx, *body.Pharmacy(), body.ShipmentMethods)
	pharmacy.Latitude = body.Latitude
	pharmacy.Longitude = body.Longitude

	if err != nil {
		ctx.Error(err)
		return
	}

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataCreatedMsg,
		Data:    response.NewCreatePharmacyDTO(*pharmacy),
	})
}

func (h *PharmacyHandler) EditPharmacy(ctx *gin.Context) {
	id := ctx.Param("id")
	pharmacyID, err := strconv.Atoi(id)
	if err != nil || pharmacyID < 1 {
		ctx.Error(apperror.InvalidParam)
		return
	}

	body := new(request.EditPharmacy)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	pharmacy, err := h.pharmacyUsecase.UpdatePharmacy(ctx, *body.Pharmacy(), uint(pharmacyID), body.ShipmentMethods)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataEditMsg,
		Data:    response.NewCreatePharmacyDTO(*pharmacy),
	})
}

func (h *PharmacyHandler) GetAllPharmacies(ctx *gin.Context) {
	collection := request.GetCollectionQuery(ctx)
	pharmaciesDTO := []response.PharmacyDto{}
	pharmacies, err := h.pharmacyUsecase.GetAllPharmacies(ctx, &collection)
	if err != nil {
		ctx.Error(err)
		return
	}

	for i := 0; i < len(pharmacies); i++ {
		pharmaciesDTO = append(pharmaciesDTO, *response.NewCreatePharmacyDTO(pharmacies[i]))
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       pharmaciesDTO,
		Pagination: response.NewPaginationDto(collection),
	})
}

func (h *PharmacyHandler) GetAllPharmaciesByManagerID(ctx *gin.Context) {
	id := ctx.Param("id")
	managerID, err := strconv.Atoi(id)
	if err != nil || managerID < 1 {
		ctx.Error(apperror.InvalidParam)
		return
	}

	collection := request.GetCollectionQuery(ctx)
	pharmaciesDTO := []response.PharmacyDto{}
	pharmacies, err := h.pharmacyUsecase.GetAllPharmaciesByManagerID(ctx, uint(managerID), &collection)
	if err != nil {
		ctx.Error(err)
		return
	}

	for i := 0; i < len(pharmacies); i++ {
		pharmaciesDTO = append(pharmaciesDTO, *response.NewCreatePharmacyDTO(pharmacies[i]))
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       pharmaciesDTO,
		Pagination: response.NewPaginationDto(collection),
	})
}

func (h *PharmacyHandler) GetPharmacyByID(ctx *gin.Context) {
	id := ctx.Param("id")
	pharmacyID, err := strconv.Atoi(id)
	if err != nil || pharmacyID < 1 {
		ctx.Error(apperror.InvalidParam)
		return
	}

	pharmacy, err := h.pharmacyUsecase.GetPharmacyByID(ctx, uint(pharmacyID))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    response.NewCreatePharmacyDTO(*pharmacy),
	})
}
