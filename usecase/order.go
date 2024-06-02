package usecase

import (
	"context"
	"errors"
	"math"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/libs/rajaongkir"
	"Alice-Seahat-Healthcare/seahat-be/repository"
	"Alice-Seahat-Healthcare/seahat-be/utils"

	"github.com/sirupsen/logrus"
)

type OrderUsecase interface {
	CreateOrder(ctx context.Context, orders []entity.Order) ([]entity.Order, error)
	UpdateConfirmOrder(ctx context.Context, body entity.Order) (*entity.Order, error)
	OrderProceed(ctx context.Context, order entity.Order) (*entity.Order, error)
	GetAllOrderByPharmacyManagerId(ctx context.Context) ([]*entity.Order, error)
	OrderSent(ctx context.Context, order entity.Order) error
	OrderCancelByPM(ctx context.Context, order entity.Order) error
}

type orderUsecaseImpl struct {
	orderRepository            repository.OrderRepository
	orderDetailRepository      repository.OrderDetailRepository
	paymentrepository          repository.PaymentRepository
	pharmacyDrugrepository     repository.PharmacyDrugRepository
	cartItemrepository         repository.CartItemRepository
	transactor                 transaction.Transactor
	stockJournalRepository     repository.StockJournalRepository
	stockRequestRepository     repository.StockRequestRepository
	stockRequestDrugRepository repository.StockRequestDrugRepository
	shipmentMethodRepository   repository.ShipmentMethodRepository
	addressRepository          repository.AddressRepository
}

func NewOrderUsecase(
	orderRepository repository.OrderRepository,
	orderDetailRepository repository.OrderDetailRepository,
	transactor transaction.Transactor,
	paymentrepository repository.PaymentRepository,
	pharmacyDrugrepository repository.PharmacyDrugRepository,
	cartItemrepository repository.CartItemRepository,
	stockJournalRepository repository.StockJournalRepository,
	stockRequestRepository repository.StockRequestRepository,
	stockRequestDrugRepository repository.StockRequestDrugRepository,
	shipmentMethodRepository repository.ShipmentMethodRepository,
	addressRepository repository.AddressRepository,

) *orderUsecaseImpl {
	return &orderUsecaseImpl{
		orderRepository:            orderRepository,
		transactor:                 transactor,
		paymentrepository:          paymentrepository,
		cartItemrepository:         cartItemrepository,
		pharmacyDrugrepository:     pharmacyDrugrepository,
		orderDetailRepository:      orderDetailRepository,
		stockJournalRepository:     stockJournalRepository,
		stockRequestRepository:     stockRequestRepository,
		stockRequestDrugRepository: stockRequestDrugRepository,
		shipmentMethodRepository:   shipmentMethodRepository,
		addressRepository:          addressRepository,
	}
}

func (u *orderUsecaseImpl) CreatePayment(ctx context.Context, payment *entity.Payment) (*entity.Payment, error) {
	payment, err := u.paymentrepository.InsertPayment(ctx, payment)
	if err != nil {
		return nil, err
	}
	return payment, err

}

