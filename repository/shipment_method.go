package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/libs/rajaongkir"
	"Alice-Seahat-Healthcare/seahat-be/utils"

	"github.com/sirupsen/logrus"
)

type ShipmentMethodRepository interface {
	SelectAllAvailByPharmacyID(ctx context.Context, pharmacyIDs []uint) ([]*entity.Pharmacy, error)
	GetDistanceKM(ctx context.Context, srcLoc string, destLoc string) (*float64, error)
	GetThirdPartyShipmentPrice(ctx context.Context, payload rajaongkir.CostPayload, etd uint) (float64, error)
	InsertManyPharmacyShipment(ctx context.Context, pharmacyID uint, shipments []string) error
	HardDeleteAllShipmentByPharmacy(ctx context.Context, pharmacyID uint) error
	SelectAllShipmentMethod(ctx context.Context) ([]*entity.ShipmentMethod, error)
	GetPharmacySMethodByShipmentIdAndPharmacyID(ctx context.Context, pharmacyID uint, shipmentId uint) (*entity.Pharmacy, error)
}

type shipmentMethodRepositoryImpl struct {
	db transaction.DBTransaction
	ro rajaongkir.RajaOngkir
}

func NewShipmentMethodRepository(db transaction.DBTransaction, ro rajaongkir.RajaOngkir) *shipmentMethodRepositoryImpl {
	return &shipmentMethodRepositoryImpl{
		db: db,
		ro: ro,
	}
}

func (r *shipmentMethodRepositoryImpl) SelectAllAvailByPharmacyID(ctx context.Context, pharmacyIDs []uint) ([]*entity.Pharmacy, error) {
	q := `
		SELECT
			p.pharmacy_id,
			p.pharmacy_name,
			ST_AsEWKT(p.pharmacy_location),
			p.subdistrict_id,
			p.created_at,
			sm.shipment_method_id,
			sm.shipment_method_name,
			sm.courier_name,
			sm.price,
			sm.duration,
			sm.created_at,
			sd.subdistrict_id, 
			sd.city_id, 
			sd.subdistrict_name, 
			sd.created_at 
		FROM pharmacies p
		LEFT JOIN pharmacy_shipments ps ON ps.pharmacy_id = p.pharmacy_id
		LEFT JOIN shipment_methods sm ON sm.shipment_method_id = ps.shipment_method_id
		LEFT JOIN subdistricts sd ON sd.subdistrict_id = p.subdistrict_id
		WHERE
			p.pharmacy_id = ANY($1::int[])
		AND
			ps.deleted_at IS NULL
		ORDER BY
			p.pharmacy_id ASC,
			sm.shipment_method_id ASC
	`

	param := new(strings.Builder)
	param.WriteString("{")

	idsLength := len(pharmacyIDs)

	for index, id := range pharmacyIDs {
		param.WriteString(fmt.Sprint(id))

		if index != idsLength-1 {
			param.WriteString(",")
		}
	}

	param.WriteString("}")

	rows, err := r.db.QueryContext(ctx, q, param.String())
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	results := make([]*entity.Pharmacy, 0)
	for rows.Next() {
		var p entity.Pharmacy
		var sm entity.ShipmentMethod

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Location,
			&p.SubdistrictID,
			&p.CreatedAt,
			&sm.ID,
			&sm.Name,
			&sm.CourierName,
			&sm.Price,
			&sm.Duration,
			&sm.CreatedAt,
			&p.Subdistrict.ID,
			&p.Subdistrict.CityID,
			&p.Subdistrict.Name,
			&p.Subdistrict.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		p.Longitude, p.Latitude, err = utils.Geo2LongLat(p.Location)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		pharmacyIndex := r.findPharmacyIndex(results, p)
		if pharmacyIndex == -1 {
			p.ShipmentMethods = []*entity.ShipmentMethod{&sm}
			results = append(results, &p)
			continue
		}

		results[pharmacyIndex].ShipmentMethods = append(results[pharmacyIndex].ShipmentMethods, &sm)
	}

	if err = rows.Err(); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return results, nil
}

