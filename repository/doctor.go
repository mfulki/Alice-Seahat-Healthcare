package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/utils"

	"github.com/sirupsen/logrus"
)

var (
	docColumnAlias = map[string]string{
		"id":                "doctor_id",
		"status":            "status",
		"specialization_id": "specialization_id",
		"created_at":        "created_at",
	}
	docSearchColumn = []string{
		"doctor_name",
	}
)

type DoctorRepository interface {
	SelectOneByEmail(ctx context.Context, email string) (*entity.Doctor, error)
	SelectOneByID(ctx context.Context, id uint) (*entity.Doctor, error)
	SelectOneWithSpecializationByID(ctx context.Context, id uint) (*entity.Doctor, error)
	SelectAll(ctx context.Context, clc *entity.Collection) ([]entity.Doctor, error)

	InsertOne(ctx context.Context, newDoctor entity.Doctor) (*entity.Doctor, error)
	UpdateOne(ctx context.Context, doctor entity.Doctor) (*entity.Doctor, error)

	UpdatePersonalByID(ctx context.Context, doctor entity.Doctor) error
	UpdatePassword(ctx context.Context, id uint, password string) error
	UpdateStatus(ctx context.Context, doctor entity.Doctor) error
	DeleteByID(ctx context.Context, doctorID uint) error
}

type doctorRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewDoctorRepository(db transaction.DBTransaction) *doctorRepositoryImpl {
	return &doctorRepositoryImpl{
		db: db,
	}
}

