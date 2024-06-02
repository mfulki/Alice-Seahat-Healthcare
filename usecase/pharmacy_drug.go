package usecase

import (
	"context"
	"errors"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/repository"
	"Alice-Seahat-Healthcare/seahat-be/utils"
)

type PharmacyDrugUsecase interface {
	GetAllWithinRadius(context.Context, entity.Address, *entity.Collection) ([]*entity.PharmacyDrug, error)
	GetByID(context.Context, uint) (*entity.PharmacyDrug, error)
	CreatePharmacyDrug(ctx context.Context, pharmacyDrug entity.PharmacyDrug) (*entity.PharmacyDrug, error)
	UpdatePharmacyDrug(ctx context.Context, pharmacyDrug entity.PharmacyDrug, id uint) error
	GetAllPharmacyDrug(ctx context.Context, clc *entity.Collection, pharmacyID uint) ([]entity.PharmacyDrug, error)
	GetNearestPharmacies(ctx context.Context, drugID uint, addr entity.Address, clc *entity.Collection) ([]entity.PharmacyDrug, error)
	GetAllDrugsOfTheDay(ctx context.Context) ([]entity.PharmacyDrug, error)
}

type pharmacyDrugUsecaseImpl struct {
	pharmacyDrugRepository repository.PharmacyDrugRepository
	addressRepository      repository.AddressRepository
	drugRepository         repository.DrugRepository
	pharmacyRepository     repository.PharmacyRepository
	categoryRepository     repository.CategoryRepository
	transactor             transaction.Transactor
	stockJournalRepository repository.StockJournalRepository
}

func NewPharmacyDrugUsecase(
	pharmacyDrugRepository repository.PharmacyDrugRepository,
	addressRepository repository.AddressRepository,
	drugRepository repository.DrugRepository,
	pharmacyRepository repository.PharmacyRepository,
	categoryRepository repository.CategoryRepository,
	transactor transaction.Transactor,
	stockJournalRepository repository.StockJournalRepository,
) *pharmacyDrugUsecaseImpl {
	return &pharmacyDrugUsecaseImpl{
		pharmacyDrugRepository: pharmacyDrugRepository,
		addressRepository:      addressRepository,
		drugRepository:         drugRepository,
		pharmacyRepository:     pharmacyRepository,
		categoryRepository:     categoryRepository,
		transactor:             transactor,
		stockJournalRepository: stockJournalRepository,
	}
}

func (u *pharmacyDrugUsecaseImpl) GetAllWithinRadius(ctx context.Context, addr entity.Address, clc *entity.Collection) ([]*entity.PharmacyDrug, error) {
	if addr.UserID != 0 {
		if mainAddr, _ := u.addressRepository.GetMainAddress(ctx, addr.UserID); mainAddr != nil {
			addr.Latitude = mainAddr.Latitude
			addr.Longitude = mainAddr.Longitude
		}
	}

	if addr.Latitude == 0 && addr.Longitude == 0 {
		addr.Latitude = constant.DefaultLatitude
		addr.Longitude = constant.DefaultLongitude
	}

	pharmacyDrugs, err := u.pharmacyDrugRepository.SelectAllWithinRadius(ctx, addr, constant.SearchRadiusMetre, clc)
	if err != nil {
		return nil, err
	}

	return pharmacyDrugs, nil
}

func (u *pharmacyDrugUsecaseImpl) GetByID(ctx context.Context, id uint) (*entity.PharmacyDrug, error) {
	pharmacyDrug, err := u.pharmacyDrugRepository.SelectById(ctx, id)
	if err != nil {
		return nil, err
	}
	return pharmacyDrug, nil
}

func (u *pharmacyDrugUsecaseImpl) GetAllPharmacyDrugs(ctx context.Context, id uint) (*entity.PharmacyDrug, error) {
	pharmacyDrug, err := u.pharmacyDrugRepository.SelectById(ctx, id)
	if err != nil {
		return nil, err
	}
	return pharmacyDrug, nil
}

func (u *pharmacyDrugUsecaseImpl) CreatePharmacyDrug(ctx context.Context, pharmacyDrug entity.PharmacyDrug) (*entity.PharmacyDrug, error) {
	managerCtx, ok := utils.CtxGetManager(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}

	_, err := u.drugRepository.SelectOneById(ctx, entity.Drug{ID: pharmacyDrug.DrugID})
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.DrugNotExist
		}

		return nil, err
	}

	_, err = u.pharmacyRepository.SelectByIDAndManagerID(ctx, pharmacyDrug.PharmacyID, managerCtx.ID)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.PharmacyNotExist
		}

		return nil, err
	}

	_, err = u.categoryRepository.SelectByID(ctx, pharmacyDrug.CategoryID)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.CategoryNotExist
		}

		return nil, err
	}

	_, err = u.pharmacyDrugRepository.SelectOneByPharmacyAndDrugID(ctx, pharmacyDrug)
	if !errors.Is(err, apperror.ErrResourceNotFound) {
		if err == nil {
			return nil, apperror.PharmacyDrugExist
		}

		return nil, err
	}

	pharmacyDrugData, err := u.pharmacyDrugRepository.CreateOne(ctx, pharmacyDrug)
	if err != nil {
		return nil, err
	}

	return pharmacyDrugData, nil
}

