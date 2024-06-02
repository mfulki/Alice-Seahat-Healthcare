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

type AddressRepository interface {
	SelectAllProvince(ctx context.Context) ([]entity.Province, error)
	SelectAllCitiesWithProvinceQuery(ctx context.Context, provinceID uint) ([]entity.City, error)
	SelectAllSubdistrictsWithCityQuery(ctx context.Context, cityID uint) ([]entity.Subdistrict, error)
	GetSubdistrictByID(ctx context.Context, subdistrictID uint) (*entity.Subdistrict, error)
	GetMainAddress(ctx context.Context, userID uint) (*entity.Address, error)
	GetAll(ctx context.Context, userID uint, isActive *bool) ([]entity.Address, error)
	InsertOne(ctx context.Context, addr entity.Address) (*entity.Address, error)
	GetByID(ctx context.Context, id uint, userID uint) (*entity.Address, error)
	UpdateByID(ctx context.Context, addr entity.Address) error
	DeleteByID(ctx context.Context, id uint, userID uint) error
	GetCityBySubdistrictID(ctx context.Context, subdistrictID uint) (*entity.City, error)
}

type addressRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewAddressRepository(db transaction.DBTransaction) *addressRepositoryImpl {
	return &addressRepositoryImpl{
		db: db,
	}
}

