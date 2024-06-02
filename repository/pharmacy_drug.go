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
	pdColumnAlias = map[string]string{
		"price":         "pd.price",
		"category_name": "c.category_name",
	}
	pdSearchColumn = []string{
		"d.drug_name",
	}
)

type PharmacyDrugRepository interface {
	SelectAllWithinRadius(context.Context, entity.Address, int, *entity.Collection) ([]*entity.PharmacyDrug, error)
	SelectOneByID(ctx context.Context, id uint) (*entity.PharmacyDrug, error)
	SelectAll(ctx context.Context, pharmacyID uint, managerID uint, clc *entity.Collection) ([]entity.PharmacyDrug, error)
	SelectById(ctx context.Context, id uint) (*entity.PharmacyDrug, error)
	CreateOne(ctx context.Context, pharmacyDrug entity.PharmacyDrug) (*entity.PharmacyDrug, error)
	UpdateStockComeOut(ctx context.Context, quantity uint, pharmacyDrugId uint) (*uint, error)
	UpdateReturnStock(ctx context.Context, orderDetails []*entity.OrderDetail) ([]entity.StockJurnal, error)
	SelectOneNearestByPharmacyManagerDrugId(ctx context.Context, pharmacyDrug entity.PharmacyDrug, quantity uint) (*entity.PharmacyDrug, error)
	UpdateSubstractionBulkStock(ctx context.Context, stockRequestDrugs map[uint][]*entity.StockRequestDrug) error
	UpdateOne(ctx context.Context, pharmacyDrug entity.PharmacyDrug, id uint) error
	SelectNearestPharmaciesByDrugID(ctx context.Context, drugID uint, addr entity.Address, radius int, clc *entity.Collection) ([]entity.PharmacyDrug, error)
	SelectMostBought(ctx context.Context, limit uint, isTheDay bool) ([]entity.PharmacyDrug, error)
	SelectAllWithLimit(ctx context.Context, limit uint) ([]entity.PharmacyDrug, error)
	SelectOneByPharmacyAndDrugID(ctx context.Context, pd entity.PharmacyDrug) (*entity.PharmacyDrug, error)
	UpdateAdditionBulkStock(ctx context.Context, stockRequestDrugs map[uint][]*entity.StockRequestDrug) error
	SelectPharmacyDrugsByCategoryId(ctx context.Context, categoryId uint) (int, error)
}

type pharmacyDrugRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewPharmacyDrugRepository(db transaction.DBTransaction) *pharmacyDrugRepositoryImpl {
	return &pharmacyDrugRepositoryImpl{
		db: db,
	}
}

