package response

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type PrescriptionDto struct {
	ID             uint      `json:"id"`
	TelemedicineID uint      `json:"telemedicine_id"`
	DrugID         uint      `json:"drug_id"`
	Quantity       uint      `json:"quantity"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	Drug           *DrugDto  `json:"drug,omitempty"`
}

func NewPrescriptionDto(p entity.Prescription) PrescriptionDto {
	var drugDto *DrugDto
	if p.Drug.Name != "" {
		drugDto = NewDrugDto(p.Drug)
	}

	return PrescriptionDto{
		ID:             p.ID,
		TelemedicineID: p.TelemedicineID,
		DrugID:         p.DrugID,
		Quantity:       p.Quantity,
		Notes:          p.Notes,
		CreatedAt:      p.CreatedAt,
		Drug:           drugDto,
	}
}

func NewMultiplePrescriptionDto(ps []entity.Prescription) []PrescriptionDto {
	dtos := make([]PrescriptionDto, 0)

	for _, p := range ps {
		dtos = append(dtos, NewPrescriptionDto(p))
	}

	return dtos
}
