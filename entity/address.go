package entity

import "time"

type Address struct {
	ID            uint
	UserID        uint
	SubdistrictID uint
	CityID        uint
	Address       string
	Latitude      float64
	Longitude     float64
	Location      string
	RawLocation   string
	IsMain        bool
	IsActive      bool
	Subdistrict   Subdistrict
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}
