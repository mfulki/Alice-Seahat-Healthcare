package handler

import (
	"net/http"

	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/dto/request"
	"Alice-Seahat-Healthcare/seahat-be/dto/response"
	"Alice-Seahat-Healthcare/seahat-be/usecase"

	"github.com/gin-gonic/gin"
)

type MessageBubbleHandler struct {
	messageBubbleUseCase usecase.MessageBubbleUsecase
}

func NewMessageBubbleHandler(messageBubbleUsecase usecase.MessageBubbleUsecase) *MessageBubbleHandler {
	return &MessageBubbleHandler{
		messageBubbleUseCase: messageBubbleUsecase,
	}
}

func (h *MessageBubbleHandler) AddChat(ctx *gin.Context) {
	body := new(request.MessageBubble)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	data, err := h.messageBubbleUseCase.AddMessageBubble(ctx, body.MessageBubble())
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataCreatedMsg,
		Data:    response.NewMessageBubbleDTO(*data),
	})
}
