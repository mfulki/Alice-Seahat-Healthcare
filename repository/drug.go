package repository

import (
	"context"
	"database/sql"
	"fmt"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/utils"

	"github.com/sirupsen/logrus"
)

var (
	dColumnAlias = map[string]string{
		"price": "pd.price",
	}
	dSearchColumn = []string{
		"drug_name",
		"manufacturer_name",
		"drug_name",
		"generic_name",
		"composition",
		"description",
	}
)

type DrugRepository interface {
	SelectAll(context.Context, bool, *entity.Collection) ([]*entity.Drug, error)
	InsertOne(ctx context.Context, drug entity.Drug) (*entity.Drug, error)
	UpdateOne(ctx context.Context, drug entity.Drug) (*entity.Drug, error)
	DeleteOne(ctx context.Context, drug entity.Drug) (*entity.Drug, error)
	CheckNewInsertDrug(ctx context.Context, drug entity.Drug) (*entity.Drug, error)
	SelectOneById(ctx context.Context, drug entity.Drug) (*entity.Drug, error)
}

type drugRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewDrugRepository(db transaction.DBTransaction) *drugRepositoryImpl {
	return &drugRepositoryImpl{
		db: db,
	}
}

func (r *drugRepositoryImpl) SelectAll(ctx context.Context, forDoctor bool, clc *entity.Collection) ([]*entity.Drug, error) {
	selectColumns := `
		drug_id,
		d.manufacturer_id,
		manufacturer_name,
		m.created_at,
		drug_name,
		generic_name,
		composition,
		description,
		classification,
		form,
		unit_in_pack,
		selling_unit,
		weight,
		height,
		length,
		width,
		image_url,
		is_active,
		d.created_at
	`

	extendFilter := ""
	defaultOrder := "d.drug_id DESC"
	if forDoctor {
		extendFilter = "d.is_active = true AND "
		defaultOrder = "d.drug_name ASC"
	}

	advanceQuery := `
		drugs d 
		JOIN 
			manufacturers m ON m.manufacturer_id = d.manufacturer_id 
		WHERE
		%s
		%s
		%s
	`

	search := utils.BuildSearchQuery(dSearchColumn, clc)
	orderBy := utils.BuildSortQuery(dColumnAlias, clc.Sort, defaultOrder)
	filter := utils.BuildFilterQuery(dColumnAlias, clc, "d.deleted_at IS NULL")

	query := utils.BuildQuery(r.db, utils.PaginateQuery{
		SelectColumns: selectColumns,
		AdvanceQuery:  fmt.Sprintf(advanceQuery, extendFilter, filter, search),
		OrderQuery:    orderBy,
	}, clc)

	rows, err := r.db.QueryContext(ctx, query, clc.Args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	var scan []*entity.Drug
	for rows.Next() {
		s := &entity.Drug{}

		err := rows.Scan(
			&s.ID,
			&s.Manufacturer.ID,
			&s.Manufacturer.Name,
			&s.Manufacturer.CreatedAt,
			&s.Name,
			&s.GenericName,
			&s.Composition,
			&s.Description,
			&s.Classification,
			&s.Form,
			&s.UnitInPack,
			&s.SellingUnit,
			&s.Weight,
			&s.Height,
			&s.Length,
			&s.Width,
			&s.ImageURL,
			&s.IsActive,
			&s.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		scan = append(scan, s)
	}

	if err = rows.Err(); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return scan, nil
}

func (r *drugRepositoryImpl) InsertOne(ctx context.Context, drug entity.Drug) (*entity.Drug, error) {
	q := `
		INSERT INTO 
			public.drugs (manufacturer_id, drug_name, generic_name, composition, description, classification, form, unit_in_pack, selling_unit, weight,height, length, width, image_url, is_active, updated_at) 
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW()) 
		RETURNING 
			drug_id
	`

	err := r.db.QueryRowContext(ctx, q,
		drug.Manufacturer.ID,
		drug.Name,
		drug.GenericName,
		drug.Composition,
		drug.Description,
		drug.Classification,
		drug.Form,
		drug.UnitInPack,
		drug.SellingUnit,
		drug.Weight,
		drug.Height,
		drug.Length,
		drug.Width,
		drug.ImageURL,
		drug.IsActive,
	).Scan(&drug.ID)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &drug, nil
}

func (r *drugRepositoryImpl) CheckNewInsertDrug(ctx context.Context, drug entity.Drug) (*entity.Drug, error) {
	q := `SELECT drug_id from public.drugs where drug_name = $1 and generic_name = $2 and manufacturer_id = $3 and composition = $4`

	err := r.db.QueryRowContext(ctx, q, drug.Name, drug.GenericName, drug.Manufacturer.ID, drug.Composition).Scan(&drug.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.ErrResourceNotFound
		}

		logrus.Error(err)
		return nil, err
	}

	return &drug, nil
}

func (r *drugRepositoryImpl) UpdateOne(ctx context.Context, drug entity.Drug) (*entity.Drug, error) {
	q := `
		UPDATE public.drugs 
		SET 
			manufacturer_id = $1,
			drug_name = $2,
			generic_name = $3,
			composition = $4,
			description = $5,
			classification = $6,
			form = $7,
			unit_in_pack = $8,
			selling_unit = $9,
			weight = $10,
			height = $11,
			length = $12,
			width = $13,
			image_url = $14,
			is_active = $15,
			updated_at = now() 
		WHERE 
			drug_id = $16
	`
	res, err := r.db.ExecContext(ctx, q,
		drug.Manufacturer.ID,
		drug.Name,
		drug.GenericName,
		drug.Composition,
		drug.Description,
		drug.Classification,
		drug.Form,
		drug.UnitInPack,
		drug.SellingUnit,
		drug.Weight,
		drug.Height,
		drug.Length,
		drug.Width,
		drug.ImageURL,
		drug.IsActive,
		drug.ID,
	)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, apperror.ErrResourceNotFound
	}

	return &drug, nil
}

func (r *drugRepositoryImpl) DeleteOne(ctx context.Context, drug entity.Drug) (*entity.Drug, error) {
	q := `UPDATE public.drugs SET update_at = now(), deleted_at = now() WHERE drug_id = $1`

	_, err := r.db.ExecContext(ctx, q, drug.ID)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &drug, nil
}

func (r *drugRepositoryImpl) SelectOneById(ctx context.Context, drug entity.Drug) (*entity.Drug, error) {
	q := `
		SELECT
			d.manufacturer_id,
			manufacturer_name,
			m.created_at,
			drug_name,
			generic_name,
			composition,
			description,
			classification,
			form,
			unit_in_pack,
			selling_unit,
			weight,
			height,
			length,
			width,
			image_url,
			is_active,
			d.created_at
		FROM drugs d 
		JOIN manufacturers m ON m.manufacturer_id = d.manufacturer_id 
		WHERE d.deleted_at IS NULL and d.drug_id = $1
	`

	err := r.db.QueryRowContext(ctx, q, drug.ID).Scan(
		&drug.Manufacturer.ID,
		&drug.Manufacturer.Name,
		&drug.Manufacturer.CreatedAt,
		&drug.Name,
		&drug.GenericName,
		&drug.Composition,
		&drug.Description,
		&drug.Classification,
		&drug.Form,
		&drug.UnitInPack,
		&drug.SellingUnit,
		&drug.Weight,
		&drug.Height,
		&drug.Length,
		&drug.Width,
		&drug.ImageURL,
		&drug.IsActive,
		&drug.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.ErrResourceNotFound
		}

		logrus.Error(err)
		return nil, err
	}

	return &drug, nil
}
