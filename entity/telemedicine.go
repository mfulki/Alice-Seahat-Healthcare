package entity

import "time"

type Telemedicine struct {
	ID                    uint
	User                  User
	Doctor                Doctor
	EndAt                 *time.Time
	Diagnose              *string
	Price                 int
	StartRestAt           *time.Time
	RestDuration          *int
	MedicalCertificateURL *string
	PrescriptionUrl *string
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             *time.Time
	Prescriptions         []Prescription
}
