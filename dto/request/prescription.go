package request

import "Alice-Seahat-Healthcare/seahat-be/entity"

type AddPrescriptionItem struct {
	DrugID   int    `json:"drug_id" binding:"required,gte=1"`
	Quantity int    `json:"quantity" binding:"required,gte=1"`
	Notes    string `json:"notes" binding:"required,gte=1"`
}

type AddPrescriptionRequest struct {
	Prescriptions []AddPrescriptionItem `json:"prescriptions" binding:"gt=0,dive"`
}

func (req *AddPrescriptionRequest) Prescription(telemedicineID uint) []entity.Prescription {
	datas := make([]entity.Prescription, 0)

	for _, p := range req.Prescriptions {
		datas = append(datas, entity.Prescription{
			TelemedicineID: telemedicineID,
			DrugID:         uint(p.DrugID),
			Quantity:       uint(p.Quantity),
			Notes:          p.Notes,
		})
	}

	return datas
}
