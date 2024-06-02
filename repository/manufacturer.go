package repository

import (
	"context"
	"fmt"

	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/utils"

	"github.com/sirupsen/logrus"
)

var (
	manufacturerColumnAlias = map[string]string{
		"id":         "manufacturer_id",
		"name":       "manufacturer_name",
		"created_at": "created_at",
	}
	manufacturerSearchColumn = []string{
		"manufacturer_name",
	}
)

type ManufacturerRepository interface {
	SelectAll(context.Context, *entity.Collection) ([]entity.Manufacturer, error)
}

type manufacturerRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewManufacturerRepository(db transaction.DBTransaction) *manufacturerRepositoryImpl {
	return &manufacturerRepositoryImpl{
		db: db,
	}
}

func (r *manufacturerRepositoryImpl) SelectAll(ctx context.Context, clc *entity.Collection) ([]entity.Manufacturer, error) {
	selectColums := `
		manufacturer_id,
		manufacturer_name,
		created_at
	`

	advanceQuery := `
		manufacturers
		WHERE
		%s
		%s
	`

	filter := utils.BuildFilterQuery(manufacturerColumnAlias, clc, "deleted_at IS NULL")
	orderBy := utils.BuildSortQuery(manufacturerColumnAlias, clc.Sort, "manufacturer_name ASC")
	search := utils.BuildSearchQuery(manufacturerSearchColumn, clc)

	q := utils.BuildQuery(r.db, utils.PaginateQuery{
		SelectColumns: selectColums,
		AdvanceQuery:  fmt.Sprintf(advanceQuery, filter, search),
		OrderQuery:    orderBy,
	}, clc)

	rows, err := r.db.QueryContext(ctx, q, clc.Args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	scan := make([]entity.Manufacturer, 0)
	for rows.Next() {
		s := entity.Manufacturer{}

		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.CreatedAt,
		)

		if err != nil {
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
