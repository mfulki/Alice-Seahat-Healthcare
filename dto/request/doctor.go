package request

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type DoctorRegister struct {
	SpecializationID  int    `json:"specialization_id" binding:"required,gte=1"`
	Name              string `json:"name" binding:"required,min=2"`
	Email             string `json:"email" binding:"required,email"`
	DateOfBirth       string `json:"date_of_birth" binding:"required,date"`
	Gender            string `json:"gender" binding:"required,oneof=male female"`
	Certificate       string `json:"certificate" binding:"required,url"`
	Price             int    `json:"price" binding:"required,gte=5000"`
	YearsOfExperience int    `json:"years_of_experience" binding:"required,gte=1"`
}

type DoctorRegisterOAuth struct {
	SpecializationID  int    `json:"specialization_id" binding:"required,gte=1"`
	Name              string `json:"name" binding:"required,min=2"`
	DateOfBirth       string `json:"date_of_birth" binding:"required,date"`
	Gender            string `json:"gender" binding:"required,oneof=male female"`
	Certificate       string `json:"certificate" binding:"required,url"`
	Price             int    `json:"price" binding:"required,gte=5000"`
	YearsOfExperience int    `json:"years_of_experience" binding:"required,gte=1"`
	GoogleToken       string `json:"google_token" binding:"required,jwt"`
}

type DoctorLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type DoctorToken struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type DoctorForgot struct {
	Email string `json:"email" binding:"required,email"`
}

type DoctorPersonalEdit struct {
	SpecializationID  int    `json:"specialization_id" binding:"required,gte=1"`
	Name              string `json:"name" binding:"required,min=2"`
	DateOfBirth       string `json:"date_of_birth" binding:"required,date"`
	Gender            string `json:"gender" binding:"required,oneof=male female"`
	Certificate       string `json:"certificate" binding:"required,url"`
	Price             int    `json:"price" binding:"required,gte=5000"`
	YearsOfExperience int    `json:"years_of_experience" binding:"required,gte=1"`
	PhotoURL          string `json:"photo_url" binding:"required,url"`
	Status            string `json:"status" binding:"required,oneof=online offline"`
}

type DoctorStatus struct {
	Status string `json:"status" binding:"required,oneof=online offline"`
}

func (req *DoctorRegister) Doctor() entity.Doctor {
	dateOfBirth, _ := time.Parse(constant.DateFormat, req.DateOfBirth)

	return entity.Doctor{
		SpecializationID:  uint(req.SpecializationID),
		Name:              req.Name,
		Email:             req.Email,
		DateOfBirth:       dateOfBirth,
		Gender:            req.Gender,
		Certificate:       req.Certificate,
		Price:             uint(req.Price),
		YearsOfExperience: uint(req.YearsOfExperience),
	}
}

func (req *DoctorRegisterOAuth) Doctor() entity.Doctor {
	dateOfBirth, _ := time.Parse(constant.DateFormat, req.DateOfBirth)

	return entity.Doctor{
		SpecializationID:  uint(req.SpecializationID),
		Name:              req.Name,
		DateOfBirth:       dateOfBirth,
		Gender:            req.Gender,
		Certificate:       req.Certificate,
		Price:             uint(req.Price),
		YearsOfExperience: uint(req.YearsOfExperience),
	}
}

func (req *DoctorLogin) Doctor() entity.Doctor {
	return entity.Doctor{
		Email:    req.Email,
		Password: &req.Password,
	}
}

func (req *DoctorPersonalEdit) Doctor() entity.Doctor {
	dateOfBirth, _ := time.Parse(constant.DateFormat, req.DateOfBirth)

	return entity.Doctor{
		SpecializationID:  uint(req.SpecializationID),
		Name:              req.Name,
		DateOfBirth:       dateOfBirth,
		Gender:            req.Gender,
		Certificate:       req.Certificate,
		Price:             uint(req.Price),
		YearsOfExperience: uint(req.YearsOfExperience),
		PhotoURL:          req.PhotoURL,
		Status:            req.Status,
	}
}

func (req *DoctorStatus) Doctor() entity.Doctor {
	return entity.Doctor{
		Status: req.Status,
	}
}
