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

type AdminHandler struct {
	adminUsecase usecase.AdminUsecase
}

func NewAdminHandler(adminUsecase usecase.AdminUsecase) *AdminHandler {
	return &AdminHandler{
		adminUsecase: adminUsecase,
	}
}

func (h *AdminHandler) Login(ctx *gin.Context) {
	body := new(request.UserLogin)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	token, err := h.adminUsecase.Login(ctx, body.Admin())
	if err != nil {
		ctx.Error(err)
		return
	}

	utils.SendLoginCookie(ctx, token)
	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.LoginPassedMsg,
	})
}

func (h *AdminHandler) GetProfile(ctx *gin.Context) {
	admin, err := h.adminUsecase.GetProfile(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    response.NewAdminDto(*admin),
	})
}

func (h *AdminHandler) GetAllDoctor(ctx *gin.Context) {
	collection := request.GetCollectionQuery(ctx)
	doctors, err := h.adminUsecase.GetAllDoctor(ctx, &collection)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       response.NewMultipleDoctor(doctors),
		Pagination: response.NewPaginationDto(collection),
	})
}

func (h *AdminHandler) GetAllManager(ctx *gin.Context) {
	collection := request.GetCollectionQuery(ctx)
	pharmacyManagers, err := h.adminUsecase.GetAllPharmacyManager(ctx, &collection)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       response.NewMultiplePharmacyManagerDto(pharmacyManagers),
		Pagination: response.NewPaginationDto(collection),
	})
}

func (h *AdminHandler) GetAllUser(ctx *gin.Context) {
	collection := request.GetCollectionQuery(ctx)
	users, err := h.adminUsecase.GetAllUser(ctx, &collection)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       response.NewMultipleUserDto(users),
		Pagination: response.NewPaginationDto(collection),
	})
}
