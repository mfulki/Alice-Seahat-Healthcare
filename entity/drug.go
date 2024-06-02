package entity

import "time"

type Drug struct {
	ID             uint
	ManufacturerID uint
	Manufacturer   Manufacturer
	Name           string
	GenericName    string
	Composition    string
	Description    string
	Classification string
	Form           string
	UnitInPack     uint
	SellingUnit    string
	Weight         uint
	Height         uint
	Length         uint
	Width          uint
	ImageURL       string
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}
