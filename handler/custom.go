package handler

import (
	"Alice-Seahat-Healthcare/seahat-be/apperror"

	"github.com/gin-gonic/gin"
)

type CustomHandler struct {
}

func NewCustomHandler() *CustomHandler {
	return &CustomHandler{}
}

func (h *CustomHandler) NoRoute(ctx *gin.Context) {
	ctx.Error(apperror.NoRoute)
	ctx.Abort()
}
