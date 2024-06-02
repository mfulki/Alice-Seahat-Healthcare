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

type CartItemUsecase interface {
	GetAllCartItem(ctx context.Context) ([]entity.CartItem, error)
	AddCartItem(ctx context.Context, item entity.CartItem) (*entity.CartItem, error)
	UpdateCartItem(ctx context.Context, item entity.CartItem) error
	DeleteBulkCartItem(ctx context.Context, ids []uint) error
}

type cartItemUsecaseImpl struct {
	cartItemRepository     repository.CartItemRepository
	pharmacyDrugRepository repository.PharmacyDrugRepository
	transactor             transaction.Transactor
}

func NewCartItemUsecase(
	cartItemRepository repository.CartItemRepository,
	pharmacyDrugRepository repository.PharmacyDrugRepository,
	transactor transaction.Transactor,
) *cartItemUsecaseImpl {
	return &cartItemUsecaseImpl{
		cartItemRepository:     cartItemRepository,
		pharmacyDrugRepository: pharmacyDrugRepository,
		transactor:             transactor,
	}
}

func (u *cartItemUsecaseImpl) GetAllCartItem(ctx context.Context) ([]entity.CartItem, error) {
	user, ok := utils.CtxGetUser(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}

	return u.cartItemRepository.GetAllWithOtherDetail(ctx, user.ID)
}

func (u *cartItemUsecaseImpl) AddCartItem(ctx context.Context, item entity.CartItem) (*entity.CartItem, error) {
	user, ok := utils.CtxGetUser(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}

	pd, err := u.pharmacyDrugRepository.SelectOneByID(ctx, item.PharmacyDrugID)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.PharmacyDrugNotExist
		}

		return nil, err
	}

	item.UserID = user.ID

	cartItem, err := u.cartItemRepository.GetByPrescriptedAndPharmacyDrugID(ctx, item)
	if err != nil && !errors.Is(err, apperror.ErrResourceNotFound) {
		return nil, err
	}

	if cartItem != nil {
		if cartItem.IsPrescripted {
			return nil, apperror.PrescriptedExist
		}

		cartItem.Quantity += item.Quantity
		if cartItem.Quantity > uint(pd.Stock) {
			return nil, apperror.InsufficientStock
		}

		if err := u.cartItemRepository.UpdateQuantityByID(ctx, *cartItem); err != nil {
			return nil, err
		}

		return cartItem, nil
	}

	if item.Quantity > uint(pd.Stock) {
		return nil, apperror.InsufficientStock
	}

	itemAdded, err := u.cartItemRepository.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	return itemAdded, nil
}

func (u *cartItemUsecaseImpl) UpdateCartItem(ctx context.Context, item entity.CartItem) error {
	user, ok := utils.CtxGetUser(ctx)
	if !ok {
		return apperror.ErrInternalServer
	}

	item.UserID = user.ID

	cartItem, err := u.cartItemRepository.SelectOneByID(ctx, item)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return apperror.ResourceNotFound
		}

		return err
	}

	if cartItem.IsPrescripted {
		return apperror.PharmacyDrugCantEdit
	}

	pd, err := u.pharmacyDrugRepository.SelectOneByID(ctx, cartItem.PharmacyDrugID)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return apperror.PharmacyDrugNotExist
		}

		return err
	}

	if item.Quantity > uint(pd.Stock) {
		return apperror.InsufficientStock
	}

	if err := u.cartItemRepository.UpdateQuantityByID(ctx, item); err != nil {
		return err
	}

	return nil
}

func (u *cartItemUsecaseImpl) DeleteBulkCartItem(ctx context.Context, ids []uint) error {
	user, ok := utils.CtxGetUser(ctx)
	if !ok {
		return apperror.ErrInternalServer
	}

	if err := u.cartItemRepository.DeleteManyByID(ctx, ids, user.ID); err != nil {
		return err
	}

	return nil
}