func (u *pharmacyDrugUsecaseImpl) UpdatePharmacyDrug(ctx context.Context, pharmacyDrug entity.PharmacyDrug, id uint) error {
	managerCtx, ok := utils.CtxGetManager(ctx)
	if !ok {
		return apperror.ErrInternalServer
	}

	_, err := u.drugRepository.SelectOneById(ctx, entity.Drug{ID: pharmacyDrug.DrugID})
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return apperror.DrugNotExist
		}

		return err
	}

	_, err = u.pharmacyRepository.SelectByIDAndManagerID(ctx, pharmacyDrug.PharmacyID, managerCtx.ID)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return apperror.PharmacyNotExist
		}

		return err
	}

	_, err = u.categoryRepository.SelectByID(ctx, pharmacyDrug.CategoryID)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return apperror.CategoryNotExist
		}

		return err
	}

	selectedPharmacyDrug, err := u.pharmacyDrugRepository.SelectOneByPharmacyAndDrugID(ctx, pharmacyDrug)
	if err != nil && !errors.Is(err, apperror.ErrResourceNotFound) {
		return err
	}

	if selectedPharmacyDrug != nil && selectedPharmacyDrug.ID != id {
		return apperror.PharmacyDrugExist
	}

	quantity := int(pharmacyDrug.Stock - selectedPharmacyDrug.Stock)
	if quantity != 0 {
		stockJournals := []entity.StockJurnal{{DrugId: pharmacyDrug.DrugID, PharmacyId: pharmacyDrug.PharmacyID, Description: constant.UpdatedStock, Quantity: quantity}}
		err := u.stockJournalRepository.InsertStockJournal(ctx, stockJournals)
		if err != nil {
			return err
		}
	}
	err = u.pharmacyDrugRepository.UpdateOne(ctx, pharmacyDrug, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *pharmacyDrugUsecaseImpl) GetAllPharmacyDrug(ctx context.Context, clc *entity.Collection, pharmacyID uint) ([]entity.PharmacyDrug, error) {
	managerID := uint(0)
	managerCtx, _ := utils.CtxGetManager(ctx)
	if managerCtx != nil {
		managerID = managerCtx.ID
	}

	pharmacyDrugs, err := u.pharmacyDrugRepository.SelectAll(ctx, pharmacyID, managerID, clc)
	if err != nil {
		return nil, err
	}

	return pharmacyDrugs, nil
}

func (u *pharmacyDrugUsecaseImpl) GetNearestPharmacies(ctx context.Context, drugID uint, addr entity.Address, clc *entity.Collection) ([]entity.PharmacyDrug, error) {
	if addr.UserID != 0 {
		if mainAddr, _ := u.addressRepository.GetMainAddress(ctx, addr.UserID); mainAddr != nil {
			addr.Latitude = mainAddr.Latitude
			addr.Longitude = mainAddr.Longitude
		}
	}

	if addr.Latitude == 0 && addr.Longitude == 0 {
		addr.Latitude = constant.DefaultLatitude
		addr.Longitude = constant.DefaultLongitude
	}

	return u.pharmacyDrugRepository.SelectNearestPharmaciesByDrugID(ctx, drugID, addr, constant.SearchRadiusMetre, clc)
}

func (u *pharmacyDrugUsecaseImpl) GetAllDrugsOfTheDay(ctx context.Context) ([]entity.PharmacyDrug, error) {
	currentMaxData := uint(constant.MaxBoughtData)

	pdMostBoughtOfTheDay, err := u.pharmacyDrugRepository.SelectMostBought(ctx, currentMaxData, true)
	if err != nil {
		return nil, err
	}

	pdLength := len(pdMostBoughtOfTheDay)
	if pdLength >= int(currentMaxData) {
		return pdMostBoughtOfTheDay, nil
	}

	currentMaxData -= uint(pdLength)
	pdMostBought, err := u.pharmacyDrugRepository.SelectMostBought(ctx, currentMaxData, false)
	if err != nil {
		return nil, err
	}

	pdLength = len(pdMostBought)
	if pdLength >= int(currentMaxData) {
		return append(pdMostBoughtOfTheDay, pdMostBought...), nil
	}

	currentMaxData -= uint(pdLength)
	pd, err := u.pharmacyDrugRepository.SelectAllWithLimit(ctx, currentMaxData)
	if err != nil {
		return nil, err
	}

	return append(pdMostBoughtOfTheDay, append(pdMostBought, pd...)...), nil
}
