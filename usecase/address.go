package usecase

import (
	"context"
	"errors"
	"math"
	"sync"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/libs/rajaongkir"
	"Alice-Seahat-Healthcare/seahat-be/repository"
	"Alice-Seahat-Healthcare/seahat-be/utils"
)

type AddressUsecase interface {
	GetAllProvinces(ctx context.Context) ([]entity.Province, error)
	GetAllCitiesWithProvinceQuery(ctx context.Context, provinceID uint) ([]entity.City, error)
	GetAllSubdistrictsWithCityQuery(ctx context.Context, cityID uint) ([]entity.Subdistrict, error)
	GetAllUserAddress(ctx context.Context, isActive *bool) ([]entity.Address, error)
	AddAddress(ctx context.Context, addr entity.Address) (*entity.Address, error)
	GetAddressByID(ctx context.Context, id uint) (*entity.Address, error)
	UpdateAddressByID(ctx context.Context, addr entity.Address) error
	DeleteAddressByID(ctx context.Context, id uint) error
	GetShipmentPriceByPharmaciesID(ctx context.Context, addressID uint, pharmacyShipmentMap map[uint]uint) ([]*entity.Pharmacy, error)
}

type addressUsecaseImpl struct {
	addressRepository        repository.AddressRepository
	shipmentMethodRepository repository.ShipmentMethodRepository
	transactor               transaction.Transactor
}

func NewAddressUsecase(
	addressRepository repository.AddressRepository,
	shipmentMethodRepository repository.ShipmentMethodRepository,
	transactor transaction.Transactor,
) *addressUsecaseImpl {
	return &addressUsecaseImpl{
		addressRepository:        addressRepository,
		shipmentMethodRepository: shipmentMethodRepository,
		transactor:               transactor,
	}
}

func (u *addressUsecaseImpl) GetAllProvinces(ctx context.Context) ([]entity.Province, error) {
	return u.addressRepository.SelectAllProvince(ctx)
}

func (u *addressUsecaseImpl) GetAllCitiesWithProvinceQuery(ctx context.Context, provinceID uint) ([]entity.City, error) {
	return u.addressRepository.SelectAllCitiesWithProvinceQuery(ctx, provinceID)
}

func (u *addressUsecaseImpl) GetAllSubdistrictsWithCityQuery(ctx context.Context, cityID uint) ([]entity.Subdistrict, error) {
	return u.addressRepository.SelectAllSubdistrictsWithCityQuery(ctx, cityID)
}

func (u *addressUsecaseImpl) GetAllUserAddress(ctx context.Context, isActive *bool) ([]entity.Address, error) {
	user, ok := utils.CtxGetUser(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}

	addrs, err := u.addressRepository.GetAll(ctx, user.ID, isActive)
	if err != nil {
		return nil, err
	}

	return addrs, nil
}

func (u *addressUsecaseImpl) AddAddress(ctx context.Context, addr entity.Address) (*entity.Address, error) {
	user, ok := utils.CtxGetUser(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}

	_, err := u.addressRepository.GetSubdistrictByID(ctx, addr.SubdistrictID)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.SubdistrictNotAvail
		}

		return nil, err
	}

	mainAddr, err := u.addressRepository.GetMainAddress(ctx, user.ID)
	if err != nil && !errors.Is(err, apperror.ErrResourceNotFound) {
		return nil, err
	}

	if mainAddr == nil && !addr.IsMain {
		return nil, apperror.MustMainAddress
	}

	if mainAddr != nil && addr.IsMain {
		mainAddr.IsMain = false

		err := u.addressRepository.UpdateByID(ctx, *mainAddr)
		if err != nil {
			return nil, err
		}
	}

	addr.UserID = user.ID
	addedAddr, err := u.addressRepository.InsertOne(ctx, addr)
	if err != nil {
		return nil, err
	}

	return addedAddr, nil
}

func (u *addressUsecaseImpl) GetAddressByID(ctx context.Context, id uint) (*entity.Address, error) {
	user, ok := utils.CtxGetUser(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}

	addr, err := u.addressRepository.GetByID(ctx, id, user.ID)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.ResourceNotFound
		}

		return nil, err
	}

	return addr, nil
}

