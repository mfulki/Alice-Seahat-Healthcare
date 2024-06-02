package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/utils"

	"github.com/sirupsen/logrus"
)

type OrderDetailRepository interface {
	InsertOrderDetail(ctx context.Context, cart []*entity.CartItem, orderId uint) ([]*entity.OrderDetail, error)
	SelectOrderDetailByOrderId(ctx context.Context, orderId uint) ([]*entity.OrderDetail, error)
}

type orderDetailRepositoryImpl struct {
	db transaction.DBTransaction
}

func NewOrderDetailRepository(db transaction.DBTransaction) *orderDetailRepositoryImpl {
	return &orderDetailRepositoryImpl{
		db: db,
	}
}

func (r *orderDetailRepositoryImpl) InsertOrderDetail(ctx context.Context, cart []*entity.CartItem, orderId uint) ([]*entity.OrderDetail, error) {
	var s strings.Builder
	var orderDetails []*entity.OrderDetail
	args := []any{}
	s.WriteString(`insert into order_details (order_id, pharmacy_drug_id, quantity, price) values `)
	for num, cartItem := range cart {
		if num > 0 {
			s.WriteString(`,`)
		}
		args = append(args, orderId, cartItem.PharmacyDrugID, cartItem.Quantity, cartItem.Price)
		Parameters := 4
		s.WriteString(`(`)
		for i := 1 + (Parameters * num); i <= (num+1)*Parameters; i++ {
			s.WriteString(fmt.Sprintf(`$%s`, strconv.Itoa(i)))
			if i != (num+1)*Parameters {
				s.WriteString(`,`)
			}
		}
		s.WriteString(`)`)
	}
	s.WriteString(` returning order_detail_id,order_id, pharmacy_drug_id, quantity, price`)
	q := s.String()
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()
	numeric := 0
	for rows.Next() {
		orderDetail := &entity.OrderDetail{}
		err := rows.Scan(
			&orderDetail.Id,
			&orderDetail.OrderId,
			&orderDetail.PharmacyDrugId,
			&orderDetail.Quantity,
			&orderDetail.Price,
		)
		numeric++
		orderDetails = append(orderDetails, orderDetail)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return orderDetails, nil

}

func (r *orderDetailRepositoryImpl) SelectOrderDetailByOrderId(ctx context.Context, orderId uint) ([]*entity.OrderDetail, error) {
	q := `
	  	select 
		od.order_detail_id,
		od.order_id,
		od.pharmacy_drug_id,
		od.quantity,
		pd.stock,
		pd.drug_id,
	  	pd.pharmacy_id,
		ST_AsEWKT(p.pharmacy_location)  
		from order_details od 
	  	JOIN pharmacy_drugs pd ON pd.pharmacy_drug_id =od.pharmacy_drug_id 
		JOIN pharmacies p ON pd.pharmacy_id =p.pharmacy_id
	  	where order_id=$1 for update
		`
	rows, err := r.db.QueryContext(ctx, q, orderId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer rows.Close()
	orderDetails := []*entity.OrderDetail{}
	numeric := 0
	for rows.Next() {
		orderDetail := &entity.OrderDetail{}
		err := rows.Scan(
			&orderDetail.Id,
			&orderDetail.OrderId,
			&orderDetail.PharmacyDrugId,
			&orderDetail.Quantity,
			&orderDetail.PharmacyDrug.Stock,
			&orderDetail.PharmacyDrug.DrugID,
			&orderDetail.PharmacyDrug.PharmacyID,
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
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return orderDetails, nil

}
