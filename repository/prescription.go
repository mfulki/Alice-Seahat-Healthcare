package repository

import (
	"context"
	"fmt"
	"strings"

	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"

	"github.com/sirupsen/logrus"
)

type PrescriptionRepository interface {
	GetAllByTelemedicineID(ctx context.Context, telemedicineID uint) ([]entity.Prescription, error)
	InsertMany(ctx context.Context, pss []entity.Prescription) ([]entity.Prescription, error)
}

type prescriptionRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewPrescriptionRepository(db transaction.DBTransaction) *prescriptionRepositoryImpl {
	return &prescriptionRepositoryImpl{
		db: db,
	}
}

func (r *prescriptionRepositoryImpl) GetAllByTelemedicineID(ctx context.Context, telemedicineID uint) ([]entity.Prescription, error) {
	q := `
		SELECT 
			p.prescription_id,
			p.telemedicine_id,
			p.drug_id,
			p.quantity,
			p.notes,
			p.created_at,
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
			d.is_active,
			d.created_at
		FROM
			prescriptions p
		LEFT JOIN
			drugs d ON d.drug_id = p.drug_id
		WHERE
			telemedicine_id = $1
	`

	rows, err := r.db.QueryContext(ctx, q, telemedicineID)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	results := make([]entity.Prescription, 0)
	for rows.Next() {
		var scan entity.Prescription
		if err := rows.Scan(&scan.ID,
			&scan.TelemedicineID,
			&scan.DrugID,
			&scan.Quantity,
			&scan.Notes,
			&scan.CreatedAt,
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
			&scan.Drug.CreatedAt,
		); err != nil {
			logrus.Error(err)
			return nil, err
		}

		results = append(results, scan)
	}

	return results, nil
}

func (r *prescriptionRepositoryImpl) InsertMany(ctx context.Context, pss []entity.Prescription) ([]entity.Prescription, error) {
	args := make([]any, 0)
	argsLen := len(args)
	q := `
		INSERT INTO
			prescriptions (telemedicine_id, drug_id, quantity, notes)
		VALUES
			%s
		RETURNING
			prescription_id, telemedicine_id, drug_id, quantity, notes, created_at
	`

	insertData := make([]string, 0)
	for _, ps := range pss {
		insertData = append(insertData, fmt.Sprintf("($%d, $%d, $%d, $%d)", argsLen+1, argsLen+2, argsLen+3, argsLen+4))
		args = append(args, ps.TelemedicineID, ps.DrugID, ps.Quantity, ps.Notes)
		argsLen = len(args)
	}

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(q, strings.Join(insertData, ",")), args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	results := make([]entity.Prescription, 0)
	for rows.Next() {
		var scan entity.Prescription
		if err := rows.Scan(&scan.ID, &scan.TelemedicineID, &scan.DrugID, &scan.Quantity, &scan.Notes, &scan.CreatedAt); err != nil {
			logrus.Error(err)
			return nil, err
		}

		results = append(results, scan)
	}

	return results, nil
}
