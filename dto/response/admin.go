package response

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type AdminDto struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAdminDto(admin entity.Admin) AdminDto {
	return AdminDto{
		ID:        admin.ID,
		Name:      admin.Name,
		Email:     admin.Email,
		CreatedAt: admin.CreatedAt,
	}
}

func NewMultipleAdminDto(admins []entity.Admin) []AdminDto {
	dtos := make([]AdminDto, 0)

	for _, admin := range admins {
		dtos = append(dtos, NewAdminDto(admin))
	}

	return dtos
}
