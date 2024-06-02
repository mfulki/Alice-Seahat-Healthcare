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
	partnerColumnAlias = map[string]string{
		"is_active":  "p.is_active",
		"created_at": "p.created_at",
	}
	partnerSearchColumn = []string{
		"p.partner_name",
		"pm.pharmacy_manager_name",
		"pm.email",
	}
)

type PartnerRepository interface {
	GetAll(ctx context.Context, clc *entity.Collection) ([]entity.Partner, error)
	InsertOne(ctx context.Context, p entity.Partner) (*entity.Partner, error)
	GetByID(ctx context.Context, id uint) (*entity.Partner, error)
	UpdateByID(ctx context.Context, p entity.Partner) (*entity.Partner, error)
}

type partnerRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewPartnerRepository(db transaction.DBTransaction) *partnerRepositoryImpl {
	return &partnerRepositoryImpl{
		db: db,
	}
}

func (r *partnerRepositoryImpl) GetAll(ctx context.Context, clc *entity.Collection) ([]entity.Partner, error) {
	selectColumns := `
		p.partner_id, 
		p.pharmacy_manager_id, 
		p.partner_name, 
		p.logo, 
		p.is_active,
		p.created_at,
		pm.pharmacy_manager_id,
		pm.pharmacy_manager_name,
		pm.email,
		pm.created_at
	`
	advanceQuery := `
			partners p
		LEFT JOIN pharmacy_managers pm
			ON p.pharmacy_manager_id = pm.pharmacy_manager_id
		WHERE
		%s
		%s
	`

	search := utils.BuildSearchQuery(partnerSearchColumn, clc)
	orderBy := utils.BuildSortQuery(partnerColumnAlias, clc.Sort, "p.partner_id desc")
	filter := utils.BuildFilterQuery(partnerColumnAlias, clc, `p.deleted_at is null and pm.deleted_at is null`)

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

	results := make([]entity.Partner, 0)
	for rows.Next() {
		var scan entity.Partner
		if err := rows.Scan(
			&scan.ID,
			&scan.PharmacyManagerID,
			&scan.Name,
			&scan.Logo,
			&scan.IsActive,
			&scan.CreatedAt,
			&scan.PharmacyManager.ID,
			&scan.PharmacyManager.Name,
			&scan.PharmacyManager.Email,
			&scan.PharmacyManager.CreatedAt,
		); err != nil {
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

func (r *partnerRepositoryImpl) InsertOne(ctx context.Context, p entity.Partner) (*entity.Partner, error) {
	q := `
		INSERT INTO
			partners (pharmacy_manager_id, partner_name, logo, is_active)
		VALUES
			($1, $2, $3, $4)
		RETURNING
			partner_id, 
			pharmacy_manager_id, 
			partner_name, 
			logo, 
			is_active,
			created_at
	`

	scan := new(entity.Partner)
	if err := r.db.QueryRowContext(ctx, q,
		p.PharmacyManagerID,
		p.Name,
		p.Logo,
		p.IsActive,
	).Scan(
		&scan.ID,
		&scan.PharmacyManagerID,
		&scan.Name,
		&scan.Logo,
		&scan.IsActive,
		&scan.CreatedAt,
	); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return scan, nil
}

func (r *partnerRepositoryImpl) GetByID(ctx context.Context, id uint) (*entity.Partner, error) {
	q := `
		SELECT
			p.partner_id, 
			p.pharmacy_manager_id, 
			p.partner_name, 
			p.logo, 
			p.is_active,
			p.created_at,
			pm.pharmacy_manager_id,
			pm.pharmacy_manager_name,
			pm.email,
			pm.created_at
		FROM
			partners p
		LEFT JOIN pharmacy_managers pm
		ON p.pharmacy_manager_id = pm.pharmacy_manager_id
		WHERE p.deleted_at IS NULL
		AND pm.deleted_at IS NULL
		AND p.partner_id = $1
		ORDER BY p.created_at DESC
	`

	var scan entity.Partner
	err := r.db.QueryRowContext(ctx, q, id).Scan(&scan.ID,
		&scan.PharmacyManagerID,
		&scan.Name,
		&scan.Logo,
		&scan.IsActive,
		&scan.CreatedAt,
		&scan.PharmacyManager.ID,
		&scan.PharmacyManager.Name,
		&scan.PharmacyManager.Email,
		&scan.PharmacyManager.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrResourceNotFound
		}

		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}

func (r *partnerRepositoryImpl) UpdateByID(ctx context.Context, p entity.Partner) (*entity.Partner, error) {
	q := `
		UPDATE partners
		SET
			partner_name = $1, 
			logo = $2, 
			is_active = $3,
			updated_at = current_timestamp
		WHERE partner_id = $4
		RETURNING
			partner_id, 
			pharmacy_manager_id, 
			partner_name, 
			logo, 
			is_active,
			created_at
	`

	var scan entity.Partner
	err := r.db.
		QueryRowContext(ctx, q, p.Name, p.Logo, p.IsActive, p.ID).
		Scan(&scan.ID, &scan.PharmacyManagerID, &scan.Name, &scan.Logo, &scan.IsActive, &scan.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrResourceNotFound
		}

		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}
