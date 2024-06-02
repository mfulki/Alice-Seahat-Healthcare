package request

import (
	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type CreatePartner struct {
	Name         string `json:"name" binding:"required,min=2"`
	Logo         string `json:"logo" binding:"required,url"`
	IsActive     *bool  `json:"is_active" binding:"required,boolean"`
	ManagerName  string `json:"manager_name" binding:"required,min=2"`
	ManagerEmail string `json:"manager_email" binding:"required,email"`
}

type UpdatePartner struct {
	Name        string `json:"name" binding:"required,min=2"`
	Logo        string `json:"logo" binding:"required,url"`
	IsActive    *bool  `json:"is_active" binding:"required,boolean"`
	ManagerName string `json:"manager_name" binding:"required,min=2"`
}

func (req CreatePartner) Partner() entity.Partner {
	return entity.Partner{
		Name:     req.Name,
		IsActive: *req.IsActive,
		Logo:     req.Logo,
		PharmacyManager: entity.PharmacyManager{
			Name:  req.ManagerName,
			Email: req.ManagerEmail,
		},
	}
}

func (req UpdatePartner) Partner() entity.Partner {
	return entity.Partner{
		Name:     req.Name,
		IsActive: *req.IsActive,
		Logo:     req.Logo,
		PharmacyManager: entity.PharmacyManager{
			Name: req.ManagerName,
		},
	}
}
