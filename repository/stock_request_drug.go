package repository

import (
	"context"

	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"

	"github.com/sirupsen/logrus"
)

type StockRequestDrugRepository interface {
	InsertStockRequestDrug(ctx context.Context, stockRequestDrug entity.StockRequestDrug) (*entity.StockRequestDrug, error)
	InsertStockRequestDrugBulk(ctx context.Context, stockRequests []*entity.StockRequest) error
}

type stockRequestDrugRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewStockRequestDrugRepository(db transaction.DBTransaction) *stockRequestDrugRepositoryImpl {
	return &stockRequestDrugRepositoryImpl{
		db: db,
	}
}

func (r *stockRequestDrugRepositoryImpl) InsertStockRequestDrug(ctx context.Context, stockRequestDrug entity.StockRequestDrug) (*entity.StockRequestDrug, error) {
	q := `
		insert into stock_request_drugs 
		(drug_id,stock_request_id,quantity)
		values
		($1,$2,$3)
		returning stock_request_drug_id,stock_request_id,drug_id,quantity;
	`
	err := r.db.QueryRowContext(ctx, q, stockRequestDrug.DrugId, stockRequestDrug.StockRequestId, stockRequestDrug.Quantity).Scan(&stockRequestDrug.Id, &stockRequestDrug.DrugId, &stockRequestDrug.StockRequestId, &stockRequestDrug.Quantity)
	if err != nil {
		return nil, err
	}
	return &stockRequestDrug, nil
}

func (r *stockRequestDrugRepositoryImpl) InsertStockRequestDrugBulk(ctx context.Context, stockRequests []*entity.StockRequest) error {
	q := `
		insert into stock_request_drugs 
		(drug_id,stock_request_id,quantity)
		values
		($1,$2,$3)
		returning stock_request_drug_id,stock_request_id,drug_id,quantity;
	`
	stmp, err := r.db.PrepareContext(ctx, q)
	if err != nil {
		logrus.Error(err)
		return err
	}
	for _, stockRequest := range stockRequests {
		stockRequestId := stockRequest.Id
		for _, stockrequestDrug := range stockRequest.StockRequestDrug {
			_, err := stmp.ExecContext(ctx, stockrequestDrug.DrugId, stockRequestId, stockrequestDrug.Quantity)
			if err != nil {
				logrus.Error(err)
				return err
			}
		}
	}
	return nil
}
