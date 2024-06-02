package response

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type TelemedicineDTO struct {
	ID                    uint              `json:"id"`
	UserID                uint              `json:"user_id"`
	DoctorID              uint              `json:"doctor_id"`
	EndAt                 *time.Time        `json:"end_at"`
	Diagnose              *string           `json:"diagnose"`
	Price                 int               `json:"price"`
	StartRestAt           *time.Time        `json:"start_rest_at"`
	RestDuration          *int              `json:"rest_duration"`
	MedicalCertificateURL *string           `json:"medical_certificate_url"`
	PrescriptionUrl       *string           `json:"prescription_url"`
	Prescriptions         []PrescriptionDto `json:"prescriptions,omitempty"`
	CreatedAt             time.Time         `json:"created_at"`
}

func NewTelemedicineDTO(p entity.Telemedicine) TelemedicineDTO {
	var pDto []PrescriptionDto
	if len(p.Prescriptions) != 0 {
		pDto = NewMultiplePrescriptionDto(p.Prescriptions)
	}

	return TelemedicineDTO{
		ID:                    p.ID,
		UserID:                p.User.ID,
		DoctorID:              p.Doctor.ID,
		EndAt:                 p.EndAt,
		Diagnose:              p.Diagnose,
		Price:                 p.Price,
		StartRestAt:           p.StartRestAt,
		RestDuration:          p.RestDuration,
		MedicalCertificateURL: p.MedicalCertificateURL,
		CreatedAt:             p.CreatedAt,
		PrescriptionUrl:       p.PrescriptionUrl,
		Prescriptions:         pDto,
	}
}

type UserDoctorTelemedicineDTO struct {
	ID                    uint              `json:"id"`
	User                  *UserDto          `json:"user"`
	Doctor                *DoctorDto        `json:"doctor"`
	EndAt                 *time.Time        `json:"end_at"`
	Diagnose              *string           `json:"diagnose"`
	Price                 int               `json:"price"`
	StartRestAt           *time.Time        `json:"start_rest_at"`
	RestDuration          *int              `json:"rest_duration"`
	MedicalCertificateURL *string           `json:"medical_certificate_url"`
	PrescriptionURL       *string           `json:"prescription_url"`
	Prescriptions         []PrescriptionDto `json:"prescriptions"`
	CreatedAt             time.Time         `json:"created_at"`
}

func NewUserDoctorTelemedicineDTO(p entity.Telemedicine) UserDoctorTelemedicineDTO {
	pDto := make([]PrescriptionDto, 0)

	if len(p.Prescriptions) != 0 {
		pDto = NewMultiplePrescriptionDto(p.Prescriptions)
	}

	doctor := NewDoctorDto(p.Doctor)
	user := NewUserDto(p.User)
	return UserDoctorTelemedicineDTO{
		ID:                    p.ID,
		User:                  &user,
		Doctor:                &doctor,
		EndAt:                 p.EndAt,
		Diagnose:              p.Diagnose,
		Price:                 p.Price,
		StartRestAt:           p.StartRestAt,
		RestDuration:          p.RestDuration,
		MedicalCertificateURL: p.MedicalCertificateURL,
		PrescriptionURL:       p.PrescriptionUrl,
		CreatedAt:             p.CreatedAt,
		Prescriptions:         pDto,
	}
}
