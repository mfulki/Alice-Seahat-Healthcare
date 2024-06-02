package response

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type CategoryDto struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func NewCategoryDto(cat entity.Category) CategoryDto {
	return CategoryDto{
		ID:        cat.ID,
		Name:      cat.Name,
		CreatedAt: cat.CreatedAt,
	}
}

func NewMultipleCategoryDto(cats []entity.Category) []CategoryDto {
	dtos := make([]CategoryDto, 0)

	for _, cat := range cats {
		dtos = append(dtos, NewCategoryDto(cat))
	}

	return dtos
}
