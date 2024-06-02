package response

import "Alice-Seahat-Healthcare/seahat-be/entity"

type SpecializationDto struct {
	ID   uint   `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func NewMultipleSpecializationDto(szs []entity.Specialization) []SpecializationDto {
	dtos := make([]SpecializationDto, 0)

	for _, sz := range szs {
		dtos = append(dtos, NewSpecializationDto(sz))
	}

	return dtos
}

func NewSpecializationDto(sz entity.Specialization) SpecializationDto {
	return SpecializationDto{
		ID:   sz.ID,
		Name: sz.Name,
	}
}
