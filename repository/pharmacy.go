package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/utils"

	"github.com/sirupsen/logrus"
)

var (
	pharmacyColumnAlias  = map[string]string{}
	pharmacySearchColumn = []string{
		"p.pharmacy_name",
		"p.address",
		"p.pharmacist_name",
		"p.pharmacist_phone_number",
	}
)

type PharmacyRepository interface {
	InsertOne(ctx context.Context, pharmacy entity.Pharmacy) (*entity.Pharmacy, error)
	UpdateOne(ctx context.Context, pharmacy entity.Pharmacy, id uint) (*entity.Pharmacy, error)
	CheckPharmacyByPharmacyId(ctx context.Context, pharmacyIds []uint, pManagerId uint) (uint, error)
	SelectByID(ctx context.Context, pharmacyID uint) (*entity.Pharmacy, error)
	SelectByIDAndManagerID(ctx context.Context, pharmacyID uint, managerID uint) (*entity.Pharmacy, error)
	GetAllPharmacies(ctx context.Context, pharmacyManagerID uint, clc *entity.Collection) ([]entity.Pharmacy, error)
	GetPharmaciesByIDWithShipments(ctx context.Context, p entity.Pharmacy) (*entity.Pharmacy, error)
}

type pharmacyRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewPharmacyRepository(db transaction.DBTransaction) *pharmacyRepositoryImpl {
	return &pharmacyRepositoryImpl{
		db: db,
	}
}

func (r *pharmacyRepositoryImpl) InsertOne(ctx context.Context, pharmacy entity.Pharmacy) (*entity.Pharmacy, error) {
	q := `
		INSERT INTO 
			public.pharmacies (pharmacy_manager_id, subdistrict_id, pharmacy_name, pharmacy_location, address, open_time, close_time, operation_day, pharmacist_name, license_number, pharmacist_phone_number)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
		RETURNING 
			pharmacy_id, pharmacy_manager_id, subdistrict_id, pharmacy_name, ST_AsEWKT(pharmacy_location), address, open_time, close_time, operation_day, pharmacist_name, license_number, pharmacist_phone_number, created_at
	`
	var scan entity.Pharmacy

	err := r.db.QueryRowContext(ctx, q,
		pharmacy.ManagerID,
		pharmacy.SubdistrictID,
		pharmacy.Name,
		fmt.Sprintf("POINT(%g %g)", pharmacy.Longitude, pharmacy.Latitude),
		pharmacy.Address,
		pharmacy.OpenTime,
		pharmacy.CloseTime,
		pharmacy.OperationDay,
		pharmacy.PharmacistName,
		pharmacy.LicenseNumber,
		pharmacy.PharmacistPhoneNumber,
	).Scan(
		&scan.ID,
		&scan.ManagerID,
		&scan.SubdistrictID,
		&scan.Name,
		&scan.Location,
		&scan.Address,
		&scan.OpenTime,
		&scan.CloseTime,
		&scan.OperationDay,
		&scan.PharmacistName,
		&scan.LicenseNumber,
		&scan.PharmacistPhoneNumber,
		&scan.CreatedAt,
	)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	scan.Longitude, scan.Latitude, err = utils.Geo2LongLat(scan.Location)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}

func (r *pharmacyRepositoryImpl) UpdateOne(ctx context.Context, pharmacy entity.Pharmacy, id uint) (*entity.Pharmacy, error) {
	q := `
		UPDATE  
			public.pharmacies 
		SET 
			subdistrict_id = $2, 
			pharmacy_name = $3, 
			pharmacy_location = $4, 
			address = $5, 
			open_time = $6,close_time = $7, 
			operation_day = $8,
			pharmacist_name = $9,
			license_number = $10,
			pharmacist_phone_number = $11  
		WHERE
			pharmacy_id = $12
		AND
			pharmacy_manager_id = $1
		RETURNING 
			pharmacy_id,
			 pharmacy_manager_id, subdistrict_id, pharmacy_name, ST_AsEWKT(pharmacy_location), address, open_time, close_time, operation_day, pharmacist_name, license_number, pharmacist_phone_number, created_at
	`

	var scan entity.Pharmacy

	err := r.db.QueryRowContext(ctx, q,
		pharmacy.ManagerID,
		pharmacy.SubdistrictID,
		pharmacy.Name,
		fmt.Sprintf("POINT(%g %g)", pharmacy.Longitude, pharmacy.Latitude),
		pharmacy.Address,
		pharmacy.OpenTime,
		pharmacy.CloseTime,
		pharmacy.OperationDay,
		pharmacy.PharmacistName,
		pharmacy.LicenseNumber,
		pharmacy.PharmacistPhoneNumber,
		id,
	).Scan(
		&scan.ID,
		&scan.ManagerID,
		&scan.SubdistrictID,
		&scan.Name,
		&scan.Location,
		&scan.Address,
		&scan.OpenTime,
		&scan.CloseTime,
		&scan.OperationDay,
		&scan.PharmacistName,
		&scan.LicenseNumber,
		&scan.PharmacistPhoneNumber,
		&scan.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrResourceNotFound
		}

		logrus.Error(err)
		return nil, err
	}

	scan.Longitude, scan.Latitude, err = utils.Geo2LongLat(scan.Location)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}

