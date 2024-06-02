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

type StockRequestUsecase interface {
	ManualStockRequest(ctx context.Context, stockRequest entity.StockRequest) ([]*entity.StockRequest, error)
	UpdateStockMutationApprove(ctx context.Context, stockRequest entity.StockRequest) (*entity.StockRequest, error)
	CancelStockMutation(ctx context.Context, stockRequest entity.StockRequest) (*entity.StockRequest, error)
	GetAllStockRequest(ctx context.Context, clc *entity.Collection) ([]*entity.StockRequest, error)
	GetDrugsWithSenderAndReceiverPharmacy(ctx context.Context, senderPharmacyID uint, receiverPharmacyID uint, clc *entity.Collection) ([]entity.DrugWithPharmacyDrug, error)
}

type stockRequestUsecaseImpl struct {
	pharmacyRepository         repository.PharmacyRepository
	stockRequestRepository     repository.StockRequestRepository
	stockRequestDrugRepository repository.StockRequestDrugRepository
	transactor                 transaction.Transactor
	pharmacyDrugrepository     repository.PharmacyDrugRepository
	stockJournalRepository     repository.StockJournalRepository
}

func NewStockRequestUsecase(
	pharmacyRepository repository.PharmacyRepository,
	stockRequestRepository repository.StockRequestRepository,
	stockRequestDrugRepository repository.StockRequestDrugRepository,
	transactor transaction.Transactor,
	pharmacyDrugrepository repository.PharmacyDrugRepository,
	stockJournalRepository repository.StockJournalRepository,
) *stockRequestUsecaseImpl {
	return &stockRequestUsecaseImpl{
		pharmacyRepository:         pharmacyRepository,
		stockRequestRepository:     stockRequestRepository,
		stockRequestDrugRepository: stockRequestDrugRepository,
		transactor:                 transactor,
		pharmacyDrugrepository:     pharmacyDrugrepository,
		stockJournalRepository:     stockJournalRepository,
	}
}

func (u *stockRequestUsecaseImpl) ManualStockRequest(ctx context.Context, stockRequest entity.StockRequest) ([]*entity.StockRequest, error) {
	if stockRequest.SenderPharmacy.ID == stockRequest.ReceiverPharmacy.ID {
		return nil, apperror.CantRequestToSamePharmacy
	}
	managerCtx, ok := utils.CtxGetManager(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}
	managerId := managerCtx.ID
	pharmacyIds := []uint{stockRequest.SenderPharmacy.ID, stockRequest.ReceiverPharmacy.ID}
	lenPharmacies, err := u.pharmacyRepository.CheckPharmacyByPharmacyId(ctx, pharmacyIds, managerId)
	if err != nil {
		return nil, err
	}
	if lenPharmacies != uint(len(pharmacyIds)) {
		return nil, apperror.CantRequestToPharmaciesNotPartner
	}

	stockRequestsTx, err := u.transactor.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		stockRequestMap := make(map[uint][]*entity.StockRequestDrug)
		stockRequestMap[stockRequest.SenderPharmacy.ID] = stockRequest.StockRequestDrug
		stockRequests, err := u.stockRequestRepository.InsertStockRequestBulk(ctx, stockRequest.ReceiverPharmacy.ID, stockRequestMap, constant.WaitingApproval)
		if err != nil {
			return nil, err
		}

		err = u.stockRequestDrugRepository.InsertStockRequestDrugBulk(ctx, stockRequests)
		if err != nil {
			return nil, err
		}

		return stockRequests, nil
	})
	if err != nil {
		return nil, err
	}
	stockRequests := stockRequestsTx.([]*entity.StockRequest)
	return stockRequests, nil

}

