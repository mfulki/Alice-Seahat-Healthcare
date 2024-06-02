package utils

import (
	"mime/multipart"
	"net/http"

	"Alice-Seahat-Healthcare/seahat-be/constant"
)

func SniffingFile(fileHeader multipart.FileHeader) (string, error) {
	buffer := make([]byte, constant.SniffingFirstResponsibleBytes)
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}
	fileType := http.DetectContentType(buffer)

	return fileType, nil
}