func (r *addressRepositoryImpl) SelectAllProvince(ctx context.Context) ([]entity.Province, error) {
	q := `
		SELECT 
			province_id, province_name, created_at 
		FROM 
			provinces
		ORDER BY
			province_name ASC
	`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	results := make([]entity.Province, 0)
	for rows.Next() {
		var scan entity.Province
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

func (r *addressRepositoryImpl) SelectAllCitiesWithProvinceQuery(ctx context.Context, provinceID uint) ([]entity.City, error) {
	args := make([]any, 0)
	builder := new(strings.Builder)

	builder.WriteString(`
		SELECT 
			city_id, province_id, city_name, city_type, created_at 
		FROM 
			cities
	`)

	if provinceID != 0 {
		args = append(args, provinceID)
		builder.WriteString(`
			WHERE
				province_id = $1
		`)
	}

	builder.WriteString(`
		ORDER BY
			city_name ASC
	`)

	rows, err := r.db.QueryContext(ctx, builder.String(), args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	results := make([]entity.City, 0)
	for rows.Next() {
		var scan entity.City
		if err := rows.Scan(
			&scan.ID,
			&scan.ProvinceID,
			&scan.Name,
			&scan.Type,
			&scan.CreatedAt,
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

func (r *addressRepositoryImpl) SelectAllSubdistrictsWithCityQuery(ctx context.Context, cityID uint) ([]entity.Subdistrict, error) {
	args := make([]any, 0)
	builder := new(strings.Builder)

	builder.WriteString(`
		SELECT 
			subdistrict_id, 
			city_id, 
			subdistrict_name, 
			created_at 
		FROM 
			subdistricts
	`)

	if cityID != 0 {
		args = append(args, cityID)
		builder.WriteString(`
			WHERE
				city_id = $1
		`)
	}

	builder.WriteString(`
		ORDER BY
			subdistrict_name ASC
	`)

	rows, err := r.db.QueryContext(ctx, builder.String(), args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	results := make([]entity.Subdistrict, 0)
	for rows.Next() {
		var scan entity.Subdistrict
		if err := rows.Scan(
			&scan.ID,
			&scan.CityID,
			&scan.Name,
			&scan.CreatedAt,
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

func (r *addressRepositoryImpl) GetSubdistrictByID(ctx context.Context, subdistrictID uint) (*entity.Subdistrict, error) {
	q := `
		SELECT 
			subdistrict_id, 
			city_id, 
			subdistrict_name, 
			created_at 
		FROM 
			subdistricts
		WHERE
			subdistrict_id = $1
		LIMIT 1
	`

	var scan entity.Subdistrict
	if err := r.db.QueryRowContext(ctx, q, subdistrictID).Scan(
		&scan.ID,
		&scan.CityID,
		&scan.Name,
		&scan.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrResourceNotFound
		}

		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}

func (r *addressRepositoryImpl) GetCityBySubdistrictID(ctx context.Context, subdistrictID uint) (*entity.City, error) {
	q := `
		SELECT 
			city_id, province_id, city_name, city_type, created_at 
		FROM 
			cities
		WHERE
			city_id = (SELECT city_id FROM subdistricts WHERE subdistrict_id = $1 LIMIT 1)
		LIMIT 1
	`

	var scan entity.City
	if err := r.db.QueryRowContext(ctx, q, subdistrictID).Scan(
		&scan.ID,
		&scan.ProvinceID,
		&scan.Name,
		&scan.Type,
		&scan.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrResourceNotFound
		}

		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}

func (r *addressRepositoryImpl) GetAll(ctx context.Context, userID uint, isActive *bool) ([]entity.Address, error) {
	args := []any{}
	builder := new(strings.Builder)

	args = append(args, userID)

	builder.WriteString(`
		SELECT 
			user_address_id,
			user_id, 
			subdistrict_id, 
			address, 
			ST_AsEWKT(user_addresses_location),
			is_main, 
			is_active,
			created_at
		FROM
			user_addresses
		WHERE
			user_id = $1
		AND
			deleted_at IS NULL
	`)

	if isActive != nil {
		args = append(args, *isActive)
		builder.WriteString(`
			AND
				is_active = $2
		`)
	}

	builder.WriteString(`
		ORDER BY
			created_at ASC
	`)

	rows, err := r.db.QueryContext(ctx, builder.String(), args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	results := make([]entity.Address, 0)
	for rows.Next() {
		var scan entity.Address
		if err := rows.Scan(
			&scan.ID,
			&scan.UserID,
			&scan.SubdistrictID,
			&scan.Address,
			&scan.Location,
			&scan.IsMain,
			&scan.IsActive,
			&scan.CreatedAt,
		); err != nil {
			logrus.Error(err)
			return nil, err
		}

		scan.Longitude, scan.Latitude, err = utils.Geo2LongLat(scan.Location)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		results = append(results, scan)
	}

	return results, nil
}

func (r *addressRepositoryImpl) GetMainAddress(ctx context.Context, userID uint) (*entity.Address, error) {
	q := `
		SELECT
			user_address_id,
			user_id, 
			subdistrict_id, 
			address, 
			ST_AsEWKT(user_addresses_location), 
			is_main, 
			is_active,
			created_at
		FROM
			user_addresses
		WHERE
			user_id = $1
		AND
			deleted_at IS NULL
		AND
			is_main = true
		LIMIT 1
	`

	var scan entity.Address
	err := r.db.QueryRowContext(ctx, q, userID).Scan(
		&scan.ID,
		&scan.UserID,
		&scan.SubdistrictID,
		&scan.Address,
		&scan.Location,
		&scan.IsMain,
		&scan.IsActive,
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

func (r *addressRepositoryImpl) InsertOne(ctx context.Context, addr entity.Address) (*entity.Address, error) {
	q := `
		INSERT INTO
			user_addresses (user_id, subdistrict_id, address, user_addresses_location, is_main, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			user_address_id,
			user_id, 
			subdistrict_id, 
			address, 
			ST_AsEWKT(user_addresses_location), 
			is_main, 
			is_active,
			created_at
	`

	var scan entity.Address
	err := r.db.QueryRowContext(ctx, q,
		addr.UserID,
		addr.SubdistrictID,
		addr.Address,
		fmt.Sprintf("POINT(%g %g)", addr.Longitude, addr.Latitude),
		addr.IsMain,
		addr.IsActive,
	).Scan(
		&scan.ID,
		&scan.UserID,
		&scan.SubdistrictID,
		&scan.Address,
		&scan.Location,
		&scan.IsMain,
		&scan.IsActive,
		&scan.CreatedAt,
	)

	scan.Latitude = addr.Latitude
	scan.Longitude = addr.Longitude

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}

func (r *addressRepositoryImpl) GetByID(ctx context.Context, id uint, userID uint) (*entity.Address, error) {
	q := `
		SELECT
			ua.user_address_id,
			ua.user_id, 
			ua.subdistrict_id, 
			ua.address, 
			ST_AsEWKT(ua.user_addresses_location),
			ua.user_addresses_location,
			ua.is_main, 
			ua.is_active,
			ua.created_at,
			s.city_id 
		FROM
			user_addresses ua
		JOIN
			subdistricts s on s.subdistrict_id =ua.subdistrict_id 
		WHERE
			ua.user_address_id = $1
		AND
			ua.user_id = $2
		AND
			ua.deleted_at IS NULL
	`

	var scan entity.Address
	err := r.db.QueryRowContext(ctx, q, id, userID).Scan(
		&scan.ID,
		&scan.UserID,
		&scan.SubdistrictID,
		&scan.Address,
		&scan.Location,
		&scan.RawLocation,
		&scan.IsMain,
		&scan.IsActive,
		&scan.CreatedAt,
		&scan.CityID,
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

func (r *addressRepositoryImpl) UpdateByID(ctx context.Context, addr entity.Address) error {
	q := `
		UPDATE	user_addresses
		SET 
			subdistrict_id = $1, 
			address = $2, 
			user_addresses_location = $3, 
			is_main = $4, 
			is_active = $5,
			updated_at = current_timestamp
		WHERE
			user_address_id = $6
		AND
			user_id = $7
		AND
			deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, q,
		addr.SubdistrictID,
		addr.Address,
		fmt.Sprintf("POINT(%g %g)", addr.Longitude, addr.Latitude),
		addr.IsMain,
		addr.IsActive,
		addr.ID,
		addr.UserID,
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

func (r *addressRepositoryImpl) DeleteByID(ctx context.Context, id uint, userID uint) error {
	q := `
		UPDATE	user_addresses
		SET 
			deleted_at = current_timestamp,
			updated_at = current_timestamp
		WHERE
			user_address_id = $1
		AND
			user_id = $2
		AND
			deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, q, id, userID)
	if err != nil {
		logrus.Error(err)
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return apperror.ErrResourceNotFound
	}

	return nil
}
