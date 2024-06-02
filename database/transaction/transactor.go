package transaction

import (
	"context"
	"database/sql"

	"Alice-Seahat-Healthcare/seahat-be/constant"

	"github.com/sirupsen/logrus"
)

type Transactor interface {
	WithTransaction(ctx context.Context, tFunc func(context.Context) (any, error)) (any, error)
}

type transactor struct {
	db *sql.DB
}

func NewTransactor(db *sql.DB) *transactor {
	return &transactor{db: db}
}

func (t *transactor) WithTransaction(ctx context.Context, tFunc func(context.Context) (any, error)) (any, error) {
	tx, err := t.db.Begin()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	txCtx := context.WithValue(ctx, constant.TxContext, tx)
	data, err := tFunc(txCtx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logrus.Error(err)
			return nil, err
		}

		return nil, err
	}

	if err := tx.Commit(); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return data, nil
}
