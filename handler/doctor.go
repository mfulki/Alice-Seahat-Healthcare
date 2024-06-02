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

type DoctorHandler struct {
	doctorUsecase usecase.DoctorUsecase
}

func NewDoctorHandler(doctorUsecase usecase.DoctorUsecase) *DoctorHandler {
	return &DoctorHandler{
		doctorUsecase: doctorUsecase,
	}
}

func (h *DoctorHandler) Login(ctx *gin.Context) {
	body := new(request.DoctorLogin)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	token, err := h.doctorUsecase.Login(ctx, body.Doctor())
	if err != nil {
		ctx.Error(err)
		return
	}

	utils.SendLoginCookie(ctx, token)
	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.LoginPassedMsg,
	})
}

func (h *DoctorHandler) LoginOAuth(ctx *gin.Context) {
	body := new(request.UserLoginOAuth)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	token, err := h.doctorUsecase.LoginOAuth(ctx, body.GoogleToken)
	if err != nil {
		ctx.Error(err)
		return
	}

	utils.SendLoginCookie(ctx, token)
	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.LoginPassedMsg,
	})
}

func (h *DoctorHandler) Register(ctx *gin.Context) {
	body := new(request.DoctorRegister)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	doctor, err := h.doctorUsecase.Register(ctx, body.Doctor())
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response.Body{
		Message: constant.RegisterMsg,
		Data:    response.NewDoctorDto(*doctor),
	})
}

func (h *DoctorHandler) RegisterOAuth(ctx *gin.Context) {
	body := new(request.DoctorRegisterOAuth)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	doctor, err := h.doctorUsecase.RegisterOAuth(ctx, body.Doctor(), body.GoogleToken)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response.Body{
		Message: constant.RegisterMsg,
		Data:    response.NewDoctorDto(*doctor),
	})
}

func (h *DoctorHandler) ForgotPassword(ctx *gin.Context) {
	body := new(request.DoctorForgot)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.doctorUsecase.ForgotPassword(ctx, body.Email); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.ForgotTokenMsg,
	})
}

func (h *DoctorHandler) Verification(ctx *gin.Context) {
	body := new(request.DoctorToken)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	err := h.doctorUsecase.Verification(ctx, body.Password, body.Token)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.VerificationMsg,
	})
}

func (h *DoctorHandler) UpdatePersonal(ctx *gin.Context) {
	body := new(request.DoctorPersonalEdit)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.doctorUsecase.UpdateProfile(ctx, body.Doctor()); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataEditMsg,
	})
}

func (h *DoctorHandler) GetAllDoctors(ctx *gin.Context) {
	collection := request.GetCollectionQuery(ctx)

	doctors, err := h.doctorUsecase.GetAllDoctors(ctx, &collection)

	if err != nil {
		ctx.Error(err)
		return
	}

	doctorDTO := []response.DoctorDto{}
	for i := 0; i < len(doctors); i++ {
		doctorDTO = append(doctorDTO, response.NewDoctorDto(doctors[i]))
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       doctorDTO,
		Pagination: response.NewPaginationDto(collection),
	})
}

func (h *DoctorHandler) GetDoctorByID(ctx *gin.Context) {
	id := ctx.Param("id")
	doctorId, err := strconv.Atoi(id)
	if err != nil || doctorId < 1 {
		ctx.Error(apperror.InvalidParam)
		return
	}

	doctor, err := h.doctorUsecase.GetDoctorByID(ctx, uint(doctorId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    response.NewDoctorDto(*doctor),
	})
}

func (h *DoctorHandler) UpdatePassword(ctx *gin.Context) {
	body := new(request.UserPasswordEdit)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.doctorUsecase.UpdatePassword(ctx, body.OldPassword, body.NewPassword); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataEditMsg,
	})
}

func (h *DoctorHandler) GetProfile(ctx *gin.Context) {
	doctor, err := h.doctorUsecase.GetProfile(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    response.NewDoctorDto(*doctor),
	})
}

func (h *DoctorHandler) UpdateStatus(ctx *gin.Context) {
	body := new(request.DoctorStatus)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.doctorUsecase.UpdateStatus(ctx, body.Doctor()); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataEditMsg,
	})
}
