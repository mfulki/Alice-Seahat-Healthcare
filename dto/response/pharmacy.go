package response

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/utils"
)

type PharmacyDto struct {
	ID                    uint                `json:"id"`
	ManagerID             uint                `json:"manager_id,omitempty"`
	SubdistrictID         uint                `json:"subdistrict_id,omitempty"`
	Name                  string              `json:"name"`
	Latitude              float64             `json:"latitude"`
	Longitude             float64             `json:"longitude"`
	Address               string              `json:"address,omitempty"`
	OpenTime              string              `json:"open_time"`
	CloseTime             string              `json:"close_time"`
	OperationDay          []bool              `json:"operation_day,omitempty"`
	PharmacistName        string              `json:"pharmacist_name,omitempty"`
	PharmacistPhoneNumber string              `json:"pharmacist_phone_number,omitempty"`
	LicenseNumber         string              `json:"license_number,omitempty"`
	CreatedAt             time.Time           `json:"created_at,omitempty"`
	Shipments             []ShipmentMethodDto `json:"shipments"`
}

type GetPharmacy struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func NewPharmacyDto(pharmacy entity.Pharmacy) (*PharmacyDto, error) {
	longitude, latitude, err := utils.Geo2LongLat(pharmacy.Location)
	if err != nil {
		return nil, err
	}

	openTime := pharmacy.OpenTime.Format(constant.FullTimeFormat)
	closeTime := pharmacy.CloseTime.Format(constant.FullTimeFormat)

	return &PharmacyDto{
		ID:                    pharmacy.ID,
		ManagerID:             pharmacy.ManagerID,
		SubdistrictID:         pharmacy.SubdistrictID,
		Name:                  pharmacy.Name,
		Longitude:             longitude,
		Latitude:              latitude,
		Address:               pharmacy.Address,
		OpenTime:              openTime,
		CloseTime:             closeTime,
		OperationDay:          utils.Num2OperationDay(pharmacy.OperationDay),
		PharmacistName:        pharmacy.PharmacistName,
		PharmacistPhoneNumber: pharmacy.PharmacistPhoneNumber,
		LicenseNumber:         pharmacy.LicenseNumber,
		CreatedAt:             pharmacy.CreatedAt,
		Shipments:             NewMultipleShipmentMethodDto(pharmacy.ShipmentMethods),
	}, nil
}

func NewCreatePharmacyDTO(pharmacy entity.Pharmacy) *PharmacyDto {
	openTime := pharmacy.OpenTime.Format(constant.FullTimeFormat)
	closeTime := pharmacy.CloseTime.Format(constant.FullTimeFormat)

	return &PharmacyDto{
		ID:                    pharmacy.ID,
		ManagerID:             pharmacy.ManagerID,
		SubdistrictID:         pharmacy.SubdistrictID,
		Name:                  pharmacy.Name,
		Longitude:             pharmacy.Longitude,
		Latitude:              pharmacy.Latitude,
		Address:               pharmacy.Address,
		OpenTime:              openTime,
		CloseTime:             closeTime,
		OperationDay:          utils.Num2OperationDay(pharmacy.OperationDay),
		PharmacistName:        pharmacy.PharmacistName,
		PharmacistPhoneNumber: pharmacy.PharmacistPhoneNumber,
		LicenseNumber:         pharmacy.LicenseNumber,
		CreatedAt:             pharmacy.CreatedAt,
		Shipments:             NewMultipleShipmentMethodDto(pharmacy.ShipmentMethods),
	}
}

func NewGetPharmacies(pharmacy entity.Pharmacy) *GetPharmacy {
	return &GetPharmacy{
		ID:   pharmacy.ID,
		Name: pharmacy.Name,
	}
}