func (r *pharmacyRepositoryImpl) CheckPharmacyByPharmacyId(ctx context.Context, pharmacyIds []uint, pManagerId uint) (uint, error) {
	q := `
		select pharmacy_id
		from pharmacies
		where pharmacy_manager_id=$1
		and deleted_at is null
		and 
	`
	var args []any
	args = append(args, pManagerId)
	for num, pharmacyId := range pharmacyIds {
		q = q + fmt.Sprintf(`pharmacy_id=$%d `, num+2)
		if num != len(pharmacyIds)-1 {
			q = q + ` or `
		}
		args = append(args, pharmacyId)
	}

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	defer rows.Close()
	var pharmacies []*entity.Pharmacy
	for rows.Next() {
		pharmacy := &entity.Pharmacy{}
		err := rows.Scan(
			&pharmacy.ID,
		)

		if err != nil {
			logrus.Error(err)
			if errors.Is(err, apperror.ErrResourceNotFound) {
				return 0, apperror.ResourceNotFound
			}
			return 0, err
		}
		pharmacies = append(pharmacies, pharmacy)

	}

	if err = rows.Err(); err != nil {
		logrus.Error(err)

		if errors.Is(err, apperror.ErrResourceNotFound) {
			return 0, apperror.ResourceNotFound
		}
		return 0, err
	}

	return uint(len(pharmacies)), nil
}

