package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sort"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/utils"

	"github.com/sirupsen/logrus"
)

var (
	srSearchColumn = []string{
		"d.drug_name",
	}
)

type StockRequestRepository interface {
	InsertStockRequest(ctx context.Context, stockRequest entity.StockRequest) (*entity.StockRequest, error)
	InsertStockRequestBulk(ctx context.Context, receiverPharmacy uint, stockRequestDrugs map[uint][]*entity.StockRequestDrug, status string) ([]*entity.StockRequest, error)
	UpdateStockMutationStatusBySender(ctx context.Context, stockrequest entity.StockRequest, updateStatus string, managerId uint) (*entity.StockRequest, error)
	GetAllStockRequestByPharmacyManagerId(ctx context.Context, pharmacyManagerId uint) ([]*entity.StockRequest, error)
	SelectDrugWithSenderAndReceiverPharmacy(ctx context.Context, pharmacySenderID uint, pharmacyReceiverID uint, clc *entity.Collection) ([]entity.DrugWithPharmacyDrug, error)
	GetStockRequestById(ctx context.Context, stockRequestId uint) (map[uint][]*entity.StockRequestDrug, map[uint][]*entity.StockRequestDrug, *entity.StockRequest, error)
}

type stockRequestRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewStockRequestRepository(db transaction.DBTransaction) *stockRequestRepositoryImpl {
	return &stockRequestRepositoryImpl{
		db: db,
	}
}

func (r *stockRequestRepositoryImpl) InsertStockRequest(ctx context.Context, stockRequest entity.StockRequest) (*entity.StockRequest, error) {
	q := `
		insert into stock_requests 
		(sender_pharmacy_id,receiver_pharmacy_id,status)
		values
		($1,$2,$3)
		returning stock_request_id,sender_pharmacy_id,receiver_pharmacy_id,status;
	`
	err := r.db.QueryRowContext(ctx, q, stockRequest.SenderPharmacy.ID, stockRequest.ReceiverPharmacy.ID, stockRequest.Status).Scan(&stockRequest.Id, &stockRequest.SenderPharmacy.ID, &stockRequest.ReceiverPharmacy.ID, &stockRequest.Status)
	if err != nil {
		return nil, err
	}
	return &stockRequest, nil
}
func (r *stockRequestRepositoryImpl) InsertStockRequestBulk(ctx context.Context, receiverPharmacy uint, stockRequestDrugs map[uint][]*entity.StockRequestDrug, status string) ([]*entity.StockRequest, error) {
	q := `
		insert into stock_requests 
		(sender_pharmacy_id,receiver_pharmacy_id,status)
		values
		($1,$2,$3)
		returning stock_request_id,sender_pharmacy_id,receiver_pharmacy_id,status;
	`
	stockRequests := make([]*entity.StockRequest, 0)
	stmt, err := r.db.PrepareContext(ctx, q)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	for senderPharmacyId, stockRequestDrug := range stockRequestDrugs {
		rows, err := stmt.QueryContext(ctx, senderPharmacyId, receiverPharmacy, status)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			stockRequest := entity.StockRequest{}
			err := rows.Scan(
				&stockRequest.Id,
				&stockRequest.SenderPharmacy.ID,
				&stockRequest.ReceiverPharmacy.ID,
				&stockRequest.Status,
			)
			stockRequest.StockRequestDrug = stockRequestDrug
			if err != nil {
				return nil, err
			}
			stockRequests = append(stockRequests, &stockRequest)
		}
	}
	return stockRequests, nil
}

func (r *stockRequestRepositoryImpl) UpdateStockMutationStatusBySender(ctx context.Context, stockrequest entity.StockRequest, updateStatus string, managerId uint) (*entity.StockRequest, error) {
	q := `update stock_requests sr 
		set status =$1
		from
			pharmacies p
		where p.pharmacy_id =sr.sender_pharmacy_id 
		and p.pharmacy_manager_id =$2
		and sr.stock_request_id =$3 
		and status ='waiting for approval'
		returning stock_request_id,sender_pharmacy_id,receiver_pharmacy_id,status;`
	err := r.db.QueryRowContext(ctx, q, updateStatus, managerId, stockrequest.Id).Scan(&stockrequest.Id, &stockrequest.SenderPharmacy.ID, &stockrequest.ReceiverPharmacy.ID, &stockrequest.Status)
	if err != nil {
		logrus.Error(err)
		if err == sql.ErrNoRows {
			if updateStatus == constant.Approved {
				return nil, apperror.CantApproveStockMutation
			}
			if updateStatus == constant.Cancelled {
				return nil, apperror.CantCancelStockMutation
			}
		}
		return nil, err
	}
	return &stockrequest, nil
}

