package response

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type UserDto struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	DateOfBirth string    `json:"date_of_birth"`
	Gender      string    `json:"gender"`
	PhotoURL    string    `json:"photo_url"`
	IsOAuth     bool      `json:"is_oauth"`
	IsVerified  bool      `json:"is_verified"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewUserDto(user entity.User) UserDto {
	return UserDto{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		DateOfBirth: user.DateOfBirth.Format(constant.DateFormat),
		Gender:      user.Gender,
		PhotoURL:    user.PhotoURL,
		IsOAuth:     user.IsOAuth,
		IsVerified:  user.IsVerified,
		CreatedAt:   user.CreatedAt,
	}
}

func NewMultipleUserDto(users []entity.User) []UserDto {
	dtos := make([]UserDto, 0)

	for _, user := range users {
		dtos = append(dtos, NewUserDto(user))
	}

	return dtos
}
