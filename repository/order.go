package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/utils"

	"github.com/sirupsen/logrus"
)

type OrderRepository interface {
	InsertOrder(ctx context.Context, orders []entity.Order) ([]entity.Order, error)
	UpdateOrderStatusByPaymentId(ctx context.Context, payment entity.Payment, futureStatus string, recentStatus string) ([]*entity.Order, error)
	UpdateOrderStatusByOrderId(ctx context.Context, order entity.Order, updateStatus string, userId uint) (*entity.Order, error)
	SelectOrdersByOrderId(ctx context.Context, order entity.Order) (*entity.Order, error)
	PMUpdateOrderStatusByOrderId(ctx context.Context, order entity.Order, updateStatus string, pMId uint) (*entity.Order, error)
	GetAllOrderByPharmacyManagerId(ctx context.Context, pharmacyManagerId uint) ([]*entity.Order, error)
	SelecOrderStatusByOrderId(ctx context.Context, orderId uint) (*string, error)
}

type orderRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewOrderRepository(db transaction.DBTransaction) *orderRepositoryImpl {
	return &orderRepositoryImpl{
		db: db,
	}
}

func (r *orderRepositoryImpl) InsertOrder(ctx context.Context, orders []entity.Order) ([]entity.Order, error) {
	var s strings.Builder
	args := []any{}
	s.WriteString(`insert into orders 
		(payment_id, pharmacy_id, order_number, total_price, status, shipment_price, shipment_method_name) 
		values `)
	for num, order := range orders {
		if num > 0 {
			s.WriteString(`,`)
		}
		orders[num].Status = constant.WaitingForPayment

		args = append(args, order.Payment.Id, order.PharmacyId, order.TotalPrice, orders[num].Status, order.ShipmentMethod.Price, order.ShipmentMethod.Name)
		Parameters := 6
		s.WriteString(`(`)
		for i := 1 + (Parameters * num); i <= (num+1)*Parameters; i++ {
			s.WriteString(fmt.Sprintf(`$%s`, strconv.Itoa(i)))
			if i != (num+1)*Parameters {
				s.WriteString(`,`)
			}
			if i == 2+(Parameters*num) {
				s.WriteString(`uuid_generate_v4()`)
				s.WriteString(`,`)
				continue
			}
		}
		s.WriteString(`)`)
	}
	s.WriteString(`Returning order_id,order_number`)
	q := s.String()

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()
	numeric := 0
	for rows.Next() {
		order := &orders[numeric]
		err := rows.Scan(
			&order.Id,
			&order.OrderNumber,
		)
		numeric++

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return orders, nil

}

func (r *orderRepositoryImpl) UpdateOrderStatusByPaymentId(ctx context.Context, payment entity.Payment, futureStatus string, recentStatus string) ([]*entity.Order, error) {
	q := `Update orders o
		SET
			status=$1,
			updated_at=now()
		FROM
			payments p
		WHERE
			o.payment_id=$2
		AND 
			o.status=$3
		AND
			p.payment_id=o.payment_id
		Returning o.order_id,o.status,o.pharmacy_id,o.order_number,o.total_price,o.shipment_price,o.shipment_method_name
		`

	rows, err := r.db.QueryContext(ctx, q, futureStatus, payment.Id, recentStatus)
	if err != nil {
		logrus.Error(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ResourceNotFound
		}
		return nil, err
	}

	defer rows.Close()
	var orders []*entity.Order
	rowsAffected := 0
	for rows.Next() {
		rowsAffected++
		order := &entity.Order{}
		order.Payment = &payment
		err := rows.Scan(
			&order.Id,
			&order.Status,
			&order.PharmacyId,
			&order.OrderNumber,
			&order.TotalPrice,
			&order.ShipmentMethod.Price,
			&order.ShipmentMethod.Name,
		)

		if err != nil {
			logrus.Error(err)
			if errors.Is(err, apperror.ErrResourceNotFound) {
				return nil, apperror.ResourceNotFound
			}
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		logrus.Error(err)

		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.ResourceNotFound
		}
		return nil, err
	}

	if rowsAffected == 0 {
		logrus.Error()
		return nil, apperror.NoValidOrderPayment
	}
	return orders, nil

}

func (r *orderRepositoryImpl) UpdateOrderStatusByOrderId(ctx context.Context, order entity.Order, updateStatus string, userId uint) (*entity.Order, error) {
	q := `Update orders o 
		SET
			status=$1,
			finished_at=now(),
			updated_at=now()
		FROM 
			payments p
		WHERE
			order_id=$2
		AND 
			status=$3
		AND
			p.user_id=$4
		AND
			p.payment_id = o.payment_id 
		Returning o.order_id,o.status,o.pharmacy_id,o.order_number,o.total_price,o.shipment_price,o.shipment_method_name
		`
	err := r.db.QueryRowContext(ctx, q, updateStatus, order.Id, order.Status, userId).Scan(
		&order.Id,
		&order.Status,
		&order.PharmacyId,
		&order.OrderNumber,
		&order.TotalPrice,
		&order.ShipmentMethod.Price,
		&order.ShipmentMethod.Name,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrResourceNotFound
		}

		logrus.Error(err)
		return nil, err
	}
	return &order, nil

}
func (r *orderRepositoryImpl) PMUpdateOrderStatusByOrderId(ctx context.Context, order entity.Order, updateStatus string, pMId uint) (*entity.Order, error) {
	futureStatus := ``
	if updateStatus == constant.Sent {
		futureStatus = `,finished_at=Now()+interval '7 days'`
	}
	if updateStatus == constant.Cancelled {
		futureStatus = `,finished_at=Now()`
	}
	q := fmt.Sprintf(`Update orders o 
		SET
			status=$1,
			updated_at=now() %s
		FROM 
			pharmacies p
		WHERE
			order_id=$2
		AND 
			status=$3
		AND
			o.pharmacy_id=p.pharmacy_id
		AND
			p.pharmacy_manager_id=$4
		AND
			o.deleted_at is null
		AND
			p.deleted_at is null
		Returning o.order_id,o.status,o.pharmacy_id,o.order_number,o.total_price,o.shipment_price,o.shipment_method_name
		`, futureStatus)
	err := r.db.QueryRowContext(ctx, q, updateStatus, order.Id, order.Status, pMId).Scan(
		&order.Id,
		&order.Status,
		&order.PharmacyId,
		&order.OrderNumber,
		&order.TotalPrice,
		&order.ShipmentMethod.Price,
		&order.ShipmentMethod.Name,
	)
	if err != nil {
		logrus.Error(err)
		if err == sql.ErrNoRows {
			return nil, apperror.ResourceNotFound
		}
		return nil, err
	}
	return &order, nil

}
func (r *orderRepositoryImpl) SelecOrderStatusByOrderId(ctx context.Context, orderId uint) (*string, error) {
	q := `
		Select status from orders where order_id=$1 and deleted_at is null
	`
	var s string
	err := r.db.QueryRowContext(ctx, q, orderId).Scan(&s)
	if err != nil {
		logrus.Error(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrResourceNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *orderRepositoryImpl) SelectOrdersByOrderId(ctx context.Context, order entity.Order) (*entity.Order, error) {
	q := `
	  	select 
		o.pharmacy_id,
		o.status,
		od.order_detail_id,
		od.order_id,
		od.pharmacy_drug_id,
		od.quantity,
		pd.stock,
		pd.drug_id,
	  	pd.pharmacy_id,
		p.pharmacy_manager_id,
		ST_AsEWKT(p.pharmacy_location)  
		from orders o
		Join order_details od ON od.order_id=o.order_id 
	  	JOIN pharmacy_drugs pd ON pd.pharmacy_drug_id =od.pharmacy_drug_id 
		JOIN pharmacies p ON o.pharmacy_id =p.pharmacy_id
	  	where o.order_id=$1 
		AND o.status=$2 
		AND o.deleted_at is null
		AND od.deleted_at is null
		AND pd.deleted_at is null
		AND p.deleted_at is null
		 for update
		`
	rows, err := r.db.QueryContext(ctx, q, order.Id, order.Status)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer rows.Close()
	numeric := 0
	orderDetails := []*entity.OrderDetail{}
	for rows.Next() {
		orderDetail := &entity.OrderDetail{}
		err := rows.Scan(
			&order.PharmacyId,
			&order.Status,
			&orderDetail.Id,
			&orderDetail.OrderId,
			&orderDetail.PharmacyDrugId,
			&orderDetail.Quantity,
			&orderDetail.PharmacyDrug.Stock,
			&orderDetail.PharmacyDrug.DrugID,
			&orderDetail.PharmacyDrug.PharmacyID,
			&orderDetail.PharmacyDrug.Pharmacy.ManagerID,
			&orderDetail.PharmacyDrug.Pharmacy.Location,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		numeric++

		orderDetail.PharmacyDrug.Pharmacy.Longitude, orderDetail.PharmacyDrug.Pharmacy.Latitude, err = utils.Geo2LongLat(orderDetail.PharmacyDrug.Pharmacy.Location)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		orderDetails = append(orderDetails, orderDetail)

	}
	order.Detail = orderDetails

	if err = rows.Err(); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &order, nil

}

func (r *orderRepositoryImpl) GetAllOrderByPharmacyManagerId(ctx context.Context, pharmacyManagerId uint) ([]*entity.Order, error) {
	q := `
	select 
	p.payment_id,
	p.user_id, 
	u.user_name,
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
	join orders o on o.payment_id =p.payment_id 
	join users u on p.user_id=u.user_id
	join pharmacies ph on ph.pharmacy_id=o.pharmacy_id
	join order_details od on od.order_id  =o.order_id
	join pharmacy_drugs pd on pd.pharmacy_drug_id =od.pharmacy_drug_id 
	join drugs d on d.drug_id =pd.drug_id 
	where  
		ph.pharmacy_manager_id=$1
	AND
	od.deleted_at is null
	AND
	d.deleted_at is null
	AND
	pd.deleted_at is null
	AND
	p.deleted_at is null
	AND
	ph.deleted_at is null
	AND
	d.deleted_at is null
	order by p.payment_id desc
	`
	rows, err := r.db.QueryContext(ctx, q, pharmacyManagerId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	oMap := map[uint]*entity.Order{}
	oDMap := map[uint][]*entity.OrderDetail{}
	orders := []*entity.Order{}
	for rows.Next() {
		p := &entity.Payment{}
		ph := &entity.Pharmacy{}
		o := &entity.Order{}
		oD := &entity.OrderDetail{}
		err := rows.Scan(
			&p.Id,
			&p.UserId,
			&p.UserName,
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
	}

	keys := make([]uint, 0, len(oMap))
	for k := range oMap {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	for _, k := range keys {
		orders = append(orders, oMap[k])
	}

	if err = rows.Err(); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return orders, nil

}
