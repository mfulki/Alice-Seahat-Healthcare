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
	paymentColumnAlias = map[string]string{
		"payment_method": "p.payment_method",
		"payment_proof":  "p.payment_proof",
		"payment_number": "p.payment_number",
		"total_price":    "p.total_price",
	}
	paymentSearchColumn = []string{
		"payment_number",
	}
)

type PaymentRepository interface {
	InsertPayment(ctx context.Context, payment *entity.Payment) (*entity.Payment, error)
	UpdatePaymentProof(ctx context.Context, payment entity.Payment) (*entity.Payment, error)
	UserDeletePayment(ctx context.Context, payment entity.Payment) (*entity.Payment, error)
	AdminDeletePayment(ctx context.Context, payment entity.Payment) (*entity.Payment, error)
	GetAllPaymentToConfirm(ctx context.Context, clc *entity.Collection) ([]*entity.Payment, error)
	GetAllPaymentByUserId(ctx context.Context, userId uint) (map[uint]*entity.Payment, error)
	UpdatePaymentExpiredAt(ctx context.Context, paymentId uint, futureStatus string) error
}

type paymentRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewPaymentRepository(db transaction.DBTransaction) *paymentRepositoryImpl {
	return &paymentRepositoryImpl{
		db: db,
	}
}

func (r *paymentRepositoryImpl) InsertPayment(ctx context.Context, payment *entity.Payment) (*entity.Payment, error) {
	q := `insert into payments 
		(user_id,payment_method,payment_expired_at,full_user_address,total_price,payment_number)
		values
		($1,'manual transfer',Now()+interval '10 minute',$2,$3,'PAY-'||to_char(now(),'YYYY')||'-'||'NUM'||to_char(now(),'MM')||'T'||'-'||uuid_generate_v4()||'-'||to_char(now(),'DD'||'H'||'A'))
		returning payment_id,payment_number;
	`

	err := r.db.QueryRowContext(ctx, q, payment.UserId, payment.FullUserAddress, payment.TotalPrice).Scan(&payment.Id, &payment.Number)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return payment, nil
}

func (r *paymentRepositoryImpl) UpdatePaymentProof(ctx context.Context, payment entity.Payment) (*entity.Payment, error) {
	q := `UPDATE
			payments 
		SET
			updated_at=now(),
			payment_proof=$1
		WHERE
			payment_id=$2
		AND
			user_id=$3
		AND
			deleted_at is null
		Returning payment_proof, payment_method,full_user_address,total_price,payment_number
		`
	err := r.db.QueryRowContext(ctx, q, payment.Proof, payment.Id, payment.UserId).Scan(&payment.Proof, &payment.Method, &payment.FullUserAddress, &payment.TotalPrice, &payment.Number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrResourceNotFound
		}

		logrus.Error(err)
		return nil, err
	}
	return &payment, err
}
func (r *paymentRepositoryImpl) UpdatePaymentExpiredAt(ctx context.Context, paymentId uint, futureStatus string) error {
	var expiredAt string
	if futureStatus == constant.WaitingForPayment {
		expiredAt = `Now()+ interval '10 minutes', payment_proof=null`
	}
	if futureStatus == constant.PaymentConfirmed {
		expiredAt = `null`
	}
	q := fmt.Sprintf(`Update 
			payments 
		set 
			updated_at=now(),
			payment_expired_at=%s
		Where
			payment_id=$1
		And
			deleted_at is null
		`, expiredAt)
	_, err := r.db.ExecContext(ctx, q, paymentId)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (r *paymentRepositoryImpl) UserDeletePayment(ctx context.Context, payment entity.Payment) (*entity.Payment, error) {
	q := `	UPDATE
				payments 
			SET
				deleted_at=now(),
				updated_at = current_timestamp

			WHERE
				payment_proof is null
			AND 
				payment_id=$1
			AND
				user_id=$2
			AND
				deleted_at is null
			Returning payment_proof, payment_method,full_user_address,total_price,payment_number
			`
	err := r.db.QueryRowContext(ctx, q, payment.Id, payment.UserId).Scan(&payment.Proof, &payment.Method, &payment.FullUserAddress, &payment.TotalPrice, &payment.Number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrResourceNotFound
		}

		logrus.Error(err)
		return nil, err
	}
	return &payment, err
}

func (r *paymentRepositoryImpl) AdminDeletePayment(ctx context.Context, payment entity.Payment) (*entity.Payment, error) {
	q := `	UPDATE
				payments 
			SET
				deleted_at=now(),
				updated_at = current_timestamp

			WHERE
				payment_id=$1
			AND
				deleted_at is null
		
			Returning payment_proof, payment_method,full_user_address,total_price,payment_number
			`
	err := r.db.QueryRowContext(ctx, q, payment.Id).Scan(&payment.Proof, &payment.Method, &payment.FullUserAddress, &payment.TotalPrice, &payment.Number)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &payment, err
}

