package usecase

import (
	"context"
	"errors"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/repository"
)

type CategoryUsecase interface {
	GetAll(ctx context.Context, clc *entity.Collection) ([]entity.Category, error)
	CreateCategory(ctx context.Context, cat entity.Category) (*entity.Category, error)
	GetCategoryByID(ctx context.Context, catID uint) (*entity.Category, error)
	UpdateCategory(ctx context.Context, cat entity.Category) error
	DeleteCategory(ctx context.Context, catID uint) error
}

type categoryUsecaseImpl struct {
	categoryRepository     repository.CategoryRepository
	transactor             transaction.Transactor
	pharmacyDrugrepository repository.PharmacyDrugRepository
}

func NewCategoryUsecase(
	categoryRepository repository.CategoryRepository,
	transactor transaction.Transactor,
	pharmacyDrugrepository repository.PharmacyDrugRepository,
) *categoryUsecaseImpl {
	return &categoryUsecaseImpl{
		categoryRepository:     categoryRepository,
		transactor:             transactor,
		pharmacyDrugrepository: pharmacyDrugrepository,
	}
}

func (u *categoryUsecaseImpl) GetAll(ctx context.Context, clc *entity.Collection) ([]entity.Category, error) {
	return u.categoryRepository.SelectAll(ctx, clc)
}

func (u *categoryUsecaseImpl) CreateCategory(ctx context.Context, cat entity.Category) (*entity.Category, error) {
	_, err := u.categoryRepository.SelectByInsensitiveName(ctx, cat.Name)
	if !errors.Is(err, apperror.ErrResourceNotFound) {
		if err == nil {
			return nil, apperror.CategoryNameExist
		}

		return nil, err
	}

	return u.categoryRepository.InsertOne(ctx, cat)
}

func (u *categoryUsecaseImpl) GetCategoryByID(ctx context.Context, catID uint) (*entity.Category, error) {
	cat, err := u.categoryRepository.SelectByID(ctx, catID)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.ResourceNotFound
		}

		return nil, err
	}

	return cat, nil
}

func (u *categoryUsecaseImpl) UpdateCategory(ctx context.Context, cat entity.Category) error {
	_, err := u.GetCategoryByID(ctx, cat.ID)
	if err != nil {
		return err
	}

	selectedCat, err := u.categoryRepository.SelectByInsensitiveName(ctx, cat.Name)
	if err != nil && !errors.Is(err, apperror.ErrResourceNotFound) {
		return err
	}

	if selectedCat != nil && selectedCat.ID != cat.ID {
		return apperror.CategoryNameExist
	}

	if err := u.categoryRepository.UpdateOne(ctx, cat); err != nil {
		return err
	}

	return nil
}

func (u *categoryUsecaseImpl) DeleteCategory(ctx context.Context, catID uint) error {
	_, err := u.GetCategoryByID(ctx, catID)
	if err != nil {
		return err
	}
	lenPharmacies, err := u.pharmacyDrugrepository.SelectPharmacyDrugsByCategoryId(ctx, catID)
	if err != nil {
		return err
	}
	if lenPharmacies > 0 {
		return apperror.CantDeleteCategory
	}
	if err := u.categoryRepository.DeleteOneByID(ctx, catID); err != nil {
		return err
	}

	return nil
}
