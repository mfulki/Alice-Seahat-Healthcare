package usecase

import (
	"context"

	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/repository"
)

type SpecializationUsecase interface {
	GetAll(ctx context.Context) ([]entity.Specialization, error)
}

type specializationUsecaseImpl struct {
	specializationRepository repository.SpecializationRepository
	transactor               transaction.Transactor
}

func NewSpecializationUsecase(
	specializationRepository repository.SpecializationRepository,
	transactor transaction.Transactor,
) *specializationUsecaseImpl {
	return &specializationUsecaseImpl{
		specializationRepository: specializationRepository,
		transactor:               transactor,
	}
}

func (u *specializationUsecaseImpl) GetAll(ctx context.Context) ([]entity.Specialization, error) {
	return u.specializationRepository.SelectAll(ctx)
}
