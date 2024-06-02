package request

import (
	"fmt"
	"time"

	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/utils"
)

type AddPharmacy struct {
	ManagerID             uint     `json:"pharmacy_manager_id" binding:"required"`
	SubdistrictID         uint     `json:"subdistrict_id" binding:"required"`
	Name                  string   `json:"pharmacy_name" binding:"required"`
	Latitude              float64  `json:"latitude" binding:"required"`
	Longitude             float64  `json:"longitude" binding:"required"`
	Address               string   `json:"address" binding:"required"`
	OpenTime              string   `json:"open_time" binding:"required,datetime"`
	CloseTime             string   `json:"close_time" binding:"required,datetime"`
	OperationDay          []bool   `json:"operation_day" binding:"required"`
	PharmacistName        string   `json:"pharmacist_name" binding:"required"`
	LicenseNumber         string   `json:"license_number" binding:"required"`
	PharmacistPhoneNumber string   `json:"pharmacist_phone_number" binding:"required"`
	ShipmentMethods       []string `json:"shipment_methods" binding:"gt=0,dive,required,oneof=instant sameday jne tiki pos"`
}

type EditPharmacy struct {
	SubdistrictID         uint     `json:"subdistrict_id" binding:"required"`
	Name                  string   `json:"pharmacy_name" binding:"required"`
	Latitude              float64  `json:"latitude" binding:"required"`
	Longitude             float64  `json:"longitude" binding:"required"`
	Address               string   `json:"address" binding:"required"`
	OpenTime              string   `json:"open_time" binding:"required,datetime"`
	CloseTime             string   `json:"close_time" binding:"required,datetime"`
	OperationDay          []bool   `json:"operation_day" binding:"required"`
	PharmacistName        string   `json:"pharmacist_name" binding:"required"`
	LicenseNumber         string   `json:"license_number" binding:"required"`
	PharmacistPhoneNumber string   `json:"pharmacist_phone_number" binding:"required"`
	ShipmentMethods       []string `json:"shipment_methods" binding:"gt=0,dive,required,oneof=instant sameday jne tiki pos"`
}

func (req *AddPharmacy) Pharmacy() *entity.Pharmacy {
	openTime, _ := time.Parse(constant.FullTimeFormat, req.OpenTime)
	closeTime, _ := time.Parse(constant.FullTimeFormat, req.CloseTime)

	return &entity.Pharmacy{
		ManagerID:             req.ManagerID,
		SubdistrictID:         req.SubdistrictID,
		Name:                  req.Name,
		Latitude:              req.Latitude,
		Longitude:             req.Longitude,
		Location:              fmt.Sprintf("POINT(%g %g)", req.Longitude, req.Latitude),
		Address:               req.Address,
		OpenTime:              openTime,
		CloseTime:             closeTime,
		OperationDay:          utils.OperationDay2Num(req.OperationDay),
		PharmacistName:        req.PharmacistName,
		LicenseNumber:         req.LicenseNumber,
		PharmacistPhoneNumber: req.PharmacistPhoneNumber,
	}
}

func (req *EditPharmacy) Pharmacy() *entity.Pharmacy {
	openTime, _ := time.Parse(constant.FullTimeFormat, req.OpenTime)
	closeTime, _ := time.Parse(constant.FullTimeFormat, req.CloseTime)

	return &entity.Pharmacy{
		SubdistrictID:         req.SubdistrictID,
		Name:                  req.Name,
		Latitude:              req.Latitude,
		Longitude:             req.Longitude,
		Location:              fmt.Sprintf("POINT(%g %g)", req.Longitude, req.Latitude),
		Address:               req.Address,
		OpenTime:              openTime,
		CloseTime:             closeTime,
		OperationDay:          utils.OperationDay2Num(req.OperationDay),
		PharmacistName:        req.PharmacistName,
		LicenseNumber:         req.LicenseNumber,
		PharmacistPhoneNumber: req.PharmacistPhoneNumber,
	}
}
