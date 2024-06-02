package usecase

import (
	"context"
	"mime/multipart"

	"Alice-Seahat-Healthcare/seahat-be/utils"
)

type UploadUsecase interface {
	UploadFile(ctx context.Context, file multipart.File) (string, error)
}

type uploadUsecaseImpl struct {
}

func NewUploadUsecase() *uploadUsecaseImpl {
	return &uploadUsecaseImpl{}
}

func (u *uploadUsecaseImpl) UploadFile(ctx context.Context, file multipart.File) (string, error) {
	uploadUrl, err := utils.UploadCloudinary(ctx, file)
	if err != nil {
		return "", err
	}

	return uploadUrl, nil
}
