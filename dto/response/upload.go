package response

type UploadDto struct {
	FileURL string `json:"file_url"`
}

func NewUploadDto(url string) UploadDto {
	return UploadDto{
		FileURL: url,
	}
}