func (u *orderUsecaseImpl) CreateOrder(ctx context.Context, orders []entity.Order) ([]entity.Order, error) {
	userCtx, ok := utils.CtxGetUser(ctx)

	if !ok {
		return nil, apperror.ErrInternalServer
	}
	orders[0].Payment.UserId = userCtx.ID
	orderTransaction, err := u.transactor.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		orders, err := u.createOrderTransaction(ctx, orders, userCtx.ID)
		if err != nil {
			return nil, err
		}
		return orders, err
	})
	if err != nil {
		return nil, err
	}

	orders = orderTransaction.([]entity.Order)

	return orders, nil
}
func (u *orderUsecaseImpl) createOrderTransaction(ctx context.Context, orders []entity.Order, userId uint) ([]entity.Order, error) {
	validOrders := make([]entity.Order, 0)
	address, err := u.addressRepository.GetByID(ctx, orders[0].Payment.Address.ID, userId)
	orders[0].Payment.FullUserAddress = address.Address
	if err != nil {
		return nil, err
	}

	for index, order := range orders {
		pharmacy, err := u.shipmentMethodRepository.GetPharmacySMethodByShipmentIdAndPharmacyID(ctx, order.PharmacyId, order.ShipmentMethod.ID)
		if err != nil {
			return nil, apperror.InvalidShipmentMethods
		}
		if pharmacy == nil {
			return nil, apperror.InvalidShipmentMethods
		}

		cart, err := u.cartItemrepository.LockRow(ctx, order.Cart, userId, order.PharmacyId)
		if err != nil {
			return nil, err
		}

		orders[index].Cart = cart
		var weight uint
		for _, oneCart := range cart {
			weight = weight + oneCart.PharmacyDrug.Drug.Weight*oneCart.Quantity
			orders[index].TotalPrice = orders[index].TotalPrice + int(oneCart.TotalPrice)

		}

		orderLen := len(orders[index].Cart)
		var shipmentPrice uint
		if orderLen != 0 {
			if order.ShipmentMethod.ID > 2 {
				payload := rajaongkir.CostPayload{Origin: pharmacy.Subdistrict.CityID, Destination: address.CityID, Weight: weight, Courier: pharmacy.ShipmentMethods[0].CourierName}
				price, err := u.shipmentMethodRepository.GetThirdPartyShipmentPrice(ctx, payload, constant.EstimatedDeliveryTime)
				if err != nil {
					return nil, err
				}
				if price == 0 {
					return nil, apperror.InvalidShipmentMethods
				}
				shipmentPrice = uint(price)

			} else {
				distance, err := u.getDistanceKM(ctx, pharmacy.Location, address.Location)
				if err != nil {
					return nil, err
				}
				shipmentPrice = distance * *pharmacy.ShipmentMethods[0].Price
			}

			orders[index].ShipmentMethod.CourierName = pharmacy.ShipmentMethods[0].CourierName
			orders[index].ShipmentMethod.Name = pharmacy.ShipmentMethods[0].Name
			orders[index].ShipmentMethod.Price = &shipmentPrice
			orders[index].TotalPrice = orders[index].TotalPrice + int(shipmentPrice)
			validOrders = append(validOrders, orders[index])
		}

	}

	orders = validOrders

	if len(orders) == 0 {
		return nil, apperror.NoValidCartOrder
	}

	for _, order := range orders {
		orders[0].Payment.TotalPrice = orders[0].Payment.TotalPrice + order.TotalPrice
	}

	_, err = u.CreatePayment(ctx, orders[0].Payment)
	if err != nil {
		return nil, err
	}
	orders, err = u.orderRepository.InsertOrder(ctx, orders)
	if err != nil {
		return nil, err
	}
	for index, order := range orders {
		order.Detail, err = u.orderDetailRepository.InsertOrderDetail(ctx, orders[index].Cart, order.Id)
		if err != nil {
			return nil, err
		}
		orders[index].Detail = order.Detail
		cartItemIds := make([]uint, 0)
		for _, cartItem := range order.Cart {
			cartItemIds = append(cartItemIds, cartItem.ID)
		}

		err = u.cartItemrepository.DeleteManyByID(ctx, cartItemIds, userId)
		if err != nil {
			return nil, err
		}
	}
	return orders, nil
}
func (u *orderUsecaseImpl) getDistanceKM(ctx context.Context, srcLoc, destLoc string) (uint, error) {
	d, err := u.shipmentMethodRepository.GetDistanceKM(ctx, srcLoc, destLoc)
	if err != nil {
		return 0, err
	}

	dRounded := math.Round(*d)
	distanceFloat := math.Max(dRounded, constant.MinShipmentDistance)
	distance := uint(distanceFloat)

	return distance, nil
}
func (u *orderUsecaseImpl) UpdateConfirmOrder(ctx context.Context, body entity.Order) (*entity.Order, error) {
	userCtx, ok := utils.CtxGetUser(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}
	userId := userCtx.ID
	body.Status = constant.Sent
	status := constant.OrderConfirmed
	orders, err := u.orderRepository.UpdateOrderStatusByOrderId(ctx, body, status, userId)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.ResourceNotFound
		}

		return nil, err
	}
	return orders, nil
}

func (u *orderUsecaseImpl) OrderProceed(ctx context.Context, order entity.Order) (*entity.Order, error) {
	managerCtx, ok := utils.CtxGetManager(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}
	managerId := managerCtx.ID
	stockRequestDrug := make(map[uint][]*entity.StockRequestDrug)
	soldStock := make(map[uint][]*entity.StockRequestDrug)
	stockJournals := make([]entity.StockJurnal, 0)
	order.Status = constant.PaymentConfirmed
	orderTx, err := u.transactor.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		orderTx, err := u.orderProceedTx(ctx, &order, managerId, stockRequestDrug, soldStock, stockJournals)
		if err != nil {
			return nil, err
		}
		return orderTx, nil
	})
	if err != nil {
		return nil, err
	}
	orderRes := orderTx.(*entity.Order)
	return orderRes, nil
}

