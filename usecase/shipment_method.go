package usecase

import (
	"context"

	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/repository"
)

type ShipmentMethodUsecase interface {
	GetAllShipmentMethod(ctx context.Context) ([]*entity.ShipmentMethod, error)
}

type shipmentMethodUsecaseImpl struct {
	shipmentMethodRepository repository.ShipmentMethodRepository
	transactor               transaction.Transactor
}

func NewShipmentMethodUsecase(
	shipmentMethodRepository repository.ShipmentMethodRepository,
	transactor transaction.Transactor,
) *shipmentMethodUsecaseImpl {
	return &shipmentMethodUsecaseImpl{
		shipmentMethodRepository: shipmentMethodRepository,
		transactor:               transactor,
	}
}

func (u *shipmentMethodUsecaseImpl) GetAllShipmentMethod(ctx context.Context) ([]*entity.ShipmentMethod, error) {
	return u.shipmentMethodRepository.SelectAllShipmentMethod(ctx)
}
