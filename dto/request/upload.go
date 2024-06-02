package request

import "mime/multipart"

type Upload struct {
	File multipart.FileHeader `form:"file" binding:"required"`
	Type string               `form:"type" binding:"required,oneof=pdf image"`
}
 