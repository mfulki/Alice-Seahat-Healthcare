package config

type UploadCloudinaryEnv struct {
	CloudName    string
	ApiKey       string
	ApiSecret    string
	UploadFolder string
}

func (e *UploadCloudinaryEnv) loadEnv() error {
	cloudName, err := getEnv("CLOUDINARY_CLOUD_NAME")
	if err != nil {
		return err
	}

	apiKey, err := getEnv("CLOUDINARY_API_KEY")
	if err != nil {
		return err
	}

	apiSecret, err := getEnv("CLOUDINARY_API_SECRET")
	if err != nil {
		return err
	}

	uploadFolder, err := getEnv("CLOUDINARY_UPLOAD_FOLDER")
	if err != nil {
		return err
	}

	e.CloudName = cloudName
	e.ApiKey = apiKey
	e.ApiSecret = apiSecret
	e.UploadFolder = uploadFolder

	return nil
}
