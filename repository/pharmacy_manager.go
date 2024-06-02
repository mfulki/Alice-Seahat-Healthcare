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

type PharmacyManagerRepository interface {
	SelectAll(ctx context.Context, clc *entity.Collection) ([]entity.PharmacyManager, error)
	SelectOneByEmail(ctx context.Context, email string) (*entity.PharmacyManager, error)
	InsertOne(ctx context.Context, pm entity.PharmacyManager) (*entity.PharmacyManager, error)
	UpdateByID(ctx context.Context, pm entity.PharmacyManager) error
}

type pharmacyManagerRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewPharmacyManagerRepository(db transaction.DBTransaction) *pharmacyManagerRepositoryImpl {
	return &pharmacyManagerRepositoryImpl{
		db: db,
	}
}

func (r *pharmacyManagerRepositoryImpl) SelectAll(ctx context.Context, clc *entity.Collection) ([]entity.PharmacyManager, error) {
	selectColumns := `
		pharmacy_manager_id,
		pharmacy_manager_name, 
		email, 
		pharmacy_manager_password, 
		created_at 
	`

	advanceQuery := `
		pharmacy_managers
		WHERE
		%s
		%s
	`

	search := utils.BuildSearchQuery(nil, clc)
	orderBy := utils.BuildSortQuery(nil, clc.Sort, "pharmacy_manager_id DESC")
	filter := utils.BuildFilterQuery(nil, clc, "deleted_at IS NULL")

	query := utils.BuildQuery(r.db, utils.PaginateQuery{
		SelectColumns: selectColumns,
		AdvanceQuery:  fmt.Sprintf(advanceQuery, filter, search),
		OrderQuery:    orderBy,
	}, clc)

	results := make([]entity.PharmacyManager, 0)
	rows, err := r.db.QueryContext(ctx, query, clc.Args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var scan entity.PharmacyManager

		if err := rows.Scan(
			&scan.ID,
			&scan.Name,
			&scan.Email,
			&scan.Password,
			&scan.CreatedAt,
		); err != nil {
			logrus.Error(err)
			return nil, err
		}

		results = append(results, scan)
	}

	return results, nil
}

func (r *pharmacyManagerRepositoryImpl) SelectOneByEmail(ctx context.Context, email string) (*entity.PharmacyManager, error) {
	q := `
		SELECT 
			pharmacy_manager_id,
			pharmacy_manager_name, 
			email, 
			pharmacy_manager_password, 
			created_at 
		FROM 
			pharmacy_managers
		WHERE
			email = $1
	`

	var scan entity.PharmacyManager
	err := r.db.QueryRowContext(ctx, q, email).Scan(
		&scan.ID,
		&scan.Name,
		&scan.Email,
		&scan.Password,
		&scan.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrResourceNotFound
		}

		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}

func (r *pharmacyManagerRepositoryImpl) InsertOne(ctx context.Context, pm entity.PharmacyManager) (*entity.PharmacyManager, error) {
	q := `
		INSERT INTO pharmacy_managers
			(pharmacy_manager_name, email, pharmacy_manager_password)
		VALUES
			($1, $2, $3)
		RETURNING
			pharmacy_manager_id,
			pharmacy_manager_name, 
			email, 
			pharmacy_manager_password, 
			created_at 
	`

	scan := new(entity.PharmacyManager)
	if err := r.db.QueryRowContext(ctx, q,
		pm.Name,
		pm.Email,
		pm.Password,
	).Scan(
		&scan.ID,
		&scan.Name,
		&scan.Email,
		&scan.Password,
		&scan.CreatedAt,
	); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return scan, nil
}

func (r *pharmacyManagerRepositoryImpl) UpdateByID(ctx context.Context, pm entity.PharmacyManager) error {
	q := `
		UPDATE pharmacy_managers
		SET
			pharmacy_manager_name = $1,
			updated_at = current_timestamp
		WHERE deleted_at IS NULL
		AND pharmacy_manager_id = $2
	`

	result, err := r.db.ExecContext(ctx, q, pm.Name, pm.ID)
	if err != nil {
		logrus.Error(err)
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return apperror.ErrResourceNotFound
	}

	return nil
}
