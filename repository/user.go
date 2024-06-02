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

type UserRepository interface {
	SelectAll(ctx context.Context, clc *entity.Collection) ([]entity.User, error)
	SelectOneByID(ctx context.Context, id uint) (*entity.User, error)
	SelectOneByEmail(ctx context.Context, email string) (*entity.User, error)
	InsertOne(ctx context.Context, newUser entity.User) (*entity.User, error)
	UpdateOne(ctx context.Context, user entity.User) (*entity.User, error)
	UpdatePersonalByID(ctx context.Context, user entity.User) error
	UpdatePassword(ctx context.Context, id uint, password string) error
}

type userRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewUserRepository(db transaction.DBTransaction) *userRepositoryImpl {
	return &userRepositoryImpl{
		db: db,
	}
}

func (r *userRepositoryImpl) SelectAll(ctx context.Context, clc *entity.Collection) ([]entity.User, error) {
	selectColumns := `
		user_id, user_name, email, user_password, date_of_birth, gender, photo_url, is_oauth, is_verified, created_at 
	`

	advanceQuery := `
		users
		WHERE
		%s
		%s
	`

	search := utils.BuildSearchQuery(nil, clc)
	orderBy := utils.BuildSortQuery(nil, clc.Sort, "user_id DESC")
	filter := utils.BuildFilterQuery(nil, clc, "deleted_at IS NULL")

	query := utils.BuildQuery(r.db, utils.PaginateQuery{
		SelectColumns: selectColumns,
		AdvanceQuery:  fmt.Sprintf(advanceQuery, filter, search),
		OrderQuery:    orderBy,
	}, clc)

	results := make([]entity.User, 0)
	rows, err := r.db.QueryContext(ctx, query, clc.Args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var scan entity.User

		if err := rows.Scan(
			&scan.ID,
			&scan.Name,
			&scan.Email,
			&scan.Password,
			&scan.DateOfBirth,
			&scan.Gender,
			&scan.PhotoURL,
			&scan.IsOAuth,
			&scan.IsVerified,
			&scan.CreatedAt,
		); err != nil {
			logrus.Error(err)
			return nil, err
		}

		results = append(results, scan)
	}

	return results, nil
}

func (r *userRepositoryImpl) SelectOneByID(ctx context.Context, id uint) (*entity.User, error) {
	q := `
		SELECT 
			user_id, user_name, email, user_password, date_of_birth, gender, photo_url, is_oauth, is_verified, created_at 
		FROM 
			users
		WHERE
			user_id = $1
	`

	var scan entity.User
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&scan.ID,
		&scan.Name,
		&scan.Email,
		&scan.Password,
		&scan.DateOfBirth,
		&scan.Gender,
		&scan.PhotoURL,
		&scan.IsOAuth,
		&scan.IsVerified,
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

func (r *userRepositoryImpl) SelectOneByEmail(ctx context.Context, email string) (*entity.User, error) {
	q := `
		SELECT 
			user_id, user_name, email, user_password, date_of_birth, gender, photo_url, is_oauth, is_verified, created_at 
		FROM 
			users
		WHERE
			email = $1
	`

	var scan entity.User
	err := r.db.QueryRowContext(ctx, q, email).Scan(
		&scan.ID,
		&scan.Name,
		&scan.Email,
		&scan.Password,
		&scan.DateOfBirth,
		&scan.Gender,
		&scan.PhotoURL,
		&scan.IsOAuth,
		&scan.IsVerified,
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

func (r *userRepositoryImpl) InsertOne(ctx context.Context, newUser entity.User) (*entity.User, error) {
	q := `
		INSERT INTO users 
			(user_name, email, user_password, date_of_birth, gender, photo_url, is_oauth, is_verified)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
			user_id, user_name, email, user_password, date_of_birth, gender, photo_url, is_oauth, is_verified, created_at
	`

	var scan entity.User
	err := r.db.QueryRowContext(ctx, q,
		newUser.Name,
		newUser.Email,
		newUser.Password,
		newUser.DateOfBirth,
		newUser.Gender,
		newUser.PhotoURL,
		newUser.IsOAuth,
		newUser.IsOAuth,
	).Scan(
		&scan.ID,
		&scan.Name,
		&scan.Email,
		&scan.Password,
		&scan.DateOfBirth,
		&scan.Gender,
		&scan.PhotoURL,
		&scan.IsOAuth,
		&scan.IsVerified,
		&scan.CreatedAt,
	)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}

func (r *userRepositoryImpl) UpdateOne(ctx context.Context, user entity.User) (*entity.User, error) {
	q := `
		UPDATE users
		SET
			user_name = $1,
			user_password = $2,
			date_of_birth = $3,
			gender = $4,
			photo_url = $5,
			is_oauth = $6,
			is_verified = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE
			email = $8
		RETURNING
			user_id, user_name, email, user_password, date_of_birth, gender, photo_url, is_oauth, is_verified, created_at
	`

	var scan entity.User
	err := r.db.QueryRowContext(ctx, q,
		user.Name,
		user.Password,
		user.DateOfBirth,
		user.Gender,
		user.PhotoURL,
		user.IsOAuth,
		user.IsVerified,
		user.Email,
	).Scan(
		&scan.ID,
		&scan.Name,
		&scan.Email,
		&scan.Password,
		&scan.DateOfBirth,
		&scan.Gender,
		&scan.PhotoURL,
		&scan.IsOAuth,
		&scan.IsVerified,
		&scan.CreatedAt,
	)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}

func (r *userRepositoryImpl) UpdatePersonalByID(ctx context.Context, user entity.User) error {
	q := `
		UPDATE users
		SET
			user_name = $1,
			date_of_birth = $2,
			gender = $3,
			photo_url = $4,
			updated_at = CURRENT_TIMESTAMP
		WHERE
			user_id = $5
	`

	result, err := r.db.ExecContext(ctx, q,
		user.Name,
		user.DateOfBirth,
		user.Gender,
		user.PhotoURL,
		user.ID,
	)

	if err != nil {
		logrus.Error(err)
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return apperror.ErrResourceNotFound
	}

	return nil
}

func (r *userRepositoryImpl) UpdatePassword(ctx context.Context, id uint, password string) error {
	q := `
		UPDATE users
		SET
			user_password = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE
			user_id = $2
	`

	result, err := r.db.ExecContext(ctx, q, password, id)
	if err != nil {
		logrus.Error(err)
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return apperror.ErrResourceNotFound
	}

	return nil
}
