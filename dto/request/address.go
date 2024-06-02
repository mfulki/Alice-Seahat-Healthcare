package request

import "Alice-Seahat-Healthcare/seahat-be/entity"

type AddressRequest struct {
	SubdistrictID int     `json:"subdistrict_id" binding:"required,gte=1"`
	Address       string  `json:"address" binding:"required"`
	Latitude      float64 `json:"latitude" binding:"required,latitude"`
	Longitude     float64 `json:"longitude" binding:"required,longitude"`
	IsMain        *bool   `json:"is_main" binding:"required,boolean"`
	IsActive      *bool   `json:"is_active" binding:"required,boolean"`
}

func (req *AddressRequest) Addr() *entity.Address {
	if *req.IsMain && !(*req.IsActive) {
		return nil
	}

	return &entity.Address{
		SubdistrictID: uint(req.SubdistrictID),
		Address:       req.Address,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		IsMain:        *req.IsMain,
		IsActive:      *req.IsActive,
	}
}
