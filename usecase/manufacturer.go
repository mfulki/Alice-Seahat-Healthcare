package usecase

import (
	"context"

	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/repository"
)

type ManufacturerUsecase interface {
	GetAll(context.Context, *entity.Collection) ([]entity.Manufacturer, error)
}

type manufacturerUsecaseImpl struct {
	manufacturerRepository repository.ManufacturerRepository
}

func NewManufacturerUsecase(
	manufacturerRepository repository.ManufacturerRepository,
) *manufacturerUsecaseImpl {
	return &manufacturerUsecaseImpl{
		manufacturerRepository: manufacturerRepository,
	}
}

func (u *manufacturerUsecaseImpl) GetAll(ctx context.Context, clc *entity.Collection) ([]entity.Manufacturer, error) {
	return u.manufacturerRepository.SelectAll(ctx, clc)
}