func (r *pharmacyRepositoryImpl) GetAllPharmacies(ctx context.Context, pharmacyManagerID uint, clc *entity.Collection) ([]entity.Pharmacy, error) {
	selectColumns := `
		p.pharmacy_id,
		p.pharmacy_name,
		ST_AsEWKT(p.pharmacy_location),
		p.address,
		p.pharmacy_manager_id,
		p.subdistrict_id,
		p.open_time,
		p.close_time,
		p.operation_day,
		p.pharmacist_name,
		p.license_number,
		p.pharmacist_phone_number,
		p.created_at
	`

	advanceQuery := `
			pharmacies p
		LEFT JOIN pharmacy_managers pm ON p.pharmacy_manager_id  = pm.pharmacy_manager_id 
		WHERE
		%s
		%s
	`

	search := utils.BuildSearchQuery(pharmacySearchColumn, clc)
	filter := utils.BuildFilterQuery(pharmacyColumnAlias, clc, "p.deleted_at is null")
	orderBy := utils.BuildSortQuery(pharmacyColumnAlias, clc.Sort, "p.pharmacy_id desc")

	if pharmacyManagerID != 0 {
		clc.Args = append(clc.Args, pharmacyManagerID)
		filter += fmt.Sprintf(`
			AND pm.pharmacy_manager_id = $%d
		`, len(clc.Args))
	}

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

	pharmacies := []entity.Pharmacy{}
	for rows.Next() {
		scan := entity.Pharmacy{}

		err := rows.Scan(
			&scan.ID,
			&scan.Name,
			&scan.Location,
			&scan.Address,
			&scan.ManagerID,
			&scan.SubdistrictID,
			&scan.OpenTime,
			&scan.CloseTime,
			&scan.OperationDay,
			&scan.PharmacistName,
			&scan.LicenseNumber,
			&scan.PharmacistPhoneNumber,
			&scan.CreatedAt,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		scan.Longitude, scan.Latitude, err = utils.Geo2LongLat(scan.Location)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		pharmacies = append(pharmacies, scan)
	}

	return pharmacies, nil
}

func (r *pharmacyRepositoryImpl) GetPharmaciesByIDWithShipments(ctx context.Context, p entity.Pharmacy) (*entity.Pharmacy, error) {
	args := []any{}
	q := new(strings.Builder)

	q.WriteString(`
		SELECT
			p.pharmacy_id,
			p.pharmacy_name,
			ST_AsEWKT(p.pharmacy_location),
			p.address,
			p.pharmacy_manager_id,
			p.subdistrict_id,
			p.open_time,
			p.close_time,
			p.operation_day,
			p.pharmacist_name,
			p.license_number,
			p.pharmacist_phone_number,
			p.created_at,
			sm.shipment_method_id,
			sm.shipment_method_name,
			sm.courier_name,
			sm.price,
			sm.duration,
			sm.created_at
		FROM
			pharmacies p
		LEFT JOIN pharmacy_managers pm ON p.pharmacy_manager_id  = pm.pharmacy_manager_id 
		LEFT JOIN pharmacy_shipments ps ON p.pharmacy_id = ps.pharmacy_id
		LEFT JOIN shipment_methods sm ON ps.shipment_method_id = sm.shipment_method_id
		WHERE
			p.deleted_at IS NULL
		AND
			p.pharmacy_id = $1
	`)

	args = append(args, p.ID)

	if p.ManagerID != 0 {
		q.WriteString("AND pm.pharmacy_manager_id = $2")
		args = append(args, p.ManagerID)
	}

	rows, err := r.db.QueryContext(ctx, q.String(), args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	pharmacy := entity.Pharmacy{}
	for rows.Next() {
		var sm struct {
			ID          *uint
			Name        *string
			CourierName *string
			Price       *uint
			Duration    *uint
			CreatedAt   *sql.NullTime
		}

		err := rows.Scan(
			&pharmacy.ID,
			&pharmacy.Name,
			&pharmacy.Location,
			&pharmacy.Address,
			&pharmacy.ManagerID,
			&pharmacy.SubdistrictID,
			&pharmacy.OpenTime,
			&pharmacy.CloseTime,
			&pharmacy.OperationDay,
			&pharmacy.PharmacistName,
			&pharmacy.LicenseNumber,
			&pharmacy.PharmacistPhoneNumber,
			&pharmacy.CreatedAt,
			&sm.ID,
			&sm.Name,
			&sm.CourierName,
			&sm.Price,
			&sm.Duration,
			&sm.CreatedAt,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		pharmacy.Longitude, pharmacy.Latitude, err = utils.Geo2LongLat(pharmacy.Location)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		if sm.Name != nil {
			pharmacy.ShipmentMethods = append(pharmacy.ShipmentMethods, &entity.ShipmentMethod{
				ID:          *sm.ID,
				Name:        *sm.Name,
				CourierName: *sm.CourierName,
				Price:       sm.Price,
				Duration:    *sm.Duration,
				CreatedAt:   sm.CreatedAt,
			})
		}
	}

	if pharmacy.ID == 0 {
		return nil, apperror.ErrResourceNotFound
	}

	return &pharmacy, nil
}

func (r *pharmacyRepositoryImpl) SelectByID(ctx context.Context, pharmacyID uint) (*entity.Pharmacy, error) {
	q := `
		SELECT
			p.pharmacy_id,
			p.pharmacy_manager_id,
			p.pharmacy_name,
			ST_AsEWKT(p.pharmacy_location),
			p.address,
			p.subdistrict_id,
			p.open_time,
			p.close_time,
			p.operation_day,
			p.pharmacist_name,
			p.license_number,
			p.pharmacist_phone_number,
			p.created_at
		FROM
			pharmacies p
		WHERE
			p.pharmacy_id = $1
	`

	var scan entity.Pharmacy
	err := r.db.QueryRowContext(ctx, q, pharmacyID).Scan(
		&scan.ID,
		&scan.ManagerID,
		&scan.Name,
		&scan.Location,
		&scan.Address,
		&scan.SubdistrictID,
		&scan.OpenTime,
		&scan.CloseTime,
		&scan.OperationDay,
		&scan.PharmacistName,
		&scan.LicenseNumber,
		&scan.PharmacistPhoneNumber,
		&scan.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrResourceNotFound
		}

		logrus.Error(err)
		return nil, err
	}

	scan.Longitude, scan.Latitude, err = utils.Geo2LongLat(scan.Location)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}

func (r *pharmacyRepositoryImpl) SelectByIDAndManagerID(ctx context.Context, pharmacyID uint, managerID uint) (*entity.Pharmacy, error) {
	q := `
	SELECT
		p.pharmacy_id,
		p.pharmacy_manager_id,
		p.pharmacy_name,
		ST_AsEWKT(p.pharmacy_location),
		p.address,
		p.subdistrict_id,
		p.open_time,
		p.close_time,
		p.operation_day,
		p.pharmacist_name,
		p.license_number,
		p.pharmacist_phone_number,
		p.created_at
	FROM
		pharmacies p
	WHERE
		p.pharmacy_id = $1
	AND
		p.pharmacy_manager_id = $2
`

	var scan entity.Pharmacy
	err := r.db.QueryRowContext(ctx, q, pharmacyID, managerID).Scan(
		&scan.ID,
		&scan.ManagerID,
		&scan.Name,
		&scan.Location,
		&scan.Address,
		&scan.SubdistrictID,
		&scan.OpenTime,
		&scan.CloseTime,
		&scan.OperationDay,
		&scan.PharmacistName,
		&scan.LicenseNumber,
		&scan.PharmacistPhoneNumber,
		&scan.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrResourceNotFound
		}

		logrus.Error(err)
		return nil, err
	}

	scan.Longitude, scan.Latitude, err = utils.Geo2LongLat(scan.Location)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}
