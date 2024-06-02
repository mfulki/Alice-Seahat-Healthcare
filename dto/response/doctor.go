package response

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type DoctorDto struct {
	ID                uint               `json:"id"`
	SpecializationID  uint               `json:"specialization_id"`
	Name              string             `json:"name"`
	Email             string             `json:"email"`
	DateOfBirth       string             `json:"date_of_birth"`
	Gender            string             `json:"gender"`
	Certificate       string             `json:"certificate"`
	PhotoURL          string             `json:"photo_url"`
	IsOAuth           bool               `json:"is_oauth"`
	IsVerified        bool               `json:"is_verified"`
	Price             uint               `json:"price"`
	Status            string             `json:"status"`
	YearsOfExperience uint               `json:"years_of_experience"`
	CreatedAt         time.Time          `json:"created_at"`
	Specialization    *SpecializationDto `json:"specialization,omitempty"`
}

func NewDoctorDto(doctor entity.Doctor) DoctorDto {
	var specialization *SpecializationDto
	if doctor.Specialization.Name != "" {
		dto := NewSpecializationDto(doctor.Specialization)
		specialization = &dto
	}

	return DoctorDto{
		ID:                doctor.ID,
		SpecializationID:  doctor.SpecializationID,
		Name:              doctor.Name,
		Email:             doctor.Email,
		DateOfBirth:       doctor.DateOfBirth.Format(constant.DateFormat),
		Gender:            doctor.Gender,
		Certificate:       doctor.Certificate,
		PhotoURL:          doctor.PhotoURL,
		IsOAuth:           doctor.IsOAuth,
		IsVerified:        doctor.IsVerified,
		Price:             doctor.Price,
		Status:            doctor.Status,
		YearsOfExperience: doctor.YearsOfExperience,
		CreatedAt:         doctor.CreatedAt,
		Specialization:    specialization,
	}
}

func NewMultipleDoctor(doctors []entity.Doctor) []DoctorDto {
	dtos := make([]DoctorDto, 0)

	for _, d := range doctors {
		dtos = append(dtos, NewDoctorDto(d))
	}

	return dtos
}
