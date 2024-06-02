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

type PharmacyUsecase interface {
	CreatePharmacy(ctx context.Context, pharmacy entity.Pharmacy, shipments []string) (*entity.Pharmacy, error)
	UpdatePharmacy(ctx context.Context, pharmacy entity.Pharmacy, id uint, shipments []string) (*entity.Pharmacy, error)
	GetAllPharmacies(ctx context.Context, clc *entity.Collection) ([]entity.Pharmacy, error)
	GetAllPharmaciesByManagerID(ctx context.Context, managerID uint, clc *entity.Collection) ([]entity.Pharmacy, error)
	GetPharmacyByID(ctx context.Context, pharmacyID uint) (*entity.Pharmacy, error)
}

type pharmacyUsecaseImpl struct {
	pharmacyRepository       repository.PharmacyRepository
	shipmentMethodRepository repository.ShipmentMethodRepository
	transactor               transaction.Transactor
}

func NewPharmacyUsecase(
	pharmacyRepository repository.PharmacyRepository,
	shipmentMethodRepository repository.ShipmentMethodRepository,
	transactor transaction.Transactor,
) *pharmacyUsecaseImpl {
	return &pharmacyUsecaseImpl{
		pharmacyRepository:       pharmacyRepository,
		shipmentMethodRepository: shipmentMethodRepository,
		transactor:               transactor,
	}
}

func (u *pharmacyUsecaseImpl) CreatePharmacy(ctx context.Context, pharmacy entity.Pharmacy, shipments []string) (*entity.Pharmacy, error) {
	pharmacyData, err := u.pharmacyRepository.InsertOne(ctx, pharmacy)
	if err != nil {
		return nil, err
	}

	if err := u.shipmentMethodRepository.InsertManyPharmacyShipment(ctx, pharmacyData.ID, shipments); err != nil {
		return nil, err
	}

	return pharmacyData, nil
}

func (u *pharmacyUsecaseImpl) UpdatePharmacy(ctx context.Context, pharmacy entity.Pharmacy, id uint, shipments []string) (*entity.Pharmacy, error) {
	manager, ok := utils.CtxGetManager(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}

	pharmacy.ManagerID = manager.ID
	pharmacyData, err := u.pharmacyRepository.UpdateOne(ctx, pharmacy, id)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.ResourceNotFound
		}

		return nil, err
	}

	if err := u.shipmentMethodRepository.HardDeleteAllShipmentByPharmacy(ctx, pharmacyData.ID); err != nil {
		return nil, err
	}

	if err := u.shipmentMethodRepository.InsertManyPharmacyShipment(ctx, pharmacyData.ID, shipments); err != nil {
		return nil, err
	}

	return pharmacyData, nil
}

func (u *pharmacyUsecaseImpl) GetAllPharmacies(ctx context.Context, clc *entity.Collection) ([]entity.Pharmacy, error) {
	managerID := uint(0)
	manager, _ := utils.CtxGetManager(ctx)
	if manager != nil {
		managerID = manager.ID
	}

	pharmacyData, err := u.pharmacyRepository.GetAllPharmacies(ctx, managerID, clc)
	if err != nil {
		return nil, err
	}

	return pharmacyData, nil
}

func (u *pharmacyUsecaseImpl) GetAllPharmaciesByManagerID(ctx context.Context, managerID uint, clc *entity.Collection) ([]entity.Pharmacy, error) {
	return u.pharmacyRepository.GetAllPharmacies(ctx, managerID, clc)
}

func (u *pharmacyUsecaseImpl) GetPharmacyByID(ctx context.Context, pharmacyID uint) (*entity.Pharmacy, error) {
	managerID := uint(0)
	manager, _ := utils.CtxGetManager(ctx)
	if manager != nil {
		managerID = manager.ID
	}

	p, err := u.pharmacyRepository.GetPharmaciesByIDWithShipments(ctx, entity.Pharmacy{
		ID:        pharmacyID,
		ManagerID: managerID,
	})

	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.ResourceNotFound
		}

		return nil, err
	}

	return p, nil
}
