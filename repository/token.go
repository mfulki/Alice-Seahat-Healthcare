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

type TokenRepository interface {
	SelectOneByToken(ctx context.Context, token entity.Token) (*entity.Token, error)
	SelectOneByDoctorID(ctx context.Context, token entity.Token) (*entity.Token, error)
	InsertOne(ctx context.Context, newData entity.Token) (*entity.Token, error)
	DeleteByUserID(ctx context.Context, userID uint, tokenType string) error
	DeleteByDoctorID(ctx context.Context, doctorID uint, tokenType string) error
	DeleteByToken(ctx context.Context, token string) error
}

type tokenRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewTokenRepository(db transaction.DBTransaction) *tokenRepositoryImpl {
	return &tokenRepositoryImpl{
		db: db,
	}
}

func (r *tokenRepositoryImpl) SelectOneByToken(ctx context.Context, token entity.Token) (*entity.Token, error) {
	q := `
		SELECT 
			token_id, user_id, doctor_id, token_type, token, expired_at, created_at
		FROM
			tokens
		WHERE
			token = $1
		AND
			expired_at > current_timestamp
		AND
			token_type = $2
		AND
			deleted_at IS NULL
	`
	var scan entity.Token
	err := r.db.QueryRowContext(ctx, q, token.Token, token.Type).Scan(
		&scan.ID,
		&scan.UserID,
		&scan.DoctorID,
		&scan.Type,
		&scan.Token,
		&scan.ExpiredAt,
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

func (r *tokenRepositoryImpl) SelectOneByDoctorID(ctx context.Context, token entity.Token) (*entity.Token, error) {
	q := `
		SELECT 
			token_id, user_id, doctor_id, token_type, token, expired_at, created_at
		FROM
			tokens
		WHERE
			doctor_id = $1
		AND
			expired_at > current_timestamp
		AND
			token_type = $2
		AND
			deleted_at IS NULL
	`
	var scan entity.Token
	err := r.db.QueryRowContext(ctx, q, token.DoctorID, token.Type).Scan(
		&scan.ID,
		&scan.UserID,
		&scan.DoctorID,
		&scan.Type,
		&scan.Token,
		&scan.ExpiredAt,
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

func (r *tokenRepositoryImpl) InsertOne(ctx context.Context, newData entity.Token) (*entity.Token, error) {
	q := `
		INSERT INTO tokens
			(user_id, doctor_id, token_type, token, expired_at)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING
			token_id, user_id, doctor_id, token_type, token, expired_at, created_at
	`

	var scan entity.Token
	err := r.db.QueryRowContext(ctx, q,
		newData.UserID,
		newData.DoctorID,
		newData.Type,
		newData.Token,
		newData.ExpiredAt,
	).Scan(
		&scan.ID,
		&scan.UserID,
		&scan.DoctorID,
		&scan.Type,
		&scan.Token,
		&scan.ExpiredAt,
		&scan.CreatedAt,
	)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}

func (r *tokenRepositoryImpl) DeleteByUserID(ctx context.Context, userID uint, tokenType string) error {
	q := `
		UPDATE
			tokens
		SET
			updated_at = CURRENT_TIMESTAMP,
			deleted_at = CURRENT_TIMESTAMP
		WHERE
			user_id = $1
		AND
			token_type = $2
		AND
		 	deleted_at IS NULL
		`

	if _, err := r.db.ExecContext(ctx, q, userID, tokenType); err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (r *tokenRepositoryImpl) DeleteByDoctorID(ctx context.Context, doctorID uint, tokenType string) error {
	q := `
		UPDATE
			tokens
		SET
			updated_at = CURRENT_TIMESTAMP,
			deleted_at = CURRENT_TIMESTAMP
		WHERE
			doctor_id = $1
		AND
			token_type = $2
		AND
		 	deleted_at IS NULL
		`

	if _, err := r.db.ExecContext(ctx, q, doctorID, tokenType); err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (r *tokenRepositoryImpl) DeleteByToken(ctx context.Context, token string) error {
	q := `
		UPDATE
			tokens
		SET
			updated_at = CURRENT_TIMESTAMP,
			deleted_at = CURRENT_TIMESTAMP
		WHERE
			token = $1
		AND
		 	deleted_at IS NULL
		`

	if _, err := r.db.ExecContext(ctx, q, token); err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}
