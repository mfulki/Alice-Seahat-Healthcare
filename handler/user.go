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

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	body := new(request.UserLogin)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	token, err := h.userUsecase.Login(ctx, body.User())
	if err != nil {
		ctx.Error(err)
		return
	}

	utils.SendLoginCookie(ctx, token)
	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.LoginPassedMsg,
	})
}

func (h *UserHandler) LoginOAuth(ctx *gin.Context) {
	body := new(request.UserLoginOAuth)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	token, err := h.userUsecase.LoginOAuth(ctx, body.GoogleToken)
	if err != nil {
		ctx.Error(err)
		return
	}

	utils.SendLoginCookie(ctx, token)
	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.LoginPassedMsg,
	})
}

func (h *UserHandler) Register(ctx *gin.Context) {
	body := new(request.UserRegister)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	user, err := h.userUsecase.Register(ctx, body.User())
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response.Body{
		Message: constant.RegisterMsg,
		Data:    response.NewUserDto(*user),
	})
}

func (h *UserHandler) RegisterOAuth(ctx *gin.Context) {
	body := new(request.UserRegisterOAuth)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	user, err := h.userUsecase.RegisterOAuth(ctx, body.User(), body.GoogleToken)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response.Body{
		Message: constant.RegisterMsg,
		Data:    response.NewUserDto(*user),
	})
}

func (h *UserHandler) Verification(ctx *gin.Context) {
	body := new(request.UserToken)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	err := h.userUsecase.Verification(ctx, body.Password, body.Token)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.VerificationMsg,
	})
}

func (h *UserHandler) ForgotPassword(ctx *gin.Context) {
	body := new(request.UserForgot)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.userUsecase.ForgotPassword(ctx, body.Email); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.ForgotTokenMsg,
	})
}

func (h *UserHandler) ResetPassword(ctx *gin.Context) {
	body := new(request.DoctorToken)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	err := h.userUsecase.ResetPassword(ctx, body.Password, body.Token)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.ResetTokenMsg,
	})
}

func (h *UserHandler) Logout(ctx *gin.Context) {
	utils.SendLogoutCookie(ctx)
	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.LogoutMsg,
	})
}

func (h *UserHandler) UpdatePersonal(ctx *gin.Context) {
	body := new(request.UserPersonalEdit)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.userUsecase.UpdateProfile(ctx, body.User()); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataEditMsg,
	})
}

func (h *UserHandler) UpdatePassword(ctx *gin.Context) {
	body := new(request.UserPasswordEdit)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	if err := h.userUsecase.UpdatePassword(ctx, body.OldPassword, body.NewPassword); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataEditMsg,
	})
}

func (h *UserHandler) GetProfile(ctx *gin.Context) {
	user, err := h.userUsecase.GetProfile(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    response.NewUserDto(*user),
	})
}

func (h *UserHandler) ResendVerification(ctx *gin.Context) {
	if err := h.userUsecase.ResendVerification(ctx); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.ForgotTokenMsg,
	})
}
