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

type TelemedicineHandler struct {
	telemedicineUsecase usecase.TelemedicineUsecase
}

func NewTelemedicineHandler(telemedicineUsecase usecase.TelemedicineUsecase) *TelemedicineHandler {
	return &TelemedicineHandler{
		telemedicineUsecase: telemedicineUsecase,
	}
}

func (h *TelemedicineHandler) AddTelemedicine(ctx *gin.Context) {
	body := new(request.AddTelemedicine)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	telemedicine, err := h.telemedicineUsecase.CreateTelemedicine(ctx, body.Telemedicine())
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataCreatedMsg,
		Data:    response.NewTelemedicineDTO(*telemedicine),
	})
}

func (h *TelemedicineHandler) GetTelemedicineByID(ctx *gin.Context) {
	id := ctx.Param("id")
	telemedicineID, err := strconv.Atoi(id)
	if err != nil || telemedicineID < 1 {
		ctx.Error(apperror.InvalidIdParams)
		return
	}

	telemedicine, err := h.telemedicineUsecase.GetTelemedicineByID(ctx, uint(telemedicineID))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    response.NewUserDoctorTelemedicineDTO(*telemedicine),
	})
}

func (h *TelemedicineHandler) GetAllTelemedicine(ctx *gin.Context) {
	collection := request.GetCollectionQuery(ctx)
	telemedicines, err := h.telemedicineUsecase.GetAllTelemedicine(ctx, &collection)
	if err != nil {
		ctx.Error(err)
		return
	}

	res := make([]response.UserDoctorTelemedicineDTO, 0)
	for _, telemedicine := range telemedicines {
		res = append(res, response.NewUserDoctorTelemedicineDTO(telemedicine))
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       res,
		Pagination: response.NewPaginationDto(collection),
	})
}

func (h *TelemedicineHandler) UpdateTelemedicine(ctx *gin.Context) {
	body := new(request.PutTelemedicine)
	id := ctx.Param("id")
	telemedicineID, err := strconv.Atoi(id)
	if err != nil || telemedicineID < 1 {
		ctx.Error(apperror.InvalidIdParams)
		return
	}

	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	err = h.telemedicineUsecase.UpdateOne(ctx, body.Telemedicine(telemedicineID))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataEditMsg,
	})
}

func (h *TelemedicineHandler) AddPrescriptions(ctx *gin.Context) {
	telemedicineID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || telemedicineID < 1 {
		ctx.Error(apperror.InvalidParam)
		return
	}

	body := new(request.AddPrescriptionRequest)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	url, err := h.telemedicineUsecase.AddManyPrescriptedDrugs(ctx, body.Prescription(uint(telemedicineID)))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataEditMsg,
		Data:    url,
	})
}
func (h *TelemedicineHandler) GenerateMedicalCertificate(ctx *gin.Context) {
	id := ctx.Param("id")
	telemedicineId, err := strconv.Atoi(id)

	if err != nil {
		ctx.Error(err)
		return
	}

	req := new(request.TelemedicineReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}
	req.TelemedicineID = uint(telemedicineId)
	body := req.Telemedicine()
	url, err := h.telemedicineUsecase.UpdateOneAndCreateMedicalCertificate(ctx, body)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataCreatedUpdatedMsg,
		Data:    url,
	})

}