func (r *pharmacyDrugRepositoryImpl) SelectAllWithinRadius(ctx context.Context, addr entity.Address, radius int, clc *entity.Collection) ([]*entity.PharmacyDrug, error) {
	selectColumns := `
		p.pharmacy_id,
		p.pharmacy_name,
		p.open_time,
		p.close_time,
		p.address,
		p.subdistrict_id,
		ST_Distance(
			ST_MakePoint($1, $2) :: geography,
			p.pharmacy_location
		) AS distance,
		ST_AsEWKT(p.pharmacy_location),
		pd.pharmacy_drug_id,
		pd.stock,
		pd.price,
		pd.category_id,
		pd.is_active,
		d.drug_id,
		d.manufacturer_id,
		m.manufacturer_name,
		d.drug_name,
		d.generic_name,
		d.composition,
		d.description,
		d.classification,
		d.form,
		d.unit_in_pack,
		d.selling_unit,
		d.weight,
		d.height,
		d.length,
		d.width,
		d.image_url,
		d.is_active,
		c.category_id,
		c.category_name,
		c.created_at
	`

	advanceQuery := `
		pharmacies p
			JOIN pharmacy_drugs pd ON pd.pharmacy_id = p.pharmacy_id
			JOIN drugs d ON pd.drug_id = d.drug_id
			JOIN manufacturers m ON m.manufacturer_id = d.manufacturer_id
			LEFT JOIN categories c ON c.category_id = pd.category_id
		WHERE
			d.is_active = true
		AND
			pd.is_active = true
		AND
		%s
		%s
	`

	clc.Args = append(clc.Args, addr.Longitude, addr.Latitude, radius)

	search := utils.BuildSearchQuery(pdSearchColumn, clc)
	orderBy := utils.BuildSortQuery(pdColumnAlias, clc.Sort, "distance ASC")
	filter := utils.BuildFilterQuery(pdColumnAlias, clc, `
		ST_DWithin(
			p.pharmacy_location,
			ST_MakePoint($1, $2) :: geography,
			$3
		)
	`)

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

	var scan []*entity.PharmacyDrug
	for rows.Next() {
		s := &entity.PharmacyDrug{}

		err := rows.Scan(
			&s.Pharmacy.ID,
			&s.Pharmacy.Name,
			&s.Pharmacy.OpenTime,
			&s.Pharmacy.CloseTime,
			&s.Pharmacy.Address,
			&s.Pharmacy.SubdistrictID,
			&s.Pharmacy.Distance,
			&s.Pharmacy.Location,
			&s.ID,
			&s.Stock,
			&s.Price,
			&s.Category.ID,
			&s.IsActive,
			&s.Drug.ID,
			&s.Drug.Manufacturer.ID,
			&s.Drug.Manufacturer.Name,
			&s.Drug.Name,
			&s.Drug.GenericName,
			&s.Drug.Composition,
			&s.Drug.Description,
			&s.Drug.Classification,
			&s.Drug.Form,
			&s.Drug.UnitInPack,
			&s.Drug.SellingUnit,
			&s.Drug.Weight,
			&s.Drug.Height,
			&s.Drug.Length,
			&s.Drug.Width,
			&s.Drug.ImageURL,
			&s.Drug.IsActive,
			&s.Category.ID,
			&s.Category.Name,
			&s.Category.CreatedAt,
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

func (r *pharmacyDrugRepositoryImpl) SelectAll(ctx context.Context, pharmacyID uint, managerID uint, clc *entity.Collection) ([]entity.PharmacyDrug, error) {
	selectColumns := `
		pd.pharmacy_drug_id,
		pd.drug_id,
		pd.category_id,
		pd.pharmacy_id,
		pd.stock,
		pd.price,
		pd.is_active,
		d.drug_id,
		d.drug_name, 
		d.generic_name,
		d.composition,
		d.description,
		d.classification,
		d.form,
		d.unit_in_pack,
		d.selling_unit,
		d.weight,
		d.height,
		d.length,
		d.width,
		d.image_url,
		d.is_active,
		c.category_id,
		c.category_name,
		c.created_at
	`

	advanceQuery := `
		pharmacy_drugs pd 
			LEFT JOIN drugs d ON pd.drug_id  = d.drug_id 
			LEFT JOIN categories c ON c.category_id = pd.category_id
		WHERE
		%s
		%s
	`

	clc.Args = append(clc.Args, pharmacyID, managerID)

	search := utils.BuildSearchQuery(pdSearchColumn, clc)
	orderBy := utils.BuildSortQuery(pdColumnAlias, clc.Sort, "pd.pharmacy_drug_id DESC")
	filter := utils.BuildFilterQuery(pdColumnAlias, clc, `
		pd.pharmacy_id = $1
		AND
		pd.pharmacy_id IN (SELECT pharmacy_id FROM pharmacies WHERE pharmacy_manager_id = $2)
	`)

	query := utils.BuildQuery(r.db, utils.PaginateQuery{
		SelectColumns: selectColumns,
		AdvanceQuery:  fmt.Sprintf(advanceQuery, filter, search),
		OrderQuery:    orderBy,
	}, clc)

	drugs := []entity.PharmacyDrug{}
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
		var scan entity.PharmacyDrug
		err := rows.Scan(
			&scan.ID,
			&scan.DrugID,
			&scan.CategoryID,
			&scan.PharmacyID,
			&scan.Stock,
			&scan.Price,
			&scan.IsActive,
			&scan.Drug.ID,
			&scan.Drug.Name,
			&scan.Drug.GenericName,
			&scan.Drug.Composition,
			&scan.Drug.Description,
			&scan.Drug.Classification,
			&scan.Drug.Form,
			&scan.Drug.UnitInPack,
			&scan.Drug.SellingUnit,
			&scan.Drug.Weight,
			&scan.Drug.Height,
			&scan.Drug.Length,
			&scan.Drug.Width,
			&scan.Drug.ImageURL,
			&scan.Drug.IsActive,
			&scan.Category.ID,
			&scan.Category.Name,
			&scan.Category.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		drugs = append(drugs, scan)
	}

	return drugs, nil
}

func (r *pharmacyDrugRepositoryImpl) SelectPharmacyDrugsByCategoryId(ctx context.Context, categoryId uint) (int, error) {
	q := `select pharmacy_drug_id from pharmacy_drugs
		where category_id=$1
		And deleted_at is null`
	rows, err := r.db.QueryContext(ctx, q, categoryId)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	defer rows.Close()
	var pharmacyIds []int
	var pharmacyId int
	for rows.Next() {
		err := rows.Scan(
			&pharmacyId,
		)
		if err != nil {
			logrus.Error(err)
			return 0, err
		}
		pharmacyIds = append(pharmacyIds, pharmacyId)
	}

	return len(pharmacyIds), nil
}

func (r *pharmacyDrugRepositoryImpl) SelectOneByID(ctx context.Context, id uint) (*entity.PharmacyDrug, error) {

	q := `
	SELECT
		pharmacy_drug_id,
		drug_id,
		pharmacy_id,
		category_id,
		stock,
		price,
		is_active,
		created_at
	FROM
		pharmacy_drugs
	WHERE
		pharmacy_drug_id = $1`

	var scan entity.PharmacyDrug
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&scan.ID,
		&scan.DrugID,
		&scan.PharmacyID,
		&scan.CategoryID,
		&scan.Stock,
		&scan.Price,
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

	return &scan, nil
}

func (r *pharmacyDrugRepositoryImpl) SelectById(ctx context.Context, id uint) (*entity.PharmacyDrug, error) {
	q := `
	SELECT
		p.pharmacy_id,
		p.pharmacy_name,
		p.open_time,
		p.close_time,
		p.address,
		p.subdistrict_id,
		ST_AsEWKT(p.pharmacy_location),
		pd.pharmacy_drug_id,
		pd.stock,
		pd.price,
		pd.category_id,
		pd.is_active,
		d.drug_id,
		d.manufacturer_id,
		m.manufacturer_name,
		d.drug_name,
		d.generic_name,
		d.composition,
		d.description,
		d.classification,
		d.form,
		d.unit_in_pack,
		d.selling_unit,
		d.weight,
		d.height,
		d.length,
		d.width,
		d.image_url,
		d.is_active
	FROM
		pharmacies p
		JOIN pharmacy_drugs pd ON pd.pharmacy_id = p.pharmacy_id
		JOIN drugs d ON pd.drug_id = d.drug_id
		JOIN manufacturers m ON m.manufacturer_id = d.manufacturer_id
	WHERE
		pd.pharmacy_drug_id = $1;`

	var scan entity.PharmacyDrug
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&scan.Pharmacy.ID,
		&scan.Pharmacy.Name,
		&scan.Pharmacy.OpenTime,
		&scan.Pharmacy.CloseTime,
		&scan.Pharmacy.Address,
		&scan.Pharmacy.SubdistrictID,
		&scan.Pharmacy.Location,
		&scan.ID,
		&scan.Stock,
		&scan.Price,
		&scan.Category.ID,
		&scan.IsActive,
		&scan.Drug.ID,
		&scan.Drug.Manufacturer.ID,
		&scan.Drug.Manufacturer.Name,
		&scan.Drug.Name,
		&scan.Drug.GenericName,
		&scan.Drug.Composition,
		&scan.Drug.Description,
		&scan.Drug.Classification,
		&scan.Drug.Form,
		&scan.Drug.UnitInPack,
		&scan.Drug.SellingUnit,
		&scan.Drug.Weight,
		&scan.Drug.Height,
		&scan.Drug.Length,
		&scan.Drug.Width,
		&scan.Drug.ImageURL,
		&scan.Drug.IsActive,
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

func (r *pharmacyDrugRepositoryImpl) UpdateOne(ctx context.Context, pharmacyDrug entity.PharmacyDrug, id uint) error {
	q := `
		UPDATE 
			pharmacy_drugs
		SET 

			drug_id = $1,
			pharmacy_id= $2,
			category_id= $3,
			stock= $4,
			price= $5,
			is_active= $6
		WHERE
			pharmacy_drug_id = $7
		RETURNING
			drug_id,
			pharmacy_id,
			category_id,
			stock,
			price,
			is_active	
		`
	var scan entity.PharmacyDrug
	err := r.db.QueryRowContext(ctx, q,
		pharmacyDrug.DrugID,
		pharmacyDrug.PharmacyID,
		pharmacyDrug.CategoryID,
		pharmacyDrug.Stock,
		pharmacyDrug.Price,
		pharmacyDrug.IsActive,
		id,
	).Scan(
		&scan.DrugID,
		&scan.PharmacyID,
		&scan.CategoryID,
		&scan.Stock,
		&scan.Price,
		&scan.IsActive,
	)

	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (r *pharmacyDrugRepositoryImpl) CreateOne(ctx context.Context, pharmacyDrug entity.PharmacyDrug) (*entity.PharmacyDrug, error) {
	q := `
		INSERT INTO pharmacy_drugs 
			(drug_id,
			pharmacy_id,
			category_id,
			stock,
			price,
			is_active)
		VALUES
			($1, $2, $3,$4,$5,$6)
		RETURNING
			pharmacy_drug_id,
			drug_id,
			pharmacy_id,
			category_id,
			stock,
			price,
			is_active	
		`
	var scan entity.PharmacyDrug
	err := r.db.QueryRowContext(ctx, q,
		pharmacyDrug.DrugID,
		pharmacyDrug.PharmacyID,
		pharmacyDrug.CategoryID,
		pharmacyDrug.Stock,
		pharmacyDrug.Price,
		pharmacyDrug.IsActive,
	).Scan(
		&scan.ID,
		&scan.DrugID,
		&scan.PharmacyID,
		&scan.CategoryID,
		&scan.Stock,
		&scan.Price,
		&scan.IsActive,
	)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}

func (r *pharmacyDrugRepositoryImpl) UpdateStockComeOut(ctx context.Context, quantity uint, pharmacyDrugId uint) (*uint, error) {
	var stock uint
	q := `
	update pharmacy_drugs
	set 
	stock =stock - $1,
	updated_at=now()
	where pharmacy_drug_id =$2
	returning stock 
	`

	err := r.db.QueryRowContext(ctx, q, quantity, pharmacyDrugId).Scan(&stock)
	if err != nil {
		return nil, err
	}
	return &stock, nil
}
func (r *pharmacyDrugRepositoryImpl) UpdateReturnStock(ctx context.Context, orderDetails []*entity.OrderDetail) ([]entity.StockJurnal, error) {
	q := `
	update pharmacy_drugs
	set stock =stock + $1,
	updated_at=now()
	where pharmacy_drug_id =$2
	and deleted_at is null
	returning drug_id
	`
	stmt, err := r.db.PrepareContext(ctx, q)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer stmt.Close()
	var stockJournals []entity.StockJurnal

	for _, orderDetail := range orderDetails {
		var stockJurnal entity.StockJurnal
		err := stmt.QueryRowContext(ctx, orderDetail.Quantity, orderDetail.PharmacyDrugId).Scan(&stockJurnal.DrugId)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		stockJurnal.PharmacyId = orderDetail.PharmacyDrug.PharmacyID
		stockJurnal.Quantity = int(orderDetail.Quantity)
		stockJurnal.Description = constant.ReturnedStock
		stockJournals = append(stockJournals, stockJurnal)
	}

	return stockJournals, nil
}
func (r *pharmacyDrugRepositoryImpl) UpdateSubstractionBulkStock(ctx context.Context, stockRequestDrugs map[uint][]*entity.StockRequestDrug) error {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE pharmacy_drugs SET stock = stock - $1, updated_at = NOW() WHERE drug_id = $2 AND pharmacy_id = $3 AND stock>= $1 ")
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer stmt.Close()

	for pharmacyId, stockRequest := range stockRequestDrugs {
		for _, stockReq := range stockRequest {
			result, err := stmt.ExecContext(ctx, stockReq.Quantity, stockReq.DrugId, pharmacyId)
			if err != nil {
				logrus.Error(err)
				return err
			}
			if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
				return apperror.InsufficientStock
			}
		}

	}

	return nil

}
func (r *pharmacyDrugRepositoryImpl) UpdateAdditionBulkStock(ctx context.Context, stockRequestDrugs map[uint][]*entity.StockRequestDrug) error {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE pharmacy_drugs SET stock = stock + $1, updated_at = NOW() WHERE drug_id = $2 AND pharmacy_id = $3")
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer stmt.Close()

	for pharmacyId, stockRequest := range stockRequestDrugs {
		for _, stockReq := range stockRequest {
			result, err := stmt.ExecContext(ctx, stockReq.Quantity, stockReq.DrugId, pharmacyId)
			if err != nil {
				logrus.Error(err)
				return err
			}
			if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
				return apperror.ResourceNotFound
			}
		}

	}

	return nil

}

func (r *pharmacyDrugRepositoryImpl) SelectOneNearestByPharmacyManagerDrugId(ctx context.Context, pharmacyDrug entity.PharmacyDrug, quantity uint) (*entity.PharmacyDrug, error) {
	q := `
		SELECT
			p.pharmacy_id,
			ST_Distance(
				ST_MakePoint($1, $2) :: geography,
				p.pharmacy_location
			) AS distance,
			ST_AsEWKT(p.pharmacy_location),
			pd.pharmacy_drug_id,
			pd.stock,
			d.drug_id
		FROM
			pharmacies p
			JOIN pharmacy_drugs pd ON pd.pharmacy_id = p.pharmacy_id
			JOIN drugs d ON pd.drug_id = d.drug_id
			JOIN manufacturers m ON m.manufacturer_id = d.manufacturer_id
		where  
			p.pharmacy_manager_id=$3
			and 
			d.drug_id=$4
			and 
			pd.stock >=$5
			and
			pd.is_active is true
			and
			d.is_active is true
		ORDER BY
			distance
		FOR UPDATE;
	`
	err := r.db.QueryRowContext(ctx, q, pharmacyDrug.Pharmacy.Longitude, pharmacyDrug.Pharmacy.Latitude, pharmacyDrug.Pharmacy.ManagerID, pharmacyDrug.DrugID, quantity).Scan(
		&pharmacyDrug.PharmacyID,
		&pharmacyDrug.Pharmacy.Distance,
		&pharmacyDrug.Pharmacy.Location,
		&pharmacyDrug.ID,
		&pharmacyDrug.Stock,
		&pharmacyDrug.DrugID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logrus.Error(err)
			return nil, apperror.InsufficientStockMutation
		}
	}
	return &pharmacyDrug, nil
}

func (r *pharmacyDrugRepositoryImpl) SelectNearestPharmaciesByDrugID(ctx context.Context, drugID uint, addr entity.Address, radius int, clc *entity.Collection) ([]entity.PharmacyDrug, error) {
	selectColumns := `
		pd.pharmacy_drug_id,
		pd.drug_id,
		pd.pharmacy_id,
		pd.category_id,
		pd.stock,
		pd.price,
		pd.is_active,
		pd.created_at,
		p.pharmacy_id,
		p.pharmacy_name,
		ST_AsEWKT(p.pharmacy_location),
		p.operation_day,
		p.open_time,
		p.close_time,
		p.created_at
	`
	advanceQuery := `
			pharmacy_drugs pd
		LEFT JOIN pharmacies p ON p.pharmacy_id = pd.pharmacy_id
		WHERE
			ST_DWithin(
				p.pharmacy_location,
				ST_MakePoint($1, $2) :: geography,
				$3
			)
		AND
			pd.drug_id = $4
		AND
		%s
		%s
	`

	clc.Args = append(clc.Args, addr.Longitude, addr.Latitude, radius, drugID)

	search := utils.BuildSearchQuery([]string{}, clc)
	orderBy := utils.BuildSortQuery(map[string]string{}, clc.Sort, `
		ST_Distance(
			ST_MakePoint($1, $2) :: geography,
			p.pharmacy_location
		) ASC
	`)

	filter := utils.BuildFilterQuery(map[string]string{}, clc, "pd.deleted_at is null")

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

	results := make([]entity.PharmacyDrug, 0)
	for rows.Next() {
		var scan entity.PharmacyDrug
		if err := rows.Scan(
			&scan.ID,
			&scan.DrugID,
			&scan.PharmacyID,
			&scan.CategoryID,
			&scan.Stock,
			&scan.Price,
			&scan.IsActive,
			&scan.CreatedAt,
			&scan.Pharmacy.ID,
			&scan.Pharmacy.Name,
			&scan.Pharmacy.Location,
			&scan.Pharmacy.OperationDay,
			&scan.Pharmacy.OpenTime,
			&scan.Pharmacy.CloseTime,
			&scan.Pharmacy.CreatedAt,
		); err != nil {
			logrus.Error(err)
			return nil, err
		}

		results = append(results, scan)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *pharmacyDrugRepositoryImpl) SelectMostBought(ctx context.Context, limit uint, isTheDay bool) ([]entity.PharmacyDrug, error) {
	addsFilter := ""
	if isTheDay {
		addsFilter = "AND DATE(p.updated_at) = CURRENT_DATE"
	}

	selectColumns := `
		od.pharmacy_drug_id,
		pd.pharmacy_id,
		pd.stock,
		pd.price,
		pd.is_active,
		p.pharmacy_id,
		p.pharmacy_name,
		p.open_time,
		p.close_time,
		p.address,
		p.subdistrict_id,
		ST_AsEWKT(p.pharmacy_location),
		d.drug_id,
		d.manufacturer_id,
		d.drug_name,
		d.generic_name,
		d.composition,
		d.description,
		d.classification,
		d.form,
		d.unit_in_pack,
		d.selling_unit,
		d.weight,
		d.height,
		d.length,
		d.width,
		d.image_url,
		d.is_active
	`

	q := fmt.Sprintf(`
		SELECT
			%s,
			sum(od.quantity)
		FROM
			order_details od
		LEFT JOIN pharmacy_drugs pd ON pd.pharmacy_drug_id = od.pharmacy_drug_id
		LEFT JOIN pharmacies p ON p.pharmacy_id  = pd.pharmacy_id
		LEFT JOIN drugs d ON d.drug_id = pd.drug_id
		WHERE
			od.order_id IN (
				SELECT o.order_id FROM orders o WHERE o.payment_id IN (
					SELECT p.payment_id FROM payments p WHERE 
						p.payment_expired_at IS NULL AND p.deleted_at IS NULL %s
				)
			)
			AND pd.is_active = true 
			AND d.is_active = true
		GROUP BY %s
		ORDER BY sum(od.quantity) DESC
		LIMIT %d
	`, selectColumns, addsFilter, selectColumns, limit)

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	results := make([]entity.PharmacyDrug, 0)
	for rows.Next() {
		var scan entity.PharmacyDrug
		var scanUint uint
		if err := rows.Scan(
			&scan.ID,
			&scan.PharmacyID,
			&scan.Stock,
			&scan.Price,
			&scan.IsActive,
			&scan.Pharmacy.ID,
			&scan.Pharmacy.Name,
			&scan.Pharmacy.OpenTime,
			&scan.Pharmacy.CloseTime,
			&scan.Pharmacy.Address,
			&scan.Pharmacy.SubdistrictID,
			&scan.Pharmacy.Location,
			&scan.Drug.ID,
			&scan.Drug.ManufacturerID,
			&scan.Drug.Name,
			&scan.Drug.GenericName,
			&scan.Drug.Composition,
			&scan.Drug.Description,
			&scan.Drug.Classification,
			&scan.Drug.Form,
			&scan.Drug.UnitInPack,
			&scan.Drug.SellingUnit,
			&scan.Drug.Weight,
			&scan.Drug.Height,
			&scan.Drug.Length,
			&scan.Drug.Width,
			&scan.Drug.ImageURL,
			&scan.Drug.IsActive,
			&scanUint,
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

func (r *pharmacyDrugRepositoryImpl) SelectAllWithLimit(ctx context.Context, limit uint) ([]entity.PharmacyDrug, error) {
	q := fmt.Sprintf(`
		SELECT
			pd.pharmacy_drug_id,
			pd.pharmacy_id,
			pd.category_id,
			pd.stock,
			pd.price,
			pd.is_active,
			p.pharmacy_id,
			p.pharmacy_name,
			p.open_time,
			p.close_time,
			p.address,
			p.subdistrict_id,
			ST_AsEWKT(p.pharmacy_location),
			d.drug_id,
			d.manufacturer_id,
			d.drug_name,
			d.generic_name,
			d.composition,
			d.description,
			d.classification,
			d.form,
			d.unit_in_pack,
			d.selling_unit,
			d.weight,
			d.height,
			d.length,
			d.width,
			d.image_url,
			d.is_active
		FROM
			pharmacy_drugs pd
		LEFT JOIN pharmacies p ON p.pharmacy_id  = pd.pharmacy_id
		LEFT JOIN drugs d ON d.drug_id = pd.drug_id
		WHERE pd.is_active = true AND d.is_active = true
		ORDER BY pd.stock DESC
		LIMIT %d
	`, limit)

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	results := make([]entity.PharmacyDrug, 0)
	for rows.Next() {
		var scan entity.PharmacyDrug
		if err := rows.Scan(
			&scan.ID,
			&scan.PharmacyID,
			&scan.CategoryID,
			&scan.Stock,
			&scan.Price,
			&scan.IsActive,
			&scan.Pharmacy.ID,
			&scan.Pharmacy.Name,
			&scan.Pharmacy.OpenTime,
			&scan.Pharmacy.CloseTime,
			&scan.Pharmacy.Address,
			&scan.Pharmacy.SubdistrictID,
			&scan.Pharmacy.Location,
			&scan.Drug.ID,
			&scan.Drug.ManufacturerID,
			&scan.Drug.Name,
			&scan.Drug.GenericName,
			&scan.Drug.Composition,
			&scan.Drug.Description,
			&scan.Drug.Classification,
			&scan.Drug.Form,
			&scan.Drug.UnitInPack,
			&scan.Drug.SellingUnit,
			&scan.Drug.Weight,
			&scan.Drug.Height,
			&scan.Drug.Length,
			&scan.Drug.Width,
			&scan.Drug.ImageURL,
			&scan.Drug.IsActive,
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

func (r *pharmacyDrugRepositoryImpl) SelectOneByPharmacyAndDrugID(ctx context.Context, pd entity.PharmacyDrug) (*entity.PharmacyDrug, error) {
	q := `
		SELECT
			pharmacy_drug_id,
			drug_id,
			pharmacy_id,
			category_id,
			stock,
			price,
			is_active,
			created_at
		FROM
			pharmacy_drugs
		WHERE
			pharmacy_id = $1
		AND
			drug_id = $2
	`

	var scan entity.PharmacyDrug
	err := r.db.QueryRowContext(ctx, q, pd.PharmacyID, pd.DrugID).Scan(
		&scan.ID,
		&scan.DrugID,
		&scan.PharmacyID,
		&scan.CategoryID,
		&scan.Stock,
		&scan.Price,
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

	return &scan, nil
}