func (r *shipmentMethodRepositoryImpl) GetPharmacySMethodByShipmentIdAndPharmacyID(ctx context.Context, pharmacyID uint, shipmentId uint) (*entity.Pharmacy, error) {
	q := `SELECT
			p.pharmacy_id,
			p.pharmacy_name,
			ST_AsEWKT(p.pharmacy_location),
			p.subdistrict_id,
			p.created_at,
			sm.shipment_method_id,
			sm.shipment_method_name,
			sm.courier_name,
			sm.price,
			sm.duration,
			sm.created_at,
			sd.subdistrict_id, 
			sd.city_id, 
			sd.subdistrict_name, 
			sd.created_at 
		FROM pharmacies p
		LEFT JOIN pharmacy_shipments ps ON ps.pharmacy_id = p.pharmacy_id
		LEFT JOIN shipment_methods sm ON sm.shipment_method_id = ps.shipment_method_id
		LEFT JOIN subdistricts sd ON sd.subdistrict_id = p.subdistrict_id
		WHERE
			p.pharmacy_id = $1
		AND 
			sm.shipment_method_id=$2
		AND
			ps.deleted_at IS NULL`

	var p entity.Pharmacy
	var sm entity.ShipmentMethod
	err := r.db.QueryRowContext(ctx, q, pharmacyID, shipmentId).Scan(
		&p.ID,
		&p.Name,
		&p.Location,
		&p.SubdistrictID,
		&p.CreatedAt,
		&sm.ID,
		&sm.Name,
		&sm.CourierName,
		&sm.Price,
		&sm.Duration,
		&sm.CreatedAt,
		&p.Subdistrict.ID,
		&p.Subdistrict.CityID,
		&p.Subdistrict.Name,
		&p.Subdistrict.CreatedAt,
	)
	p.ShipmentMethods = []*entity.ShipmentMethod{&sm}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.ErrResourceNotFound
		}
		logrus.Error(err)
		return nil, err
	}
	return &p, nil
}

func (r *shipmentMethodRepositoryImpl) findPharmacyIndex(pharmacies []*entity.Pharmacy, pharmacy entity.Pharmacy) int {
	for index, p := range pharmacies {
		if p.ID == pharmacy.ID {
			return index
		}
	}

	return -1
}

func (r *shipmentMethodRepositoryImpl) GetDistanceKM(ctx context.Context, srcLoc string, destLoc string) (*float64, error) {
	q := `
		SELECT ST_Distance($1::geography, $2::geography) AS distance
	`

	var metre float64
	err := r.db.QueryRowContext(ctx, q, srcLoc, destLoc).Scan(&metre)
	if err != nil {
		return nil, err
	}

	metre /= constant.MetreToKilometre
	return &metre, nil
}

func (r *shipmentMethodRepositoryImpl) GetThirdPartyShipmentPrice(ctx context.Context, payload rajaongkir.CostPayload, etd uint) (float64, error) {
	data, err := r.ro.Post(ctx, "/cost", payload)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	costRes := new(rajaongkir.CostResponse)
	if err := json.Unmarshal(data, costRes); err != nil {
		logrus.Error(err)
		return 0, err
	}

	couriers := costRes.RajaOngkir.Results
	if len(couriers) == 0 {
		return 0, nil
	}

	for _, service := range couriers[0].Costs {
		for _, cost := range service.Cost {
			if strings.Contains(cost.Etd, strconv.Itoa(int(etd))) {
				return cost.Value, nil
			}
		}
	}

	return 0, nil
}

func (r *shipmentMethodRepositoryImpl) InsertManyPharmacyShipment(ctx context.Context, pharmacyID uint, shipments []string) error {
	q := `
		INSERT INTO 
			pharmacy_shipments (pharmacy_id, shipment_method_id) 
		SELECT $1, sm.shipment_method_id
		FROM shipment_methods sm
		WHERE sm.courier_name = ANY($2::varchar[])
	`
	param := "{" + strings.Join(shipments, ",") + "}"
	if _, err := r.db.ExecContext(ctx, q, pharmacyID, param); err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (r *shipmentMethodRepositoryImpl) HardDeleteAllShipmentByPharmacy(ctx context.Context, pharmacyID uint) error {
	q := `
		DELETE FROM pharmacy_shipments
		WHERE pharmacy_id = $1
	`

	if _, err := r.db.ExecContext(ctx, q, pharmacyID); err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (r *shipmentMethodRepositoryImpl) SelectAllShipmentMethod(ctx context.Context) ([]*entity.ShipmentMethod, error) {
	q := `
		SELECT
			shipment_method_id, shipment_method_name, courier_name, price, duration, created_at
		FROM
			shipment_methods
		ORDER BY
			created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()

	results := make([]*entity.ShipmentMethod, 0)
	for rows.Next() {
		var scan entity.ShipmentMethod
		if err := rows.Scan(&scan.ID, &scan.Name, &scan.CourierName, &scan.Price, &scan.Duration, &scan.CreatedAt); err != nil {
			logrus.Error(err)
			return nil, err
		}

		results = append(results, &scan)
	}

	if err := rows.Err(); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return results, nil
}
