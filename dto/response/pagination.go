package response

import "Alice-Seahat-Healthcare/seahat-be/entity"

type PaginationDto struct {
	CurrentPage  uint `json:"current_page"`
	PerPage      uint `json:"per_page"`
	TotalPages   uint `json:"total_pages"`
	TotalRecords uint `json:"total_records"`
}

func NewPaginationDto(collection entity.Collection) *PaginationDto {
	if collection.Limit == 0 || collection.Page == 0 {
		return nil
	}

	pgnDto := PaginationDto{
		CurrentPage:  collection.Page,
		PerPage:      collection.Limit,
		TotalRecords: collection.TotalRecords,
	}

	pgnDto.TotalPages = pgnDto.TotalRecords / pgnDto.PerPage
	if pgnDto.TotalRecords%pgnDto.PerPage != 0 {
		pgnDto.TotalPages++
	}

	return &pgnDto
}
