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

type PharmacyManagerUsecase interface {
	Login(ctx context.Context, body entity.PharmacyManager) (string, error)
	GetProfile(ctx context.Context) (*entity.PharmacyManager, error)
}

type pharmacyManagerUsecaseImpl struct {
	pharmacyManagerRepository repository.PharmacyManagerRepository
	partnerRepository         repository.PartnerRepository
	transactor                transaction.Transactor
}

func NewPharmacyManagerUsecase(
	pharmacyManagerRepository repository.PharmacyManagerRepository,
	partnerRepository repository.PartnerRepository,
	transactor transaction.Transactor,
) *pharmacyManagerUsecaseImpl {
	return &pharmacyManagerUsecaseImpl{
		pharmacyManagerRepository: pharmacyManagerRepository,
		partnerRepository:         partnerRepository,
		transactor:                transactor,
	}
}

func (u *pharmacyManagerUsecaseImpl) Login(ctx context.Context, body entity.PharmacyManager) (string, error) {
	pm, err := u.pharmacyManagerRepository.SelectOneByEmail(ctx, body.Email)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return "", apperror.InvalidCredential
		}

		return "", err
	}

	if !utils.HashCompareDefault(body.Password, &pm.Password) {
		return "", apperror.InvalidCredential
	}

	jwtData := map[string]any{
		"ID":    pm.ID,
		"Email": pm.Email,
	}

	token, err := utils.JwtGenerateManager(jwtData)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *pharmacyManagerUsecaseImpl) GetProfile(ctx context.Context) (*entity.PharmacyManager, error) {
	managerCtx, ok := utils.CtxGetManager(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}

	pm, err := u.pharmacyManagerRepository.SelectOneByEmail(ctx, managerCtx.Email)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.ResourceNotFound
		}

		return nil, err
	}

	return pm, nil
}
