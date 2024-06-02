package entity

import "time"

type Doctor struct {
	ID                uint
	SpecializationID  uint
	Specialization    Specialization
	Name              string
	Email             string
	Password          *string
	DateOfBirth       time.Time
	Gender            string
	Certificate       string
	PhotoURL          string
	IsOAuth           bool
	IsVerified        bool
	Price             uint
	Status            string
	YearsOfExperience uint
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}