func (u *orderUsecaseImpl) orderProceedTx(ctx context.Context, order *entity.Order, managerId uint, stockRequestDrug map[uint][]*entity.StockRequestDrug, soldStock map[uint][]*entity.StockRequestDrug, stockJournals []entity.StockJurnal) (*entity.Order, error) {
	order, err := u.orderRepository.SelectOrdersByOrderId(ctx, *order)
	if err != nil {
		return nil, err
	}
	if len(order.Detail) == 0 {
		logrus.Error(apperror.NoValidPayment)
		return nil, apperror.NoValidPayment
	}

	for _, orderDetail := range order.Detail {
		if managerId != orderDetail.PharmacyDrug.Pharmacy.ManagerID || order.PharmacyId != orderDetail.PharmacyDrug.PharmacyID {
			logrus.Error(apperror.Unauthorized)
			return nil, apperror.Unauthorized
		}

		if orderDetail.Quantity <= orderDetail.PharmacyDrug.Stock {
			description := "sold"
			stockJournals = append(stockJournals, entity.StockJurnal{DrugId: orderDetail.PharmacyDrug.DrugID, PharmacyId: order.PharmacyId, Quantity: int(orderDetail.Quantity) * -1, Description: description})
			soldStock[order.PharmacyId] = append(stockRequestDrug[order.PharmacyId], &entity.StockRequestDrug{DrugId: orderDetail.PharmacyDrug.DrugID, Quantity: int(orderDetail.Quantity)})

			if err != nil {
				return nil, err
			}
			continue
		}

		pharmacyDrug, err := u.pharmacyDrugrepository.SelectOneNearestByPharmacyManagerDrugId(ctx, orderDetail.PharmacyDrug, orderDetail.Quantity)
		if err != nil {
			return nil, err
		}

		stockJournals = append(stockJournals, entity.StockJurnal{DrugId: pharmacyDrug.DrugID, PharmacyId: pharmacyDrug.PharmacyID, Quantity: int(orderDetail.Quantity) * -1, Description: "transfered"})
		stockJournals = append(stockJournals, entity.StockJurnal{DrugId: pharmacyDrug.DrugID, PharmacyId: order.PharmacyId, Quantity: int(orderDetail.Quantity), Description: "received"})
		stockRequestDrug[pharmacyDrug.PharmacyID] = append(stockRequestDrug[pharmacyDrug.PharmacyID], &entity.StockRequestDrug{DrugId: pharmacyDrug.DrugID, Quantity: int(orderDetail.Quantity)})

	}

	err = u.stockRequesMutationAuto(ctx, order.PharmacyId, stockRequestDrug)
	if err != nil {
		return nil, err
	}
	stockRequestDrug[order.PharmacyId] = soldStock[order.PharmacyId]
	err = u.updateStock(ctx, stockRequestDrug, stockJournals)
	if err != nil {
		return nil, err
	}
	order, err = u.orderRepository.PMUpdateOrderStatusByOrderId(ctx, *order, constant.Processed, managerId)
	if err != nil {
		return nil, err
	}
	return order, nil
}
func (u *orderUsecaseImpl) stockRequesMutationAuto(ctx context.Context, pharmacyId uint, stockRequestDrug map[uint][]*entity.StockRequestDrug) error {
	stockRequest, err := u.stockRequestRepository.InsertStockRequestBulk(ctx, pharmacyId, stockRequestDrug, constant.Approved)
	if err != nil {
		return err
	}

	err = u.stockRequestDrugRepository.InsertStockRequestDrugBulk(ctx, stockRequest)
	if err != nil {
		return err
	}
	return nil
}

func (u *orderUsecaseImpl) updateStock(ctx context.Context, stockRequestDrug map[uint][]*entity.StockRequestDrug, stockJournals []entity.StockJurnal) error {
	err := u.pharmacyDrugrepository.UpdateSubstractionBulkStock(ctx, stockRequestDrug)
	if err != nil {
		return err
	}
	err = u.stockJournalRepository.InsertStockJournal(ctx, stockJournals)
	if err != nil {
		return err
	}
	return nil
}
func (u *orderUsecaseImpl) GetAllOrderByPharmacyManagerId(ctx context.Context) ([]*entity.Order, error) {
	managerCtx, ok := utils.CtxGetManager(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}
	managerId := managerCtx.ID
	orders, err := u.orderRepository.GetAllOrderByPharmacyManagerId(ctx, managerId)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
func (u *orderUsecaseImpl) OrderSent(ctx context.Context, order entity.Order) error {
	managerCtx, ok := utils.CtxGetManager(ctx)
	if !ok {
		return apperror.ErrInternalServer
	}
	managerId := managerCtx.ID
	order.Status = constant.Processed
	_, err := u.orderRepository.PMUpdateOrderStatusByOrderId(ctx, order, constant.Sent, managerId)
	if err != nil {
		return err
	}
	return nil
}

func (u *orderUsecaseImpl) OrderCancelByPM(ctx context.Context, order entity.Order) error {
	status, err := u.orderRepository.SelecOrderStatusByOrderId(ctx, order.Id)
	if err != nil {
		return err
	}

	if *status != constant.Processed && *status != constant.PaymentConfirmed {
		return apperror.CantCancelOrder
	}
	managerCtx, ok := utils.CtxGetManager(ctx)
	if !ok {
		return apperror.ErrInternalServer
	}
	managerId := managerCtx.ID
	order.Status = *status
	if *status == constant.Processed {
		orderDetails, err := u.orderDetailRepository.SelectOrderDetailByOrderId(ctx, order.Id)
		if err != nil {
			return err
		}
		stockJournals, err := u.pharmacyDrugrepository.UpdateReturnStock(ctx, orderDetails)
		if err != nil {
			return err
		}
		err = u.stockJournalRepository.InsertStockJournal(ctx, stockJournals)
		if err != nil {
			return err
		}
	}
	_, err = u.orderRepository.PMUpdateOrderStatusByOrderId(ctx, order, constant.Cancelled, managerId)
	if err != nil {
		return err
	}
	return nil
}
