package entity

import "time"

type Pharmacy struct {
	ID                    uint
	Distance              float64
	ManagerID             uint
	SubdistrictID         uint
	Name                  string
	Latitude              float64
	Longitude             float64
	Location              string
	Address               string
	OpenTime              time.Time
	CloseTime             time.Time
	OperationDay          uint
	PharmacistName        string
	LicenseNumber         string
	PharmacistPhoneNumber string
	Subdistrict           Subdistrict
	ShipmentMethods       []*ShipmentMethod
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             *time.Time
}
