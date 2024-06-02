package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"

	"github.com/sirupsen/logrus"
)

type CartItemRepository interface {
	GetAllWithOtherDetail(ctx context.Context, userID uint) ([]entity.CartItem, error)
	SelectOneByID(ctx context.Context, item entity.CartItem) (*entity.CartItem, error)
	InsertOne(ctx context.Context, item entity.CartItem) (*entity.CartItem, error)
	GetByPrescriptedAndPharmacyDrugID(ctx context.Context, item entity.CartItem) (*entity.CartItem, error)
	UpdateQuantityByID(ctx context.Context, item entity.CartItem) error
	DeleteManyByID(ctx context.Context, ids []uint, userID uint) error
	LockRow(ctx context.Context, cart []*entity.CartItem, userId uint, pharmacyId uint) ([]*entity.CartItem, error)
}

type cartItemRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewCartItemRepository(db transaction.DBTransaction) *cartItemRepositoryImpl {
	return &cartItemRepositoryImpl{
		db: db,
	}
}

func (r *cartItemRepositoryImpl) GetAllWithOtherDetail(ctx context.Context, userID uint) ([]entity.CartItem, error) {
	q := `
		SELECT
			c.cart_item_id,
			c.user_id, 
			c.pharmacy_drug_id, 
			c.quantity, 
			c.is_prescripted,
			c.created_at,
			pd.pharmacy_drug_id,
			pd.stock,
			pd.price,
			p.pharmacy_id,
			p.pharmacy_name,
			d.drug_id,
			d.drug_name,
			d.generic_name,
			d.composition,
			d.description,
			d.classification,
			d.form,
			d.unit_in_pack,
			d.selling_unit,	
			d.image_url,
			d.is_active,
			d.weight,
			d.height,
			d.length,
			d.width,
			d.created_at
		FROM
			cart_items c
		LEFT JOIN pharmacy_drugs pd ON c.pharmacy_drug_id = pd.pharmacy_drug_id
		LEFT JOIN pharmacies p ON pd.pharmacy_id = p.pharmacy_id
		LEFT JOIN drugs d ON pd.drug_id = d.drug_id
		WHERE
			c.user_id = $1
		AND
			c.deleted_at IS NULL 
		ORDER BY
			c.is_prescripted DESC,
			c.created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	results := make([]entity.CartItem, 0)
	for rows.Next() {
		var scan entity.CartItem
		if err := rows.Scan(
			&scan.ID,
			&scan.UserID,
			&scan.PharmacyDrugID,
			&scan.Quantity,
			&scan.IsPrescripted,
			&scan.CreatedAt,
			&scan.PharmacyDrug.ID,
			&scan.PharmacyDrug.Stock,
			&scan.PharmacyDrug.Price,
			&scan.PharmacyDrug.Pharmacy.ID,
			&scan.PharmacyDrug.Pharmacy.Name,
			&scan.PharmacyDrug.Drug.ID,
			&scan.PharmacyDrug.Drug.Name,
			&scan.PharmacyDrug.Drug.GenericName,
			&scan.PharmacyDrug.Drug.Composition,
			&scan.PharmacyDrug.Drug.Description,
			&scan.PharmacyDrug.Drug.Classification,
			&scan.PharmacyDrug.Drug.Form,
			&scan.PharmacyDrug.Drug.UnitInPack,
			&scan.PharmacyDrug.Drug.SellingUnit,
			&scan.PharmacyDrug.Drug.ImageURL,
			&scan.PharmacyDrug.Drug.IsActive,
			&scan.PharmacyDrug.Drug.Weight,
			&scan.PharmacyDrug.Drug.Height,
			&scan.PharmacyDrug.Drug.Length,
			&scan.PharmacyDrug.Drug.Width,
			&scan.PharmacyDrug.Drug.CreatedAt,
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

func (r *cartItemRepositoryImpl) SelectOneByID(ctx context.Context, item entity.CartItem) (*entity.CartItem, error) {
	q := `
		SELECT
			cart_item_id,
			user_id, 
			pharmacy_drug_id, 
			quantity, 
			is_prescripted,
			created_at
		FROM
			cart_items
		WHERE
			cart_item_id = $1
		AND
			user_id = $2
		AND
			deleted_at IS NULL
	`

	var scan entity.CartItem
	if err := r.db.QueryRowContext(ctx, q,
		item.ID,
		item.UserID,
	).Scan(
		&scan.ID,
		&scan.UserID,
		&scan.PharmacyDrugID,
		&scan.Quantity,
		&scan.IsPrescripted,
		&scan.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrResourceNotFound
		}

		return nil, err
	}

	return &scan, nil
}

func (r *cartItemRepositoryImpl) InsertOne(ctx context.Context, item entity.CartItem) (*entity.CartItem, error) {
	q := `
		INSERT INTO
			cart_items (user_id, pharmacy_drug_id, quantity, is_prescripted)
		VALUES
			($1, $2, $3, $4)
		RETURNING
			cart_item_id,
			user_id, 
			pharmacy_drug_id, 
			quantity, 
			is_prescripted,
			created_at
	`

	var scan entity.CartItem
	if err := r.db.QueryRowContext(ctx, q,
		item.UserID,
		item.PharmacyDrugID,
		item.Quantity,
		item.IsPrescripted,
	).Scan(
		&scan.ID,
		&scan.UserID,
		&scan.PharmacyDrugID,
		&scan.Quantity,
		&scan.IsPrescripted,
		&scan.CreatedAt,
	); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &scan, nil
}

func (r *cartItemRepositoryImpl) GetByPrescriptedAndPharmacyDrugID(ctx context.Context, item entity.CartItem) (*entity.CartItem, error) {
	q := `
		SELECT
			cart_item_id,
			user_id, 
			pharmacy_drug_id, 
			quantity, 
			is_prescripted,
			created_at
		FROM
			cart_items
		WHERE
			pharmacy_drug_id = $1
		AND
			user_id = $2
		AND
			is_prescripted = $3
		AND
			deleted_at IS NULL
	`

	var scan entity.CartItem
	if err := r.db.QueryRowContext(ctx, q,
		item.PharmacyDrugID,
		item.UserID,
		item.IsPrescripted,
	).Scan(
		&scan.ID,
		&scan.UserID,
		&scan.PharmacyDrugID,
		&scan.Quantity,
		&scan.IsPrescripted,
		&scan.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrResourceNotFound
		}

		return nil, err
	}

	return &scan, nil
}

func (r *cartItemRepositoryImpl) UpdateQuantityByID(ctx context.Context, item entity.CartItem) error {
	q := `
		UPDATE
			cart_items
		SET
			quantity = $1,
			updated_at = current_timestamp
		WHERE
			cart_item_id = $2
		AND
			user_id = $3
		AND
			deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, q, item.Quantity, item.ID, item.UserID)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return apperror.ErrResourceNotFound
	}

	return nil
}

