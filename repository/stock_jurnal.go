package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/utils"

	"github.com/sirupsen/logrus"
)

var (
	journalColumnAlias  = map[string]string{}
	journalSearchColumn = []string{
		"s.description",
	}
)

type StockJournalRepository interface {
	InsertStockJournal(ctx context.Context, stockJurnals []entity.StockJurnal) error
	GetAllStockJournalByPharmacyId(ctx context.Context, stockJournal entity.StockJurnal, pMId uint, clc *entity.Collection) ([]*entity.StockJurnal, error)
}

type stockJournalRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewStockJournalRepository(db transaction.DBTransaction) *stockJournalRepositoryImpl {
	return &stockJournalRepositoryImpl{
		db: db,
	}
}

func (r *stockJournalRepositoryImpl) GetAllStockJournalByPharmacyId(ctx context.Context, stockJournal entity.StockJurnal, pMId uint, clc *entity.Collection) ([]*entity.StockJurnal, error) {
	selectColumns := `
		s.stock_journal_id,
		s.drug_id,
		d.drug_name,
		s.pharmacy_id,
		s.quantity,
		s.description,
		s.created_at,
		s.updated_at,
		s.deleted_at
	`
	advanceQuery := `
			stock_journals s 
		join pharmacies p on p.pharmacy_id=s.pharmacy_id
		join drugs d on s.drug_id=d.drug_id
		WHERE
			s.pharmacy_id=$1
		AND
			p.pharmacy_manager_id=$2
		AND
		%s
		%s
	`

	clc.Args = append(clc.Args, stockJournal.PharmacyId, pMId)

	search := utils.BuildSearchQuery(journalSearchColumn, clc)
	orderBy := utils.BuildSortQuery(journalColumnAlias, clc.Sort, "s.stock_journal_id desc")
	filter := utils.BuildFilterQuery(journalColumnAlias, clc, "s.deleted_at is null")

	query := utils.BuildQuery(r.db, utils.PaginateQuery{
		SelectColumns: selectColumns,
		AdvanceQuery:  fmt.Sprintf(advanceQuery, filter, search),
		OrderQuery:    orderBy,
	}, clc)

	rows, err := r.db.QueryContext(ctx, query, clc.Args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	var scan []*entity.StockJurnal
	for rows.Next() {
		s := &entity.StockJurnal{}
		err := rows.Scan(
			&s.Id,
			&s.DrugId,
			&s.DrugName,
			&s.PharmacyId,
			&s.Quantity,
			&s.Description,
			&s.CreatedAt,
			&s.UpdatedAt,
			&s.DeletedAt,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, apperror.ErrResourceNotFound
			}

			logrus.Error(err)
			return nil, err
		}
		scan = append(scan, s)

	}

	if err = rows.Err(); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return scan, nil
}

func (r *stockJournalRepositoryImpl) InsertStockJournal(ctx context.Context, stockJurnals []entity.StockJurnal) error {
	q := `insert into stock_journals 
			(drug_id,pharmacy_id,quantity,description)
			values
			($1,$2,$3,$4)
		`
	stmt, err := r.db.PrepareContext(ctx, q)
	if err != nil {
		logrus.Error(err)
		return err
	}
	for _, stockJurnal := range stockJurnals {
		_, err := stmt.ExecContext(ctx, stockJurnal.DrugId, stockJurnal.PharmacyId, stockJurnal.Quantity, stockJurnal.Description)
		if err != nil {
			logrus.Error(err)
			return err
		}

	}
	return nil

}
