package handler

import (
	"net/http"

	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/dto/request"
	"Alice-Seahat-Healthcare/seahat-be/dto/response"
	"Alice-Seahat-Healthcare/seahat-be/usecase"
	"Alice-Seahat-Healthcare/seahat-be/utils"

	"github.com/gin-gonic/gin"
)

type PharmacyManagerHandler struct {
	pharmacyManager usecase.PharmacyManagerUsecase
}

func NewPharmacyManagerHandler(pharmacyManager usecase.PharmacyManagerUsecase) *PharmacyManagerHandler {
	return &PharmacyManagerHandler{
		pharmacyManager: pharmacyManager,
	}
}

func (h *PharmacyManagerHandler) Login(ctx *gin.Context) {
	body := new(request.UserLogin)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	token, err := h.pharmacyManager.Login(ctx, body.PharmacyManager())
	if err != nil {
		ctx.Error(err)
		return
	}

	utils.SendLoginCookie(ctx, token)
	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.LoginPassedMsg,
	})
}

func (h *PharmacyManagerHandler) GetProfile(ctx *gin.Context) {
	manager, err := h.pharmacyManager.GetProfile(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    response.NewPharmacyManagerDto(*manager),
	})
}
