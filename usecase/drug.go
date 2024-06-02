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

type DrugUsecase interface {
	GetAll(context.Context, *entity.Collection) ([]*entity.Drug, error)
	InsertOne(ctx context.Context, drug entity.Drug) (*entity.Drug, error)
	UpdateOne(ctx context.Context, drug entity.Drug) (*entity.Drug, error)
	DeleteOne(ctx context.Context, drug entity.Drug) (*entity.Drug, error)
	GetOneById(ctx context.Context, drug entity.Drug) (*entity.Drug, error)
}

type drugUsecaseImpl struct {
	drugRepository repository.DrugRepository
	transactor     transaction.Transactor
}

func NewDrugUsecase(
	drugRepository repository.DrugRepository,
	transactor transaction.Transactor,
) *drugUsecaseImpl {
	return &drugUsecaseImpl{
		drugRepository: drugRepository,
		transactor:     transactor,
	}
}

func (u *drugUsecaseImpl) GetAll(ctx context.Context, clc *entity.Collection) ([]*entity.Drug, error) {
	_, isDoctor := utils.CtxGetDoctor(ctx)

	drugs, err := u.drugRepository.SelectAll(ctx, isDoctor, clc)
	if err != nil {
		return nil, err
	}

	return drugs, nil
}

func (u *drugUsecaseImpl) InsertOne(ctx context.Context, drug entity.Drug) (*entity.Drug, error) {
	checkedDrugs, err := u.drugRepository.CheckNewInsertDrug(ctx, drug)
	if checkedDrugs != nil && err == nil {
		return nil, apperror.DrugExist
	}

	if err != nil {
		if err != apperror.ErrResourceNotFound {
			return nil, err
		}
	}

	drugNew, err := u.drugRepository.InsertOne(ctx, drug)
	if err != nil {
		return nil, err
	}

	return drugNew, nil
}

func (u *drugUsecaseImpl) UpdateOne(ctx context.Context, drug entity.Drug) (*entity.Drug, error) {
	drugUpdate, err := u.drugRepository.UpdateOne(ctx, drug)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.ResourceNotFound
		}

		return nil, err
	}

	return drugUpdate, nil
}

func (u *drugUsecaseImpl) DeleteOne(ctx context.Context, drug entity.Drug) (*entity.Drug, error) {
	drugNew, err := u.drugRepository.UpdateOne(ctx, drug)
	if err != nil {
		return nil, err
	}

	return drugNew, nil
}
func (u *drugUsecaseImpl) GetOneById(ctx context.Context, drug entity.Drug) (*entity.Drug, error) {
	drugItem, err := u.drugRepository.SelectOneById(ctx, drug)
	if err != nil {
		return nil, err
	}
	return drugItem, nil
}