func (u *stockRequestUsecaseImpl) UpdateStockMutationApprove(ctx context.Context, stockRequest entity.StockRequest) (*entity.StockRequest, error) {
	managerCtx, ok := utils.CtxGetManager(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}
	managerId := managerCtx.ID
	stockReqTx, err := u.transactor.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		senderSDrugs, receiverSDrugs, stockReqs, err := u.stockRequestRepository.GetStockRequestById(ctx, stockRequest.Id)

		if err != nil {
			return nil, err
		}
		var stockJournals []entity.StockJurnal
		stockReqDrugs := senderSDrugs[stockReqs.SenderPharmacy.ID]

		for _, reqDrug := range stockReqDrugs {
			stockJournal := entity.StockJurnal{DrugId: reqDrug.DrugId, PharmacyId: stockReqs.SenderPharmacy.ID, Quantity: int(reqDrug.Quantity) * -1, Description: constant.SendStockMutation}
			stockJournals = append(stockJournals, stockJournal)
			stockJournal = entity.StockJurnal{DrugId: reqDrug.DrugId, PharmacyId: stockReqs.ReceiverPharmacy.ID, Quantity: int(reqDrug.Quantity), Description: constant.ReceiveStockMutation}
			stockJournals = append(stockJournals, stockJournal)
		}

		err = u.updateStock(ctx, senderSDrugs, receiverSDrugs, stockJournals)
		if err != nil {
			return nil, err
		}
		stockReq, err := u.stockRequestRepository.UpdateStockMutationStatusBySender(ctx, stockRequest, constant.Approved, managerId)
		if err != nil {
			return nil, err
		}
		return stockReq, nil

	})
	if err != nil {
		return nil, err
	}
	stockReq := stockReqTx.(*entity.StockRequest)
	return stockReq, nil

}

func (u *stockRequestUsecaseImpl) updateStock(ctx context.Context, stockRequestDrugSender map[uint][]*entity.StockRequestDrug, stockRequestDrugReceiver map[uint][]*entity.StockRequestDrug, stockJournals []entity.StockJurnal) error {
	err := u.pharmacyDrugrepository.UpdateSubstractionBulkStock(ctx, stockRequestDrugSender)
	if err != nil {
		return err
	}
	err = u.pharmacyDrugrepository.UpdateAdditionBulkStock(ctx, stockRequestDrugReceiver)
	if err != nil {
		return err
	}
	err = u.stockJournalRepository.InsertStockJournal(ctx, stockJournals)
	if err != nil {
		return err
	}
	return nil
}

func (u *stockRequestUsecaseImpl) CancelStockMutation(ctx context.Context, stockRequest entity.StockRequest) (*entity.StockRequest, error) {
	managerCtx, ok := utils.CtxGetManager(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}
	managerId := managerCtx.ID
	stockReq, err := u.stockRequestRepository.UpdateStockMutationStatusBySender(ctx, stockRequest, constant.Cancelled, managerId)
	if err != nil {
		return nil, err
	}
	return stockReq, nil
}

func (u *stockRequestUsecaseImpl) GetAllStockRequest(ctx context.Context, clc *entity.Collection) ([]*entity.StockRequest, error) {
	managerCtx, ok := utils.CtxGetManager(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}

	managerId := managerCtx.ID
	stockRequests, err := u.stockRequestRepository.GetAllStockRequestByPharmacyManagerId(ctx, managerId)
	if err != nil {
		return nil, err
	}

	return utils.HardPagination(stockRequests, clc), nil
}

func (u *stockRequestUsecaseImpl) GetDrugsWithSenderAndReceiverPharmacy(ctx context.Context, pharmacySenderID uint, pharmacyReceiverID uint, clc *entity.Collection) ([]entity.DrugWithPharmacyDrug, error) {
	if pharmacySenderID == pharmacyReceiverID {
		return nil, apperror.SenderAndReceiverCantSame
	}

	managerCtx, ok := utils.CtxGetManager(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}

	_, err := u.pharmacyRepository.SelectByIDAndManagerID(ctx, pharmacySenderID, managerCtx.ID)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.PharmacyNotExist
		}

		return nil, err
	}

	_, err = u.pharmacyRepository.SelectByIDAndManagerID(ctx, pharmacyReceiverID, managerCtx.ID)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.PharmacyNotExist
		}

		return nil, err
	}

	return u.stockRequestRepository.SelectDrugWithSenderAndReceiverPharmacy(ctx, pharmacySenderID, pharmacyReceiverID, clc)
}
