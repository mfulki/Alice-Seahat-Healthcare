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
	catColumnAlias = map[string]string{
		"name":       "category_name",
		"created_at": "created_at",
	}
	catSearchColumn = []string{
		"category_name",
	}
)

type CategoryRepository interface {
	SelectAll(ctx context.Context, clc *entity.Collection) ([]entity.Category, error)
	InsertOne(ctx context.Context, cat entity.Category) (*entity.Category, error)
	SelectByID(ctx context.Context, categoryID uint) (*entity.Category, error)
	SelectByInsensitiveName(ctx context.Context, categoryName string) (*entity.Category, error)
	UpdateOne(ctx context.Context, cat entity.Category) error
	DeleteOneByID(ctx context.Context, categoryID uint) error
}

type categoryRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewCategoryRepository(db transaction.DBTransaction) *categoryRepositoryImpl {
	return &categoryRepositoryImpl{
		db: db,
	}
}

func (r *categoryRepositoryImpl) SelectAll(ctx context.Context, clc *entity.Collection) ([]entity.Category, error) {
	selectColumns := `
		category_id,
		category_name,
		created_at
	`

	advanceQuery := `
		categories
		WHERE
		%s
		%s
	`

	search := utils.BuildSearchQuery(catSearchColumn, clc)
	orderBy := utils.BuildSortQuery(catColumnAlias, clc.Sort, "category_id desc")
	filter := utils.BuildFilterQuery(catColumnAlias, clc, "deleted_at IS NULL")

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

	results := make([]entity.Category, 0)
	for rows.Next() {
		var scan entity.Category
		if err := rows.Scan(&scan.ID, &scan.Name, &scan.CreatedAt); err != nil {
			logrus.Error(err)
			return nil, err
		}

		results = append(results, scan)
	}

	if err := rows.Err(); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return results, nil
}

func (r *categoryRepositoryImpl) InsertOne(ctx context.Context, cat entity.Category) (*entity.Category, error) {
	q := `
		INSERT INTO
			categories (category_name)
		VALUES
			($1)
		RETURNING
			category_id,
			category_name,
			created_at
	`

	var scan entity.Category
	if err := r.db.QueryRowContext(ctx, q, cat.Name).Scan(&scan.ID, &scan.Name, &scan.CreatedAt); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}

func (r *categoryRepositoryImpl) SelectByID(ctx context.Context, categoryID uint) (*entity.Category, error) {
	q := `
		SELECT 
			category_id,
			category_name,
			created_at
		FROM 
			categories
		WHERE
			category_id = $1
		AND
			deleted_at IS NULL
	`

	var scan entity.Category
	if err := r.db.QueryRowContext(ctx, q, categoryID).Scan(&scan.ID, &scan.Name, &scan.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrResourceNotFound
		}

		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}

func (r *categoryRepositoryImpl) SelectByInsensitiveName(ctx context.Context, categoryName string) (*entity.Category, error) {
	q := `
		SELECT 
			category_id,
			category_name,
			created_at
		FROM 
			categories
		WHERE
			LOWER(category_name) = LOWER($1)
		AND
			deleted_at IS NULL
	`

	var scan entity.Category
	if err := r.db.QueryRowContext(ctx, q, categoryName).Scan(&scan.ID, &scan.Name, &scan.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrResourceNotFound
		}

		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}

func (r *categoryRepositoryImpl) UpdateOne(ctx context.Context, cat entity.Category) error {
	q := `
		UPDATE
			categories
		SET
			category_name = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE
			category_id = $2
		RETURNING
			category_id,
			category_name,
			created_at
	`

	if _, err := r.db.ExecContext(ctx, q, cat.Name, cat.ID); err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (r *categoryRepositoryImpl) DeleteOneByID(ctx context.Context, categoryID uint) error {
	q := `
		UPDATE
			categories
		SET
			updated_at = CURRENT_TIMESTAMP,
			deleted_at = CURRENT_TIMESTAMP
		WHERE
			category_id = $1
		RETURNING
			category_id,
			category_name,
			created_at
	`

	if _, err := r.db.ExecContext(ctx, q, categoryID); err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}
