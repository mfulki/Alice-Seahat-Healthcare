package request

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type UserRegister struct {
	Name        string `json:"name" binding:"required,min=2"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	DateOfBirth string `json:"date_of_birth" binding:"required,date"`
	Gender      string `json:"gender" binding:"required,oneof=male female"`
}
type UserRegisterOAuth struct {
	Name        string `json:"name" binding:"required,min=2"`
	DateOfBirth string `json:"date_of_birth" binding:"required,date"`
	Gender      string `json:"gender" binding:"required,oneof=male female"`
	GoogleToken string `json:"google_token" binding:"required,jwt"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserLoginOAuth struct {
	GoogleToken string `json:"google_token" binding:"required,jwt"`
}

type UserToken struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserForgot struct {
	Email string `json:"email" binding:"required,email"`
}

type UserPersonalEdit struct {
	Name        string `json:"name" binding:"required,min=2"`
	DateOfBirth string `json:"date_of_birth" binding:"required,date"`
	Gender      string `json:"gender" binding:"required,oneof=male female"`
	PhotoURL    string `json:"photo_url" binding:"required,url"`
}

type UserPasswordEdit struct {
	OldPassword string `json:"old_password" binding:"required,min=8"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

func (req *UserRegister) User() entity.User {
	dateOfBirth, _ := time.Parse(constant.DateFormat, req.DateOfBirth)

	return entity.User{
		Name:        req.Name,
		Email:       req.Email,
		Password:    &req.Password,
		DateOfBirth: dateOfBirth,
		Gender:      req.Gender,
	}
}

func (req *UserRegisterOAuth) User() entity.User {
	dateOfBirth, _ := time.Parse(constant.DateFormat, req.DateOfBirth)

	return entity.User{
		Name:        req.Name,
		DateOfBirth: dateOfBirth,
		Gender:      req.Gender,
	}
}

func (req *UserLogin) User() entity.User {
	return entity.User{
		Email:    req.Email,
		Password: &req.Password,
	}
}

func (req *UserPersonalEdit) User() entity.User {
	dateOfBirth, _ := time.Parse(constant.DateFormat, req.DateOfBirth)

	return entity.User{
		Name:        req.Name,
		DateOfBirth: dateOfBirth,
		Gender:      req.Gender,
		PhotoURL:    req.PhotoURL,
	}
}

func (req *UserLogin) PharmacyManager() entity.PharmacyManager {
	return entity.PharmacyManager{
		Email:    req.Email,
		Password: req.Password,
	}
}

func (req *UserLogin) Admin() entity.Admin {
	return entity.Admin{
		Email:    req.Email,
		Password: req.Password,
	}
}
