package utils

import (
	"context"

	"Alice-Seahat-Healthcare/seahat-be/config"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/sirupsen/logrus"
)

func UploadCloudinary(ctx context.Context, file interface{}) (string, error) {
	cld, err := cloudinary.NewFromParams(config.Cloudinary.CloudName, config.Cloudinary.ApiKey, config.Cloudinary.ApiSecret)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	uploadParam, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{Folder: config.Cloudinary.UploadFolder})
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	return uploadParam.SecureURL, nil
}
