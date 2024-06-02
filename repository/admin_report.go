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
	reportColumnAlias = map[string]string{
		"total_price": "sum(od.quantity)",
		"pharmacy_id": "p.pharmacy_id",
	}
	reportDrugSearchColumn = []string{
		"d.drug_name",
		"m.manufacturer_name",
	}
	reportCategorySearchColumn = []string{
		"c.category_name",
	}
)

type AdminReportRepository interface {
	GetDrugReportByCurrentMonth(ctx context.Context, managerID uint, clc *entity.Collection) ([]entity.AdminReportByDrug, error)
	GetCategoryReportByCurrentMonth(ctx context.Context, managerID uint, clc *entity.Collection) ([]entity.AdminReportByCategory, error)
}

type adminReportRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewAdminReportRepository(db transaction.DBTransaction) *adminReportRepositoryImpl {
	return &adminReportRepositoryImpl{
		db: db,
	}
}

func (r *adminReportRepositoryImpl) GetDrugReportByCurrentMonth(ctx context.Context, managerID uint, clc *entity.Collection) ([]entity.AdminReportByDrug, error) {
	selectColumns := `
		d.drug_id,
		d.drug_name,
		m.manufacturer_name,
		d.classification,
		d.form,
		d.unit_in_pack,
		d.selling_unit,
		d.image_url,
		sum(od.quantity) as total_quantity,
		sum(od.price) as total_price
	`

	advanceQuery := `
		order_details od
		INNER JOIN pharmacy_drugs pd ON pd.pharmacy_drug_id = od.pharmacy_drug_id
		INNER JOIN pharmacies p ON p.pharmacy_id = pd.pharmacy_id
		INNER JOIN drugs d ON d.drug_id = pd.drug_id
		INNER JOIN manufacturers m ON m.manufacturer_id = d.manufacturer_id
		WHERE
			od.order_id IN (
				SELECT o.order_id FROM orders o WHERE o.payment_id IN (
					SELECT p.payment_id FROM payments p 
					WHERE 
						p.payment_expired_at IS NULL AND p.deleted_at IS NULL
					AND
						p.updated_at BETWEEN
							date_trunc('month', current_date)
						AND
							date_trunc('month', current_date) + interval '1 month' - interval '1 day'
				)
			)
			%s
			AND
			%s
			%s
		GROUP BY
			d.drug_id,
			d.drug_name,
			m.manufacturer_name,
			d.classification,
			d.form,
			d.unit_in_pack,
			d.selling_unit,
			d.image_url
	`

	extendQuery := ""
	if managerID != 0 {
		extendQuery = "AND p.pharmacy_manager_id = $1"
		clc.Args = append(clc.Args, managerID)
	}

	search := utils.BuildSearchQuery(reportDrugSearchColumn, clc)
	orderBy := utils.BuildSortQuery(reportColumnAlias, clc.Sort, "sum(od.quantity) desc")
	filter := utils.BuildFilterQuery(reportColumnAlias, clc, "od.deleted_at IS NULL")

	query := utils.BuildQueryWithGroupBy(r.db, utils.PaginateQuery{
		SelectColumns: selectColumns,
		AdvanceQuery:  fmt.Sprintf(advanceQuery, extendQuery, filter, search),
		OrderQuery:    orderBy,
	}, clc)

	reports := make([]entity.AdminReportByDrug, 0)
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
		report := entity.AdminReportByDrug{}
		err := rows.Scan(
			&report.DrugID,
			&report.DrugName,
			&report.ManufacturerName,
			&report.Classification,
			&report.Form,
			&report.UnitInPack,
			&report.SellingUnit,
			&report.ImageURL,
			&report.TotalQuantiy,
			&report.TotalPrice,
		)

		if err != nil {
			return nil, err
		}

		reports = append(reports, report)
	}

	return reports, nil
}

func (r *adminReportRepositoryImpl) GetCategoryReportByCurrentMonth(ctx context.Context, managerID uint, clc *entity.Collection) ([]entity.AdminReportByCategory, error) {
	selectColumns := `
		c.category_id,
		c.category_name,
		sum(od.quantity) as total_quantity,
		sum(od.price) as total_price
	`

	advanceQuery := `
		order_details od
		INNER JOIN pharmacy_drugs pd ON pd.pharmacy_drug_id = od.pharmacy_drug_id
		INNER JOIN pharmacies p ON p.pharmacy_id = pd.pharmacy_id
		INNER JOIN categories c ON c.category_id = pd.category_id
		WHERE
			od.order_id IN (
				SELECT o.order_id FROM orders o WHERE o.payment_id IN (
					SELECT p.payment_id FROM payments p 
					WHERE 
						p.payment_expired_at IS NULL AND p.deleted_at IS NULL
					AND
						p.updated_at BETWEEN
							date_trunc('month', current_date)
						AND
							date_trunc('month', current_date) + interval '1 month' - interval '1 day'
				)
			)
			%s
			AND
			%s
			%s
		GROUP BY
			c.category_id,
			c.category_name
	`

	extendQuery := ""
	if managerID != 0 {
		extendQuery = "AND p.pharmacy_manager_id = $1"
		clc.Args = append(clc.Args, managerID)
	}

	search := utils.BuildSearchQuery(reportCategorySearchColumn, clc)
	orderBy := utils.BuildSortQuery(reportColumnAlias, clc.Sort, "sum(od.quantity) desc")
	filter := utils.BuildFilterQuery(reportColumnAlias, clc, "od.deleted_at IS NULL")

	query := utils.BuildQueryWithGroupBy(r.db, utils.PaginateQuery{
		SelectColumns: selectColumns,
		AdvanceQuery:  fmt.Sprintf(advanceQuery, extendQuery, filter, search),
		OrderQuery:    orderBy,
	}, clc)

	reports := make([]entity.AdminReportByCategory, 0)
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
		report := entity.AdminReportByCategory{}
		err := rows.Scan(
			&report.CategoryID,
			&report.CategoryName,
			&report.TotalQuantiy,
			&report.TotalPrice,
		)

		if err != nil {
			return nil, err
		}

		reports = append(reports, report)
	}

	return reports, nil
}
