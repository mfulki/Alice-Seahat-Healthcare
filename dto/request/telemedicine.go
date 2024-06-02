package request

import (
	"time"

	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/entity"
)

type TelemedicineReq struct {
	EndAt          string  `json:"end_at" `
	Diagnose       *string `json:"diagnose" binding:"required,min=4"`
	StartRestAt    string  `json:"start_rest_at"`
	RestDuration   *int    `json:"rest_duration"`
	TelemedicineID uint
}

func (req *TelemedicineReq) Telemedicine() entity.Telemedicine {
	layoutFormat := constant.FullTimeFormat

	startRestAt, _ := time.Parse(layoutFormat, req.StartRestAt)
	endAt, _ := time.Parse(layoutFormat, req.EndAt)
	return entity.Telemedicine{
		ID:           req.TelemedicineID,
		EndAt:        &endAt,
		Diagnose:     req.Diagnose,
		StartRestAt:  &startRestAt,
		RestDuration: req.RestDuration,
	}
}

type AddTelemedicine struct {
	DoctorID uint `json:"doctor_id" binding:"required"`
}

type PutTelemedicine struct {
	EndAt                 string `json:"end_at" binding:"required"`
	Diagnose              string `json:"diagnose" binding:"required"`
	StartRestAt           string `json:"start_rest_at" binding:"required"`
	RestDuration          int    `json:"rest_duration" binding:"required"`
	MedicalCertificateURL string `json:"medical_certificate_url" binding:"required"`
}

func (req *AddTelemedicine) Telemedicine() entity.Telemedicine {
	doctor := entity.Doctor{ID: req.DoctorID}
	return entity.Telemedicine{
		Doctor: doctor,
	}
}

func (req *PutTelemedicine) Telemedicine(id int) entity.Telemedicine {
	layoutFormat := "2006-01-02 15:04:05"
	startRestAt, _ := time.Parse(layoutFormat, req.StartRestAt)
	endAt, _ := time.Parse(layoutFormat, req.EndAt)

	return entity.Telemedicine{
		ID:                    uint(id),
		Diagnose:              &req.Diagnose,
		StartRestAt:           &startRestAt,
		RestDuration:          &req.RestDuration,
		MedicalCertificateURL: &req.MedicalCertificateURL,
		EndAt:                 &endAt,
	}
}