func (r *doctorRepositoryImpl) SelectOneByEmail(ctx context.Context, email string) (*entity.Doctor, error) {
	q := `
		SELECT 
			doctor_id,
			specialization_id, 
			doctor_name, 
			email, 
			doctor_password, 
			date_of_birth, 
			gender, 
			doctor_certificate,
			photo_url, 
			is_oauth, 
			is_verified, 
			price,
			status,
			years_of_experience,
			created_at 
		FROM 
			doctors
		WHERE
			email = $1
		AND
			deleted_at IS NULL
	`

	var scan entity.Doctor
	err := r.db.QueryRowContext(ctx, q, email).Scan(
		&scan.ID,
		&scan.SpecializationID,
		&scan.Name,
		&scan.Email,
		&scan.Password,
		&scan.DateOfBirth,
		&scan.Gender,
		&scan.Certificate,
		&scan.PhotoURL,
		&scan.IsOAuth,
		&scan.IsVerified,
		&scan.Price,
		&scan.Status,
		&scan.YearsOfExperience,
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

func (r *doctorRepositoryImpl) SelectOneByID(ctx context.Context, id uint) (*entity.Doctor, error) {
	q := `
		SELECT 
			doctor_id,
			specialization_id, 
			doctor_name, 
			email, 
			doctor_password, 
			date_of_birth, 
			gender, 
			doctor_certificate,
			photo_url, 
			is_oauth, 
			is_verified, 
			price,
			status,
			years_of_experience,
			created_at 
		FROM 
			doctors
		WHERE
			doctor_id = $1
	`

	var scan entity.Doctor
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&scan.ID,
		&scan.SpecializationID,
		&scan.Name,
		&scan.Email,
		&scan.Password,
		&scan.DateOfBirth,
		&scan.Gender,
		&scan.Certificate,
		&scan.PhotoURL,
		&scan.IsOAuth,
		&scan.IsVerified,
		&scan.Price,
		&scan.Status,
		&scan.YearsOfExperience,
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

func (r *doctorRepositoryImpl) SelectOneWithSpecializationByID(ctx context.Context, id uint) (*entity.Doctor, error) {
	q := `
		SELECT 
			d.doctor_id,
			d.specialization_id, 
			d.doctor_name, 
			d.email, 
			d.doctor_password, 
			d.date_of_birth, 
			d.gender, 
			d.doctor_certificate,
			d.photo_url, 
			d.is_oauth, 
			d.is_verified, 
			d.price,
			d.status,
			d.years_of_experience,
			d.created_at,
			s.specialization_id,
			s.specialization_name,
			s.created_at
		FROM 
			doctors d
		LEFT JOIN specializations s ON s.specialization_id = d.specialization_id
		WHERE
			doctor_id = $1
	`

	var scan entity.Doctor
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&scan.ID,
		&scan.SpecializationID,
		&scan.Name,
		&scan.Email,
		&scan.Password,
		&scan.DateOfBirth,
		&scan.Gender,
		&scan.Certificate,
		&scan.PhotoURL,
		&scan.IsOAuth,
		&scan.IsVerified,
		&scan.Price,
		&scan.Status,
		&scan.YearsOfExperience,
		&scan.CreatedAt,
		&scan.Specialization.ID,
		&scan.Specialization.Name,
		&scan.Specialization.CreatedAt,
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

func (r *doctorRepositoryImpl) InsertOne(ctx context.Context, newDoctor entity.Doctor) (*entity.Doctor, error) {
	q := `
		INSERT INTO doctors (
			specialization_id, 
			doctor_name, 
			email, 
			doctor_password, 
			date_of_birth, 
			gender, 
			doctor_certificate,
			photo_url, 
			is_oauth, 
			is_verified, 
			price,
			status,
			years_of_experience 
		) 
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING
			doctor_id,
			specialization_id, 
			doctor_name, 
			email, 
			doctor_password, 
			date_of_birth, 
			gender, 
			doctor_certificate,
			photo_url, 
			is_oauth, 
			is_verified, 
			price,
			status,
			years_of_experience,
			created_at 
	`

	var scan entity.Doctor
	err := r.db.QueryRowContext(ctx, q,
		newDoctor.SpecializationID,
		newDoctor.Name,
		newDoctor.Email,
		newDoctor.Password,
		newDoctor.DateOfBirth,
		newDoctor.Gender,
		newDoctor.Certificate,
		newDoctor.PhotoURL,
		newDoctor.IsOAuth,
		newDoctor.IsVerified,
		newDoctor.Price,
		constant.StatusOffline,
		newDoctor.YearsOfExperience,
	).Scan(
		&scan.ID,
		&scan.SpecializationID,
		&scan.Name,
		&scan.Email,
		&scan.Password,
		&scan.DateOfBirth,
		&scan.Gender,
		&scan.Certificate,
		&scan.PhotoURL,
		&scan.IsOAuth,
		&scan.IsVerified,
		&scan.Price,
		&scan.Status,
		&scan.YearsOfExperience,
		&scan.CreatedAt,
	)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}

func (r *doctorRepositoryImpl) UpdateOne(ctx context.Context, doctor entity.Doctor) (*entity.Doctor, error) {
	q := `
		UPDATE doctors
		SET
			specialization_id = $1, 
			doctor_name = $2, 
			doctor_password = $3, 
			date_of_birth = $4, 
			gender = $5, 
			doctor_certificate = $6,
			photo_url = $6, 
			is_oauth = $7, 
			is_verified = $8, 
			price = $9,
			status = $10,
			years_of_experience = $11,
			updated_at = CURRENT_TIMESTAMP
		WHERE
			email = $12
		RETURNING
			doctor_id,
			specialization_id, 
			doctor_name, 
			email, 
			doctor_password, 
			date_of_birth, 
			gender, 
			doctor_certificate,
			photo_url, 
			is_oauth, 
			is_verified, 
			price,
			status,
			years_of_experience,
			created_at 
	`

	var scan entity.Doctor
	err := r.db.QueryRowContext(ctx, q,
		doctor.SpecializationID,
		doctor.Name,
		doctor.Password,
		doctor.DateOfBirth,
		doctor.Gender,
		doctor.Certificate,
		doctor.PhotoURL,
		doctor.IsOAuth,
		doctor.IsVerified,
		doctor.Price,
		doctor.Status,
		doctor.YearsOfExperience,
		doctor.Email,
	).Scan(
		&scan.ID,
		&scan.SpecializationID,
		&scan.Name,
		&scan.Email,
		&scan.Password,
		&scan.DateOfBirth,
		&scan.Gender,
		&scan.Certificate,
		&scan.PhotoURL,
		&scan.IsOAuth,
		&scan.IsVerified,
		&scan.Price,
		&scan.Status,
		&scan.YearsOfExperience,
		&scan.CreatedAt,
	)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}

func (r *doctorRepositoryImpl) SelectAll(ctx context.Context, clc *entity.Collection) ([]entity.Doctor, error) {
	selectColumns := `
		doctor_id,
		specialization_id, 
		doctor_name, 
		email, 
		doctor_password, 
		date_of_birth, 
		gender, 
		doctor_certificate,
		photo_url, 
		is_oauth, 
		is_verified, 
		price,
		status,
		years_of_experience,
		created_at
	`
	advanceQuery := `
			doctors 
		WHERE
		%s
		%s
	`

	doctors := []entity.Doctor{}

	search := utils.BuildSearchQuery(docSearchColumn, clc)
	orderBy := utils.BuildSortQuery(docColumnAlias, clc.Sort, "status asc")
	filter := utils.BuildFilterQuery(docColumnAlias, clc, "deleted_at is null")

	query := utils.BuildQuery(r.db, utils.PaginateQuery{
		SelectColumns: selectColumns,
		AdvanceQuery:  fmt.Sprintf(advanceQuery, filter, search),
		OrderQuery:    orderBy,
	}, clc)

	rows, err := r.db.QueryContext(ctx, query, clc.Args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrResourceNotFound
		}
		logrus.Error(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		doctor := entity.Doctor{}
		err := rows.Scan(
			&doctor.ID,
			&doctor.SpecializationID,
			&doctor.Name,
			&doctor.Email,
			&doctor.Password,
			&doctor.DateOfBirth,
			&doctor.Gender,
			&doctor.Certificate,
			&doctor.PhotoURL,
			&doctor.IsOAuth,
			&doctor.IsVerified,
			&doctor.Price,
			&doctor.Status,
			&doctor.YearsOfExperience,
			&doctor.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		doctors = append(doctors, doctor)
	}
	return doctors, nil
}

func (r *doctorRepositoryImpl) UpdatePersonalByID(ctx context.Context, doctor entity.Doctor) error {
	q := `
		UPDATE doctors
		SET
			specialization_id = $1, 
			doctor_name = $2, 
			date_of_birth = $3, 
			gender = $4, 
			doctor_certificate = $5,
			photo_url = $6, 
			price = $7,
			status = $8,
			years_of_experience = $9,
			updated_at = CURRENT_TIMESTAMP
		WHERE
			doctor_id = $10
	`

	result, err := r.db.ExecContext(ctx, q,
		doctor.SpecializationID,
		doctor.Name,
		doctor.DateOfBirth,
		doctor.Gender,
		doctor.Certificate,
		doctor.PhotoURL,
		doctor.Price,
		doctor.Status,
		doctor.YearsOfExperience,
		doctor.ID,
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

func (r *doctorRepositoryImpl) UpdatePassword(ctx context.Context, id uint, password string) error {
	q := `
		UPDATE doctors
		SET
			doctor_password = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE
			doctor_id = $2
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

func (r *doctorRepositoryImpl) UpdateStatus(ctx context.Context, doctor entity.Doctor) error {
	q := `
		UPDATE doctors
		SET
			status = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE
			doctor_id = $2
	`

	result, err := r.db.ExecContext(ctx, q, doctor.Status, doctor.ID)
	if err != nil {
		logrus.Error(err)
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return apperror.ErrResourceNotFound
	}

	return nil
}

func (r *doctorRepositoryImpl) DeleteByID(ctx context.Context, doctorID uint) error {
	q := `
		UPDATE
			doctors
		SET
			updated_at = now(),
			deleted_at = now()
		WHERE
			doctor_id = $1
		AND
			deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, q, doctorID)
	if err != nil {
		logrus.Error(err)
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return apperror.ErrResourceNotFound
	}

	return nil
}