func (r *stockRequestRepositoryImpl) GetAllStockRequestByPharmacyManagerId(ctx context.Context, pharmacyManagerId uint) ([]*entity.StockRequest, error) {
	q := `
		select 
		sr.stock_request_id,
		sr.sender_pharmacy_id ,
		ps.pharmacy_name,
		sr.receiver_pharmacy_id ,
		pr.pharmacy_name,
		sr.status ,
		sr.created_at ,
		sr.updated_at ,
		sr.deleted_at ,
		srd.stock_request_drug_id ,
		srd.stock_request_id ,
		srd.drug_id,
		d.drug_name,
		srd.quantity,
		srd.created_at,
		srd.updated_at,
		srd.deleted_at 
		from stock_requests sr 
		join pharmacies ps on sr.sender_pharmacy_id =ps.pharmacy_id 
		join pharmacies pr on sr.receiver_pharmacy_id =pr.pharmacy_id  
		join stock_request_drugs srd 
		on sr.stock_request_id =srd.stock_request_id
		join drugs d on d.drug_id=srd.drug_id
		where sr.deleted_at is null
		and srd.deleted_at is null
		and pr.pharmacy_manager_id=$1 or ps.pharmacy_manager_id=$1 
		order by srd.stock_request_drug_id  desc ;
		`
	rows, err := r.db.QueryContext(ctx, q, pharmacyManagerId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	sRMap := map[uint]*entity.StockRequest{}
	sRDMap := map[uint][]*entity.StockRequestDrug{}
	stockrequests := []*entity.StockRequest{}
	for rows.Next() {
		sR := &entity.StockRequest{}
		sRD := &entity.StockRequestDrug{}
		err := rows.Scan(
			&sR.Id,
			&sR.SenderPharmacy.ID,
			&sR.SenderPharmacy.Name,
			&sR.ReceiverPharmacy.ID,
			&sR.ReceiverPharmacy.Name,
			&sR.Status,
			&sR.CreatedAt,
			&sR.UpdatedAt,
			&sR.DeletedAt,
			&sRD.Id,
			&sRD.StockRequestId,
			&sRD.DrugId,
			&sRD.Drug.Name,
			&sRD.Quantity,
			&sRD.CreatedAt,
			&sRD.UpdatedAt,
			&sRD.DeletedAt,
		)
		sRD.Drug.ID = sRD.DrugId

		sRDMap[sR.Id] = append(sRDMap[sR.Id], sRD)
		sR.StockRequestDrug = sRDMap[sR.Id]

		sRMap[sR.Id] = sR
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}

	keys := make([]uint, 0, len(sRMap))
	for k := range sRMap {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	for _, k := range keys {
		stockrequests = append(stockrequests, sRMap[k])
	}

	if err = rows.Err(); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return stockrequests, nil
}

func (r *stockRequestRepositoryImpl) GetStockRequestById(ctx context.Context, stockRequestId uint) (map[uint][]*entity.StockRequestDrug, map[uint][]*entity.StockRequestDrug, *entity.StockRequest, error) {
	q := `
		select 
		sr.stock_request_id,
		sr.sender_pharmacy_id ,
		ps.pharmacy_name,
		sr.receiver_pharmacy_id ,
		pr.pharmacy_name,
		sr.status ,
		sr.created_at ,
		sr.updated_at ,
		sr.deleted_at ,
		srd.stock_request_drug_id ,
		srd.stock_request_id ,
		srd.drug_id,
		d.drug_name,
		srd.quantity,
		srd.created_at,
		srd.updated_at,
		srd.deleted_at 
		from stock_requests sr 
		join pharmacies ps on sr.sender_pharmacy_id =ps.pharmacy_id 
		join pharmacies pr on sr.receiver_pharmacy_id =pr.pharmacy_id  
		join stock_request_drugs srd 
		on sr.stock_request_id =srd.stock_request_id
		join drugs d on d.drug_id=srd.drug_id
		where sr.deleted_at is null
		and srd.deleted_at is null
		and sr.stock_request_id=$1 
		order by srd.stock_request_drug_id  desc
		For Update ;
		`

	rows, err := r.db.QueryContext(ctx, q, stockRequestId)
	if err != nil {
		logrus.Error(err)
		return nil, nil, nil, err
	}

	defer rows.Close()

	sRDSenderMap := map[uint][]*entity.StockRequestDrug{}
	sRDReceiverMap := map[uint][]*entity.StockRequestDrug{}
	sR := &entity.StockRequest{}
	for rows.Next() {
		sRD := &entity.StockRequestDrug{}
		err := rows.Scan(
			&sR.Id,
			&sR.SenderPharmacy.ID,
			&sR.SenderPharmacy.Name,
			&sR.ReceiverPharmacy.ID,
			&sR.ReceiverPharmacy.Name,
			&sR.Status,
			&sR.CreatedAt,
			&sR.UpdatedAt,
			&sR.DeletedAt,
			&sRD.Id,
			&sRD.StockRequestId,
			&sRD.DrugId,
			&sRD.Drug.Name,
			&sRD.Quantity,
			&sRD.CreatedAt,
			&sRD.UpdatedAt,
			&sRD.DeletedAt,
		)
		sRD.Drug.ID = sRD.DrugId
		sRDSenderMap[sR.SenderPharmacy.ID] = append(sRDSenderMap[sR.SenderPharmacy.ID], sRD)
		sRDReceiverMap[sR.ReceiverPharmacy.ID] = append(sRDReceiverMap[sR.ReceiverPharmacy.ID], sRD)
		sR.StockRequestDrug = sRDSenderMap[sR.SenderPharmacy.ID]
		if err != nil {
			logrus.Error(err)
			return nil, nil, nil, err
		}
	}

	if err = rows.Err(); err != nil {
		logrus.Error(err)
		return nil, nil, nil, err
	}

	return sRDSenderMap, sRDReceiverMap, sR, nil
}

func (r *stockRequestRepositoryImpl) SelectDrugWithSenderAndReceiverPharmacy(ctx context.Context, pharmacySenderID uint, pharmacyReceiverID uint, clc *entity.Collection) ([]entity.DrugWithPharmacyDrug, error) {
	selectColumns := `
		DISTINCT
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
		d.created_at,
		pdSender.pharmacy_drug_id,
		pdSender.drug_id,
		pdSender.pharmacy_id,
		pdSender.category_id,
		pdSender.stock,
		pdSender.price,
		pdSender.is_active,
		pdSender.created_at,
		pdReceiver.pharmacy_drug_id,
		pdReceiver.drug_id,
		pdReceiver.pharmacy_id,
		pdReceiver.category_id,
		pdReceiver.stock,
		pdReceiver.price,
		pdReceiver.is_active,
		pdReceiver.created_at
	`

	advanceQuery := `
		drugs d
			INNER JOIN pharmacy_drugs pdSender ON pdSender.drug_id = d.drug_id
			INNER JOIN pharmacy_drugs pdReceiver ON pdReceiver.drug_id = d.drug_id
		WHERE
			d.drug_id IN (SELECT drug_id FROM pharmacy_drugs WHERE pharmacy_id = $1)
		AND
			d.drug_id IN (SELECT drug_id FROM pharmacy_drugs WHERE pharmacy_id = $2)
		AND
			pdSender.pharmacy_id = $1
		AND
			pdReceiver.pharmacy_id = $2
		AND
		%s
		%s
	`

	clc.Args = append(clc.Args, pharmacySenderID, pharmacyReceiverID)

	search := utils.BuildSearchQuery(srSearchColumn, clc)
	orderBy := utils.BuildSortQuery(nil, "", "d.drug_name DESC")
	filter := utils.BuildFilterQuery(map[string]string{}, clc, "d.deleted_at IS NULL")

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

	results := make([]entity.DrugWithPharmacyDrug, 0)
	for rows.Next() {
		var scan entity.DrugWithPharmacyDrug
		if err := rows.Scan(
			&scan.ID,
			&scan.Name,
			&scan.GenericName,
			&scan.Composition,
			&scan.Description,
			&scan.Classification,
			&scan.Form,
			&scan.UnitInPack,
			&scan.SellingUnit,
			&scan.Weight,
			&scan.Height,
			&scan.Length,
			&scan.Width,
			&scan.ImageURL,
			&scan.IsActive,
			&scan.CreatedAt,
			&scan.SenderPharmacyDrug.ID,
			&scan.SenderPharmacyDrug.DrugID,
			&scan.SenderPharmacyDrug.PharmacyID,
			&scan.SenderPharmacyDrug.CategoryID,
			&scan.SenderPharmacyDrug.Stock,
			&scan.SenderPharmacyDrug.Price,
			&scan.SenderPharmacyDrug.IsActive,
			&scan.SenderPharmacyDrug.CreatedAt,
			&scan.ReceiverPharmacyDrug.ID,
			&scan.ReceiverPharmacyDrug.DrugID,
			&scan.ReceiverPharmacyDrug.PharmacyID,
			&scan.ReceiverPharmacyDrug.CategoryID,
			&scan.ReceiverPharmacyDrug.Stock,
			&scan.ReceiverPharmacyDrug.Price,
			&scan.ReceiverPharmacyDrug.IsActive,
			&scan.ReceiverPharmacyDrug.CreatedAt,
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
