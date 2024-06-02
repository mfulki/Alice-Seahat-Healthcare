package usecase

import (
	"context"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/repository"
	"Alice-Seahat-Healthcare/seahat-be/utils"
)

type StockJournalUsecase interface {
	GetAllStockJournalBYPharmacyId(ctx context.Context, stockJournal entity.StockJurnal, clc *entity.Collection) ([]*entity.StockJurnal, error)
}

type stockJournalUsecaseImpl struct {
	orderRepository            repository.OrderRepository
	orderDetailRepository      repository.OrderDetailRepository
	paymentrepository          repository.PaymentRepository
	pharmacyDrugrepository     repository.PharmacyDrugRepository
	cartItemrepository         repository.CartItemRepository
	transactor                 transaction.Transactor
	stockJournalRepository     repository.StockJournalRepository
	stockRequestRepository     repository.StockRequestRepository
	stockRequestDrugRepository repository.StockRequestDrugRepository
}

func NewStockJournalUsecase(
	orderRepository repository.OrderRepository,
	orderDetailRepository repository.OrderDetailRepository,
	transactor transaction.Transactor,
	paymentrepository repository.PaymentRepository,
	pharmacyDrugrepository repository.PharmacyDrugRepository,
	cartItemrepository repository.CartItemRepository,
	stockJournalRepository repository.StockJournalRepository,
	stockRequestRepository repository.StockRequestRepository,
	stockRequestDrugRepository repository.StockRequestDrugRepository,

) *stockJournalUsecaseImpl {
	return &stockJournalUsecaseImpl{
		orderRepository:            orderRepository,
		transactor:                 transactor,
		paymentrepository:          paymentrepository,
		cartItemrepository:         cartItemrepository,
		pharmacyDrugrepository:     pharmacyDrugrepository,
		orderDetailRepository:      orderDetailRepository,
		stockJournalRepository:     stockJournalRepository,
		stockRequestRepository:     stockRequestRepository,
		stockRequestDrugRepository: stockRequestDrugRepository,
	}
}

func (u *stockJournalUsecaseImpl) GetAllStockJournalBYPharmacyId(ctx context.Context, stockJournal entity.StockJurnal, clc *entity.Collection) ([]*entity.StockJurnal, error) {
	managerID := uint(0)
	manager, _ := utils.CtxGetManager(ctx)
	if manager != nil {
		managerID = manager.ID
	}
	stockJournals, err := u.stockJournalRepository.GetAllStockJournalByPharmacyId(ctx, stockJournal, managerID, clc)
	if err != nil {
		return nil, err
	}
	if stockJournals == nil {
		return nil, apperror.Unauthorized
	}

	return stockJournals, err
}
