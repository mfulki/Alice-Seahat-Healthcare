package usecase

import (
	"context"
	"errors"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/repository"
	"Alice-Seahat-Healthcare/seahat-be/utils"
)

type AdminUsecase interface {
	Login(ctx context.Context, body entity.Admin) (string, error)
	GetProfile(ctx context.Context) (*entity.Admin, error)
	GetAllUser(ctx context.Context, clc *entity.Collection) ([]entity.User, error)
	GetAllDoctor(ctx context.Context, clc *entity.Collection) ([]entity.Doctor, error)
	GetAllPharmacyManager(ctx context.Context, clc *entity.Collection) ([]entity.PharmacyManager, error)
}

type adminUsecaseImpl struct {
	userRepository            repository.UserRepository
	doctorRepository          repository.DoctorRepository
	pharmacyManagerRepository repository.PharmacyManagerRepository
	adminRepository           repository.AdminRepository
	transactor                transaction.Transactor
}

func NewAdminUsecase(
	userRepository repository.UserRepository,
	doctorRepository repository.DoctorRepository,
	pharmacyManagerRepository repository.PharmacyManagerRepository,
	adminRepository repository.AdminRepository,
	transactor transaction.Transactor,
) *adminUsecaseImpl {
	return &adminUsecaseImpl{
		userRepository:            userRepository,
		doctorRepository:          doctorRepository,
		pharmacyManagerRepository: pharmacyManagerRepository,
		adminRepository:           adminRepository,
		transactor:                transactor,
	}
}

func (u *adminUsecaseImpl) Login(ctx context.Context, body entity.Admin) (string, error) {
	admin, err := u.adminRepository.SelectOneByEmail(ctx, body.Email)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return "", apperror.InvalidCredential
		}

		return "", err
	}

	if !utils.HashCompareDefault(body.Password, &admin.Password) {
		return "", apperror.InvalidCredential
	}

	jwtData := map[string]any{
		"ID":    admin.ID,
		"Email": admin.Email,
	}

	token, err := utils.JwtGenerateAdmin(jwtData)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *adminUsecaseImpl) GetProfile(ctx context.Context) (*entity.Admin, error) {
	adminCtx, ok := utils.CtxGetAdmin(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}

	admin, err := u.adminRepository.SelectOneByEmail(ctx, adminCtx.Email)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.ResourceNotFound
		}

		return nil, err
	}

	return admin, nil
}

func (u *adminUsecaseImpl) GetAllUser(ctx context.Context, clc *entity.Collection) ([]entity.User, error) {
	return u.userRepository.SelectAll(ctx, clc)
}

func (u *adminUsecaseImpl) GetAllDoctor(ctx context.Context, clc *entity.Collection) ([]entity.Doctor, error) {
	if clc.Sort == "" {
		clc.Sort = "id,desc"
	}

	return u.doctorRepository.SelectAll(ctx, clc)
}

func (u *adminUsecaseImpl) GetAllPharmacyManager(ctx context.Context, clc *entity.Collection) ([]entity.PharmacyManager, error) {
	return u.pharmacyManagerRepository.SelectAll(ctx, clc)
}