func (u *addressUsecaseImpl) UpdateAddressByID(ctx context.Context, addr entity.Address) error {
	user, ok := utils.CtxGetUser(ctx)
	if !ok {
		return apperror.ErrInternalServer
	}

	_, err := u.addressRepository.GetSubdistrictByID(ctx, addr.SubdistrictID)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return apperror.SubdistrictNotAvail
		}

		return err
	}

	mainAddr, err := u.addressRepository.GetMainAddress(ctx, user.ID)
	if err != nil {
		return err
	}

	if !addr.IsMain && mainAddr.ID == addr.ID {
		return apperror.MustMainAddress
	}

	if addr.IsMain {
		mainAddr.IsMain = false

		err := u.addressRepository.UpdateByID(ctx, *mainAddr)
		if err != nil {
			return err
		}
	}

	addr.UserID = user.ID

	if err := u.addressRepository.UpdateByID(ctx, addr); err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return apperror.ResourceNotFound
		}

		return err
	}

	return nil
}

func (u *addressUsecaseImpl) DeleteAddressByID(ctx context.Context, id uint) error {
	user, ok := utils.CtxGetUser(ctx)
	if !ok {
		return apperror.ErrInternalServer
	}

	mainAddr, err := u.addressRepository.GetMainAddress(ctx, user.ID)
	if err != nil {
		return err
	}

	if mainAddr.ID == id {
		return apperror.CantDeleteMainAddress
	}

	if err := u.addressRepository.DeleteByID(ctx, id, user.ID); err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return apperror.ResourceNotFound
		}

		return err
	}

	return nil
}

func (u *addressUsecaseImpl) GetShipmentPriceByPharmaciesID(ctx context.Context, addressID uint, pharmacyShipmentMap map[uint]uint) ([]*entity.Pharmacy, error) {
	pharmaciesID := make([]uint, 0)
	for key := range pharmacyShipmentMap {
		pharmaciesID = append(pharmaciesID, key)
	}

	if len(pharmaciesID) == 0 {
		return nil, nil
	}

	user, ok := utils.CtxGetUser(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}

	addr, err := u.addressRepository.GetByID(ctx, addressID, user.ID)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.AddressNotExist
		}

		return nil, err
	}

	pharmacies, err := u.shipmentMethodRepository.SelectAllAvailByPharmacyID(ctx, pharmaciesID)
	if err != nil {
		return nil, err
	}

	userCity, err := u.addressRepository.GetCityBySubdistrictID(ctx, addr.SubdistrictID)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for _, p := range pharmacies {
		distance, err := u.GetDistanceKM(ctx, addr.Location, p.Location)
		if err != nil {
			return nil, err
		}

		for _, sm := range p.ShipmentMethods {
			if sm.Price != nil {
				*sm.Price *= distance
				if *sm.Price > constant.MaxShipmentPrice {
					*sm.Price = 0
				}

				continue
			}

			wg.Add(1)
			go func(wg *sync.WaitGroup, p *entity.Pharmacy, sm *entity.ShipmentMethod) {
				defer wg.Done()

				price, _ := u.shipmentMethodRepository.GetThirdPartyShipmentPrice(ctx, rajaongkir.CostPayload{
					Origin:      userCity.ID,
					Destination: p.Subdistrict.CityID,
					Weight:      pharmacyShipmentMap[p.ID],
					Courier:     sm.CourierName,
				}, 1)

				uintPrice := uint(price)
				sm.Price = &uintPrice
			}(&wg, p, sm)
		}
	}

	wg.Wait()

	return pharmacies, nil
}

func (u *addressUsecaseImpl) GetDistanceKM(ctx context.Context, srcLoc, destLoc string) (uint, error) {
	d, err := u.shipmentMethodRepository.GetDistanceKM(ctx, srcLoc, destLoc)
	if err != nil {
		return 0, err
	}

	dRounded := math.Round(*d)
	distanceFloat := math.Max(dRounded, constant.MinShipmentDistance)
	distance := uint(distanceFloat)

	return distance, nil
}
