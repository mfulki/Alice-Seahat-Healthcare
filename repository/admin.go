package repository

import (
	"context"
	"database/sql"
	"errors"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"

	"github.com/sirupsen/logrus"
)

type AdminRepository interface {
	SelectOneByEmail(ctx context.Context, email string) (*entity.Admin, error)
}

type adminRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewAdminRepository(db transaction.DBTransaction) *adminRepositoryImpl {
	return &adminRepositoryImpl{
		db: db,
	}
}

func (r *adminRepositoryImpl) SelectOneByEmail(ctx context.Context, email string) (*entity.Admin, error) {
	q := `
		SELECT 
			admin_id,
			admin_name, 
			email, 
			admin_password, 
			created_at 
		FROM 
			admins
		WHERE
			email = $1
	`

	var scan entity.Admin
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