func (r *paymentRepositoryImpl) GetAllPaymentToConfirm(ctx context.Context, clc *entity.Collection) ([]*entity.Payment, error) {
	selectColumns := `
		distinct  
		p.payment_id,
		p.user_id, 
		u.user_name,
		p.payment_method,
		p.payment_proof,
		p.payment_number,
		p.total_price,
		p.updated_at,
		o.status 
	`
	advanceQuery := `
			payments p
		join orders o on p.payment_id = o.payment_id 
		join users u on p.user_id = u.user_id
		where 
			o.status = $1
		AND 
		%s
		%s
	`

	clc.Args = append(clc.Args, constant.WaitingForPaymentConfirmation)

	search := utils.BuildSearchQuery(paymentSearchColumn, clc)
	orderBy := utils.BuildSortQuery(paymentColumnAlias, clc.Sort, "p.updated_at desc")
	filter := utils.BuildFilterQuery(paymentColumnAlias, clc, "p.deleted_at is null")

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

	var scan []*entity.Payment
	for rows.Next() {
		s := &entity.Payment{}
		err := rows.Scan(
			&s.Id,
			&s.UserId,
			&s.UserName,
			&s.Method,
			&s.Proof,
			&s.Number,
			&s.TotalPrice,
			&s.UpdatedAt,
			&s.Status,
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

func (r *paymentRepositoryImpl) GetAllPaymentByUserId(ctx context.Context, userId uint) (map[uint]*entity.Payment, error) {
	q := `
	select 
	p.payment_id,
	p.user_id, 
	p.payment_method,
	p.payment_expired_at,
	p.payment_proof, 
	p.full_user_address, 
	p.total_price, 
	p.payment_number,
	p.deleted_at,
	p.created_at,
	o.status,
	o.order_id,
	o.pharmacy_id ,
	ph.pharmacy_name,
	o.order_number ,
	o.shipment_method_name,
	o.shipment_price,
	o.finished_at,
	o.total_price, 
	o.created_at,
	pd.pharmacy_drug_id,
	pd.created_at,
	d.drug_id ,
	d.drug_name ,
	d.image_url ,
	od.order_detail_id , 
	od.price ,
	od.quantity 
	from payments p
	join orders o on o.payment_id = p.payment_id 
	join pharmacies ph on ph.pharmacy_id = o.pharmacy_id
	join order_details od on od.order_id = o.order_id
	join pharmacy_drugs pd on pd.pharmacy_drug_id = od.pharmacy_drug_id 
	join drugs d on d.drug_id = pd.drug_id 
	where  p.user_id =$1 
	AND
	od.deleted_at is null
	AND
	d.deleted_at is null
	AND
	pd.deleted_at is null
	order by p.payment_id desc
	`
	rows, err := r.db.QueryContext(ctx, q, userId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	oMap := map[uint]*entity.Order{}
	oDMap := map[uint][]*entity.OrderDetail{}
	pMap := map[uint]*entity.Payment{}
	oMaps := map[uint][]*entity.Order{}
	for rows.Next() {
		p := &entity.Payment{}
		ph := &entity.Pharmacy{}
		o := &entity.Order{}
		oD := &entity.OrderDetail{}
		err := rows.Scan(
			&p.Id,
			&p.UserId,
			&p.Method,
			&p.ExpiredAt,
			&p.Proof,
			&p.FullUserAddress,
			&p.TotalPrice,
			&p.Number,
			&p.DeletedAt,
			&p.CreatedAt,
			&o.Status,
			&o.Id,
			&ph.ID,
			&ph.Name,
			&o.OrderNumber,
			&o.ShipmentMethod.Name,
			&o.ShipmentMethod.Price,
			&o.FinishedAt,
			&o.TotalPrice,
			&o.CreatedAt,
			&oD.PharmacyDrug.ID,
			&oD.PharmacyDrug.CreatedAt,
			&oD.PharmacyDrug.Drug.ID,
			&oD.PharmacyDrug.Drug.Name,
			&oD.PharmacyDrug.Drug.ImageURL,
			&oD.Id,
			&oD.Price,
			&oD.Quantity,
		)

		oDMap[o.Id] = append(oDMap[o.Id], oD)
		o.Detail = oDMap[o.Id]

		o.Payment = p
		o.Pharmacy = ph
		o.PharmacyId = ph.ID
		oMap[o.Id] = o
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		pMap[p.Id] = p
	}

	for _, order := range oMap {
		pId := order.Payment.Id
		order.Payment = nil
		oMaps[pId] = append(oMaps[pId], order)
	}

	for pId, orders := range oMaps {
		pMap[pId].Orders = orders
	}

	if err = rows.Err(); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return pMap, nil

}