func (r *cartItemRepositoryImpl) DeleteManyByID(ctx context.Context, ids []uint, userID uint) error {
	paramBuilder := new(strings.Builder)
	paramBuilder.WriteString("(")

	for index, id := range ids {
		paramBuilder.WriteString(strconv.Itoa(int(id)))

		if index != len(ids)-1 {
			paramBuilder.WriteString(",")
		}
	}

	paramBuilder.WriteString(")")

	q := fmt.Sprintf(`
		UPDATE
			cart_items
		SET
			updated_at = CURRENT_TIMESTAMP,	
			deleted_at = CURRENT_TIMESTAMP
		WHERE
			cart_item_id IN %s
		AND
			user_id = $1
		AND
			deleted_at IS NULL
	`, paramBuilder.String())

	if _, err := r.db.ExecContext(ctx, q, userID); err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}
func (r *cartItemRepositoryImpl) LockRow(ctx context.Context, cart []*entity.CartItem, userId uint, pharmacyId uint) ([]*entity.CartItem, error) {
	var s strings.Builder
	args := []any{}
	s.WriteString(`SELECT
					ci.user_id,
					ci.cart_item_id, 
					ci.pharmacy_drug_id, 
					ci.quantity, 
					ci.is_prescripted,
					pd.price,
					pd.pharmacy_id,
					d.weight
					FROM
					cart_items ci  
					join 
						pharmacy_drugs pd on ci.pharmacy_drug_id =pd.pharmacy_drug_id 
					join
						drugs d on pd.drug_id=d.drug_id
					WHERE 
						ci.deleted_at is null 
					AND
						ci.user_id = $1 
					AND
						pd.pharmacy_id=$2 
					AND(`)
	args = append(args, userId, pharmacyId)

	for num, cartItem := range cart {
		args = append(args, cartItem.ID)
		if num != 0 {
			s.WriteString(` OR `)
		}
		s.WriteString(fmt.Sprintf(` cart_item_id =$%s `, strconv.Itoa(num+3)))
	}
	s.WriteString(`) FOR UPDATE`)
	q := s.String()

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()
	numeric := 0
	var validCarts []*entity.CartItem

	for rows.Next() {
		cartItemData := &entity.CartItem{}
		err := rows.Scan(
			&cartItemData.UserID,
			&cartItemData.ID,
			&cartItemData.PharmacyDrugID,
			&cartItemData.Quantity,
			&cartItemData.IsPrescripted,
			&cartItemData.Price,
			&cartItemData.PharmacyDrug.PharmacyID,
			&cartItemData.PharmacyDrug.Drug.Weight,
		)
		numeric++
		cartItemData.TotalPrice = cartItemData.Price * cartItemData.Quantity
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		if cartItemData.PharmacyDrugID != 0 {
			validCarts = append(validCarts, cartItemData)
		}
	}

	if err = rows.Err(); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return validCarts, nil

}
