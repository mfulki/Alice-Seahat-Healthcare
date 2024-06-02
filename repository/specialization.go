package repository

import (
	"context"

	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"

	"github.com/sirupsen/logrus"
)

type SpecializationRepository interface {
	SelectAll(ctx context.Context) ([]entity.Specialization, error)
}

type specializationRepostioryImpl struct {
	db transaction.DBTransaction
}

func NewSpecializationRepository(db transaction.DBTransaction) *specializationRepostioryImpl {
	return &specializationRepostioryImpl{
		db: db,
	}
}

func (r *specializationRepostioryImpl) SelectAll(ctx context.Context) ([]entity.Specialization, error) {
	q := `
		SELECT 
			specialization_id,
			specialization_name,
			created_at
		FROM 
			specializations
		ORDER BY
			specialization_name ASC
	`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	results := make([]entity.Specialization, 0)
	for rows.Next() {
		var scan entity.Specialization
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
