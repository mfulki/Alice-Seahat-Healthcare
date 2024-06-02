package usecase

import (
	"bytes"
	"context"
	"errors"
	"time"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/repository"
	"Alice-Seahat-Healthcare/seahat-be/utils"
)

type TelemedicineUsecase interface {
	CreateTelemedicine(ctx context.Context, telemedicine entity.Telemedicine) (*entity.Telemedicine, error)
	GetTelemedicineByID(ctx context.Context, id uint) (*entity.Telemedicine, error)
	GetAllTelemedicine(ctx context.Context, clc *entity.Collection) ([]entity.Telemedicine, error)
	UpdateOne(ctx context.Context, updateTelemedicine entity.Telemedicine) error
	AddManyPrescriptedDrugs(ctx context.Context, prescriptions []entity.Prescription) (*string, error)
	UpdateOneAndCreateMedicalCertificate(ctx context.Context, updateTelemedicine entity.Telemedicine) (*string, error)
}

type telemedicineUsecaseImpl struct {
	telemedicineRepository repository.TelemedicineRepository
	userRepository         repository.UserRepository
	doctorRepository       repository.DoctorRepository
	prescriptionRepository repository.PrescriptionRepository
}

func NewTelemedicineUsecase(
	telemedicineRepository repository.TelemedicineRepository,
	userRepository repository.UserRepository,
	doctorRepository repository.DoctorRepository,
	prescriptionRepository repository.PrescriptionRepository,
) *telemedicineUsecaseImpl {
	return &telemedicineUsecaseImpl{
		telemedicineRepository: telemedicineRepository,
		userRepository:         userRepository,
		doctorRepository:       doctorRepository,
		prescriptionRepository: prescriptionRepository,
	}
}

func (u *telemedicineUsecaseImpl) CreateTelemedicine(ctx context.Context, telemedicine entity.Telemedicine) (*entity.Telemedicine, error) {
	userCtx, ok := utils.CtxGetUser(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}

	doctor, err := u.doctorRepository.SelectOneByID(ctx, telemedicine.Doctor.ID)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.DoctorNotExist
		}

		return nil, err
	}

	_, err = u.telemedicineRepository.SelectOneOngoingByUserAndDoctorID(ctx, userCtx.ID, telemedicine.Doctor.ID)
	if !errors.Is(err, apperror.ErrResourceNotFound) {
		if err == nil {
			return nil, apperror.TelemedicineOngoingWithDoctor
		}

		return nil, err
	}

	telemedicine.User.ID = userCtx.ID
	telemedicine.Price = int(doctor.Price)
	telemedicineData, err := u.telemedicineRepository.InsertOne(ctx, telemedicine)
	if err != nil {
		return nil, err
	}

	return telemedicineData, nil
}

func (u *telemedicineUsecaseImpl) GetTelemedicineByID(ctx context.Context, id uint) (*entity.Telemedicine, error) {
	user, _ := utils.CtxGetUser(ctx)
	doctor, _ := utils.CtxGetDoctor(ctx)

	var userID *uint
	var doctorID *uint

	if user != nil {
		userID = &user.ID
	} else {
		doctorID = &doctor.ID
	}

	telemedicineData, err := u.telemedicineRepository.SelectOneByID(ctx, id, userID, doctorID)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.ResourceNotFound
		}

		return nil, err
	}

	user, err = u.userRepository.SelectOneByID(ctx, telemedicineData.User.ID)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.ResourceNotFound
		}

		return nil, err
	}
	telemedicineData.User = *user
	doctor, err = u.doctorRepository.SelectOneWithSpecializationByID(ctx, telemedicineData.Doctor.ID)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.ResourceNotFound
		}

		return nil, err
	}
	telemedicineData.Doctor = *doctor
	prescriptions, err := u.prescriptionRepository.GetAllByTelemedicineID(ctx, telemedicineData.ID)
	if err != nil {
		return nil, err
	}

	telemedicineData.Prescriptions = prescriptions

	return telemedicineData, nil
}

