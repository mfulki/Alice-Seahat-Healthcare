package handler

import (
	"net/http"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/dto/request"
	"Alice-Seahat-Healthcare/seahat-be/dto/response"
	"Alice-Seahat-Healthcare/seahat-be/usecase"
	"Alice-Seahat-Healthcare/seahat-be/utils"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadUsecase usecase.UploadUsecase
}

func NewUploadHandler(uploadUsecase usecase.UploadUsecase) *UploadHandler {
	return &UploadHandler{
		uploadUsecase: uploadUsecase,
	}
}

func (h *UploadHandler) UploadCloudinary(ctx *gin.Context) {
	body := new(request.Upload)
	if err := ctx.ShouldBind(body); err != nil {
		ctx.Error(err)
		return
	}

	if int(body.File.Size) > constant.MaxFileSize[body.Type] {
		ctx.Error(apperror.FileTooLarge)
		return
	}

	fileType, err := utils.SniffingFile(body.File)
	if err != nil {
		ctx.Error(err)
		return
	}

	if fileType != constant.FileType[body.Type] {
		ctx.Error(apperror.FileInvalidType)
		return
	}

	file, err := body.File.Open()
	if err != nil {
		ctx.Error(err)
		return
	}

	defer file.Close()

	uploadUrl, err := h.uploadUsecase.UploadFile(ctx, file)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response.Body{
		Message: constant.DataCreatedMsg,
		Data:    response.NewUploadDto(uploadUrl),
	})
}
