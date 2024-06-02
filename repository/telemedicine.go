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
	telemedicineColumnAlias  = map[string]string{}
	telemedicineSearchColumn = []string{}
)

type TelemedicineRepository interface {
	SelectOneByID(ctx context.Context, id uint, userID *uint, doctorID *uint) (*entity.Telemedicine, error)
	SelectOneOngoingByUserAndDoctorID(ctx context.Context, userID uint, doctorID uint) (*entity.Telemedicine, error)
	InsertOne(ctx context.Context, newTelemedicine entity.Telemedicine) (*entity.Telemedicine, error)
	GetAllTelemedicine(ctx context.Context, userID *uint, doctorID *uint, clc *entity.Collection) ([]entity.Telemedicine, error)
	UpdateOne(ctx context.Context, updateTelemedicine entity.Telemedicine) error
	SelectOneByIdJoinDoctorUser(ctx context.Context, telemedicine entity.Telemedicine) (*entity.Telemedicine, error)
}

type telemedicineRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewTelemedicineRepository(db transaction.DBTransaction) *telemedicineRepositoryImpl {
	return &telemedicineRepositoryImpl{
		db: db,
	}
}

func (r *telemedicineRepositoryImpl) SelectOneByID(ctx context.Context, id uint, userID *uint, doctorID *uint) (*entity.Telemedicine, error) {
	q := `
		SELECT 
			telemedicine_id, user_id , doctor_id , end_at,diagnose , price, start_rest_at , rest_duration , medical_certificate_url, created_at,prescription_certificate_url	 
		FROM 
			telemedicines
		WHERE
			telemedicine_id = $1
	`

	if userID != nil {
		userIDClause := fmt.Sprintf("AND user_id = %d", *userID)
		q += userIDClause

	}

	if doctorID != nil {
		doctorIDClause := fmt.Sprintf("AND doctor_id = %d", *doctorID)
		q += doctorIDClause
	}

	var scan entity.Telemedicine
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&scan.ID,
		&scan.User.ID,
		&scan.Doctor.ID,
		&scan.EndAt,
		&scan.Diagnose,
		&scan.Price,
		&scan.StartRestAt,
		&scan.RestDuration,
		&scan.MedicalCertificateURL,
		&scan.CreatedAt,
		&scan.PrescriptionUrl,
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

func (r *telemedicineRepositoryImpl) SelectOneOngoingByUserAndDoctorID(ctx context.Context, userID uint, doctorID uint) (*entity.Telemedicine, error) {
	q := `
		SELECT 
			telemedicine_id, 
			user_id , 
			doctor_id, 
			end_at,
			diagnose, 
			price, 
			start_rest_at, 
			rest_duration, 
			medical_certificate_url, 
			created_at,
			prescription_certificate_url	 
		FROM 
			telemedicines
		WHERE
			user_id = $1
		AND
			doctor_id = $2
		AND
			end_at IS NULL
		AND
			deleted_at IS NULL
	`

	var scan entity.Telemedicine
	err := r.db.QueryRowContext(ctx, q, userID, doctorID).Scan(
		&scan.ID,
		&scan.User.ID,
		&scan.Doctor.ID,
		&scan.EndAt,
		&scan.Diagnose,
		&scan.Price,
		&scan.StartRestAt,
		&scan.RestDuration,
		&scan.MedicalCertificateURL,
		&scan.CreatedAt,
		&scan.PrescriptionUrl,
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

func (r *telemedicineRepositoryImpl) GetAllTelemedicine(ctx context.Context, userID *uint, doctorID *uint, clc *entity.Collection) ([]entity.Telemedicine, error) {
	selectColumns := `
		telemedicine_id,
		user_id,
		doctor_id,
		end_at,
		diagnose,
		price,
		start_rest_at,
		rest_duration,
		medical_certificate_url,
		created_at,
		prescription_certificate_url
	`
	advanceQuery := `
			telemedicines
		WHERE
		%s
		%s
		%s
	`

	extendFilter := ""
	if userID != nil {
		extendFilter = "user_id = $1 AND "
		clc.Args = append(clc.Args, userID)
	} else {
		extendFilter = "doctor_id = $1 AND "
		clc.Args = append(clc.Args, doctorID)
	}

	search := utils.BuildSearchQuery(telemedicineSearchColumn, clc)
	orderBy := utils.BuildSortQuery(telemedicineColumnAlias, clc.Sort, "end_at desc, telemedicine_id desc")
	filter := utils.BuildFilterQuery(telemedicineColumnAlias, clc, "deleted_at is null")

	query := utils.BuildQuery(r.db, utils.PaginateQuery{
		SelectColumns: selectColumns,
		AdvanceQuery:  fmt.Sprintf(advanceQuery, extendFilter, filter, search),
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

	telemedicines := make([]entity.Telemedicine, 0)
	for rows.Next() {
		telemedicine := entity.Telemedicine{}
		err := rows.Scan(
			&telemedicine.ID,
			&telemedicine.User.ID,
			&telemedicine.Doctor.ID,
			&telemedicine.EndAt,
			&telemedicine.Diagnose,
			&telemedicine.Price,
			&telemedicine.StartRestAt,
			&telemedicine.RestDuration,
			&telemedicine.MedicalCertificateURL,
			&telemedicine.CreatedAt,
			&telemedicine.PrescriptionUrl,
		)

		if err != nil {
			return nil, err
		}
		telemedicines = append(telemedicines, telemedicine)
	}
	return telemedicines, nil
}
func (r *telemedicineRepositoryImpl) SelectOneByIdJoinDoctorUser(ctx context.Context, telemedicine entity.Telemedicine) (*entity.Telemedicine, error) {
	q := `Select
			t.telemedicine_id, 
			t.user_id,
			u.user_name,
			u.date_of_birth,
			u.gender,
			t.doctor_id,
			d.doctor_name,
			s.specialization_name ,
			t.price,
			t.created_at,
			t.end_at
		from 
			telemedicines t 
		join 
			users u on t.user_id=u.user_id
		join 
			doctors d on d.doctor_id=t.doctor_id
		join 
			specializations s on d.specialization_id=s.specialization_id 
		where
			t.telemedicine_id =$1
		AND
			t.doctor_id=$2;`

	err := r.db.QueryRowContext(ctx, q,
		telemedicine.ID,
		telemedicine.Doctor.ID,
	).Scan(
		&telemedicine.ID,
		&telemedicine.User.ID,
		&telemedicine.User.Name,
		&telemedicine.User.DateOfBirth,
		&telemedicine.User.Gender,
		&telemedicine.Doctor.ID,
		&telemedicine.Doctor.Name,
		&telemedicine.Doctor.Specialization.Name,
		&telemedicine.Price,
		&telemedicine.CreatedAt,
		&telemedicine.EndAt,
	)

	if err != nil {
		logrus.Error(err)
		if err == sql.ErrNoRows {
			return nil, apperror.ResourceNotFound
		}
		return nil, err
	}
	return &telemedicine, nil
}

func (r *telemedicineRepositoryImpl) InsertOne(ctx context.Context, newTelemedicine entity.Telemedicine) (*entity.Telemedicine, error) {
	q := `
		INSERT INTO telemedicines 
			(user_id, doctor_id, price)
		VALUES
			($1, $2, $3)
		RETURNING
			telemedicine_id, user_id, doctor_id, end_at, diagnose, price, start_rest_at, rest_duration, medical_certificate_url, created_at,prescription_certificate_url
	`
	var scan entity.Telemedicine
	err := r.db.QueryRowContext(ctx, q,
		newTelemedicine.User.ID,
		newTelemedicine.Doctor.ID,
		newTelemedicine.Price,
	).Scan(
		&scan.ID,
		&scan.User.ID,
		&scan.Doctor.ID,
		&scan.EndAt,
		&scan.Diagnose,
		&scan.Price,
		&scan.StartRestAt,
		&scan.RestDuration,
		&scan.MedicalCertificateURL,
		&scan.CreatedAt,
		&scan.PrescriptionUrl,
	)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}

func (r *telemedicineRepositoryImpl) UpdateOne(ctx context.Context, updateTelemedicine entity.Telemedicine) error {

	q := `
		UPDATE 	
			telemedicines t
		SET 
			end_at = $1,
			diagnose = $2,
			start_rest_at = $3,
			rest_duration = $4,
			medical_certificate_url = $5,
			updated_at = now(),
			prescription_certificate_url =$6
		from
			users u
		WHERE 
			t.telemedicine_id = $7
		AND
			t.doctor_id = $8
		and 
			t.user_id =u.user_id  
	`
	result, err := r.db.ExecContext(ctx, q, updateTelemedicine.EndAt, updateTelemedicine.Diagnose, updateTelemedicine.StartRestAt, updateTelemedicine.RestDuration, updateTelemedicine.MedicalCertificateURL, updateTelemedicine.PrescriptionUrl, updateTelemedicine.ID, updateTelemedicine.Doctor.ID)
	if err != nil {
		logrus.Error(err)
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return apperror.ErrResourceNotFound
	}

	return nil

}
