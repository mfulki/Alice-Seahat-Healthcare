package request

import "Alice-Seahat-Healthcare/seahat-be/entity"

type CategoryRequest struct {
	Name string `json:"name" binding:"required,min=2"`
}

func (req CategoryRequest) Category() entity.Category {
	return entity.Category{
		Name: req.Name,
	}
}