func (u *telemedicineUsecaseImpl) GetAllTelemedicine(ctx context.Context, clc *entity.Collection) ([]entity.Telemedicine, error) {
	user, _ := utils.CtxGetUser(ctx)
	doctor, _ := utils.CtxGetDoctor(ctx)

	var userID *uint
	var doctorID *uint

	if user != nil {
		userID = &user.ID
	} else {
		doctorID = &doctor.ID
	}

	telemedicines, err := u.telemedicineRepository.GetAllTelemedicine(ctx, userID, doctorID, clc)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(telemedicines); i++ {
		user, err := u.userRepository.SelectOneByID(ctx, telemedicines[i].User.ID)
		if err != nil {
			return nil, err
		}
		telemedicines[i].User = *user
		doctor, err := u.doctorRepository.SelectOneWithSpecializationByID(ctx, telemedicines[i].Doctor.ID)
		if err != nil {
			return nil, err
		}
		telemedicines[i].Doctor = *doctor
	}

	return telemedicines, nil
}

func (u *telemedicineUsecaseImpl) UpdateOne(ctx context.Context, updateTelemedicine entity.Telemedicine) error {
	doctor, ok := utils.CtxGetDoctor(ctx)
	if !ok {
		return apperror.ErrInternalServer
	}

	updateTelemedicine.Doctor.ID = doctor.ID

	err := u.telemedicineRepository.UpdateOne(ctx, updateTelemedicine)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return apperror.ResourceNotFound
		}

		return err
	}

	return nil
}

func (u *telemedicineUsecaseImpl) UpdateOneAndCreateMedicalCertificate(ctx context.Context, updateTelemedicine entity.Telemedicine) (*string, error) {
	doctor, ok := utils.CtxGetDoctor(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}

	updateTelemedicine.Doctor.ID = doctor.ID

	zeroduration := 0
	if updateTelemedicine.RestDuration == &zeroduration {
		updateTelemedicine.StartRestAt = &time.Time{}
	}

	if *updateTelemedicine.Diagnose == constant.DiagnoseStatus {
		s := constant.EndChat
		err := u.UpdateOne(ctx, updateTelemedicine)
		if err != nil {
			return nil, err
		}
		return &s, nil
	}

	telemedicine, err := u.telemedicineRepository.SelectOneByIdJoinDoctorUser(ctx, updateTelemedicine)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.ResourceNotFound
		}

		return nil, err
	}
	telemedicine.EndAt = updateTelemedicine.EndAt
	pdf := bytes.NewBuffer(nil)
	err = utils.GenerateMedicalCertificate(pdf, *telemedicine)
	if err != nil {
		return nil, err
	}

	url, err := utils.UploadCloudinary(ctx, pdf)
	if err != nil {
		return nil, err
	}
	telemedicine.MedicalCertificateURL = &url
	err = u.UpdateOne(ctx, *telemedicine)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (u *telemedicineUsecaseImpl) AddManyPrescriptedDrugs(ctx context.Context, prescriptions []entity.Prescription) (*string, error) {
	doctor, ok := utils.CtxGetDoctor(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}

	telemedicineID := prescriptions[0].TelemedicineID
	bodyTelemedicine := entity.Telemedicine{ID: telemedicineID, Doctor: entity.Doctor{ID: doctor.ID}}
	telemedicine, err := u.telemedicineRepository.SelectOneByIdJoinDoctorUser(ctx, bodyTelemedicine)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.ResourceNotFound
		}

		return nil, err
	}

	if telemedicine.EndAt != nil {
		return nil, apperror.TelemedicineHasBeenEnded
	}

	_, err = u.prescriptionRepository.InsertMany(ctx, prescriptions)
	if err != nil {
		return nil, err
	}

	prescriptions, err = u.prescriptionRepository.GetAllByTelemedicineID(ctx, telemedicine.ID)
	if err != nil {
		return nil, err
	}
	telemedicine.Prescriptions = prescriptions

	pdf := bytes.NewBuffer(nil)
	err = utils.GeneratePrescription(pdf, *telemedicine)
	if err != nil {
		return nil, err
	}

	url, err := utils.UploadCloudinary(ctx, pdf)
	if err != nil {
		return nil, err
	}
	telemedicine.PrescriptionUrl = &url
	err = u.telemedicineRepository.UpdateOne(ctx, *telemedicine)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.ResourceNotFound
		}

		return nil, err
	}
	return &url, nil
}
