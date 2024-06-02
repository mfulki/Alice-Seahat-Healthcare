package usecase

import (
	"context"
	"errors"
	"sort"
	"time"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/repository"
	"Alice-Seahat-Healthcare/seahat-be/utils"
)

type PaymentUsecase interface {
	UpdatePaymentProof(ctx context.Context, body entity.Payment) ([]*entity.Order, error)
	UserCancelPayment(ctx context.Context, body entity.Payment) ([]*entity.Order, error)
	GetAllPaymentToConfirm(ctx context.Context, clc *entity.Collection) ([]*entity.Payment, error)
	PaymentConfirmation(ctx context.Context, body entity.Payment) ([]*entity.Order, error)
	GetAllPaymentByUserId(ctx context.Context, clc *entity.Collection) ([]*entity.Payment, error)
	AdminCancelPayment(ctx context.Context, body entity.Payment) ([]*entity.Order, error)
	AdminRejectPayment(ctx context.Context, body entity.Payment) error
}

type paymentUsecaseImpl struct {
	paymentrepository repository.PaymentRepository
	orderRepository   repository.OrderRepository
	transactor        transaction.Transactor
}

func NewPaymentUsecase(
	paymentrepository repository.PaymentRepository,
	orderRepository repository.OrderRepository,
	transactor transaction.Transactor,

) *paymentUsecaseImpl {
	return &paymentUsecaseImpl{
		paymentrepository: paymentrepository,
		orderRepository:   orderRepository,
		transactor:        transactor,
	}
}

func (u *paymentUsecaseImpl) UpdatePaymentProof(ctx context.Context, body entity.Payment) ([]*entity.Order, error) {
	userCtx, ok := utils.CtxGetUser(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}
	body.UserId = userCtx.ID
	ordersTx, err := u.transactor.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		payment, err := u.paymentrepository.UpdatePaymentProof(ctx, body)
		if err != nil {
			if errors.Is(err, apperror.ErrResourceNotFound) {
				return nil, apperror.ResourceNotFound
			}

			return nil, err
		}
		futureStatus := constant.WaitingForPaymentConfirmation
		recentStatus := constant.WaitingForPayment
		orders, err := u.orderRepository.UpdateOrderStatusByPaymentId(ctx, *payment, futureStatus, recentStatus)
		if err != nil {
			return nil, err
		}
		return orders, nil
	})
	if err != nil {
		return nil, err
	}
	orders := ordersTx.([]*entity.Order)
	return orders, nil
}
func (u *paymentUsecaseImpl) UserCancelPayment(ctx context.Context, body entity.Payment) ([]*entity.Order, error) {
	userCtx, ok := utils.CtxGetUser(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}
	body.UserId = userCtx.ID
	ordersTx, err := u.transactor.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		payment, err := u.paymentrepository.UserDeletePayment(ctx, body)
		if err != nil {
			if errors.Is(err, apperror.ErrResourceNotFound) {
				return nil, apperror.ResourceNotFound
			}

			return nil, err
		}
		futureStatus := constant.Cancelled
		recentStatus := constant.WaitingForPayment
		orders, err := u.orderRepository.UpdateOrderStatusByPaymentId(ctx, *payment, futureStatus, recentStatus)
		if err != nil {
			return nil, err
		}
		return orders, nil
	})
	if err != nil {
		return nil, err
	}
	orders := ordersTx.([]*entity.Order)
	return orders, nil
}
func (u *paymentUsecaseImpl) AdminRejectPayment(ctx context.Context, body entity.Payment) error {
	_, err := u.transactor.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		futureStatus := constant.WaitingForPayment
		recentStatus := constant.WaitingForPaymentConfirmation
		err := u.paymentrepository.UpdatePaymentExpiredAt(ctx, body.Id, futureStatus)
		if err != nil {
			if errors.Is(err, apperror.ErrResourceNotFound) {
				return nil, apperror.ResourceNotFound
			}

			return nil, err
		}
		orders, err := u.orderRepository.UpdateOrderStatusByPaymentId(ctx, body, futureStatus, recentStatus)
		if err != nil {
			return nil, err
		}
		return orders, nil

	})
	return err
}

func (u *paymentUsecaseImpl) AdminCancelPayment(ctx context.Context, body entity.Payment) ([]*entity.Order, error) {
	ordersTx, err := u.transactor.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		payment, err := u.paymentrepository.AdminDeletePayment(ctx, body)
		if err != nil {
			if errors.Is(err, apperror.ErrResourceNotFound) {
				return nil, apperror.ResourceNotFound
			}

			return nil, err
		}
		futureStatus := constant.Cancelled
		recentStatus := constant.WaitingForPayment
		orders, err := u.orderRepository.UpdateOrderStatusByPaymentId(ctx, *payment, futureStatus, recentStatus)
		if err != nil {
			return nil, err
		}
		return orders, nil
	})
	if err != nil {
		return nil, err
	}
	orders := ordersTx.([]*entity.Order)
	return orders, nil
}
func (u *paymentUsecaseImpl) PaymentConfirmation(ctx context.Context, body entity.Payment) ([]*entity.Order, error) {
	futureStatus := constant.PaymentConfirmed
	recentStatus := constant.WaitingForPaymentConfirmation
	ordersTx, err := u.transactor.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		orders, err := u.orderRepository.UpdateOrderStatusByPaymentId(ctx, body, futureStatus, recentStatus)
		if err != nil {
			return nil, err
		}
		err = u.paymentrepository.UpdatePaymentExpiredAt(ctx, body.Id, futureStatus)
		if err != nil {
			return nil, err
		}
		return orders, nil
	})
	if err != nil {
		return nil, err
	}
	orders := ordersTx.([]*entity.Order)
	return orders, nil
}

func (u *paymentUsecaseImpl) GetAllPaymentToConfirm(ctx context.Context, clc *entity.Collection) ([]*entity.Payment, error) {
	payments, err := u.paymentrepository.GetAllPaymentToConfirm(ctx, clc)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (u *paymentUsecaseImpl) GetAllPaymentByUserId(ctx context.Context, clc *entity.Collection) ([]*entity.Payment, error) {
	userCtx, ok := utils.CtxGetUser(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}
	userId := userCtx.ID

	payments, err := u.paymentrepository.GetAllPaymentByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	payment := make([]*entity.Payment, 0)
	keys := u.sortingPaymentMapKey(payments)

	for _, k := range keys {
		paymentWithStatus := u.getPaymentStatus(payments[k])
		payment = append(payment, paymentWithStatus)

	}

	return utils.HardPagination(payment, clc), nil
}

func (u *paymentUsecaseImpl) getPaymentStatus(payment *entity.Payment) *entity.Payment {
	if payment.DeletedAt != nil {
		payment.Status = constant.Cancelled
		return payment
	}
	if payment.ExpiredAt != nil {
		if payment.ExpiredAt.Time.After(time.Now()) && payment.Proof != nil {
			payment.Status = constant.WaitingForPaymentConfirmation
			return payment
		}
		if payment.ExpiredAt.Time.After(time.Now()) {
			payment.Status = constant.WaitingForPayment
			return payment
		}
		if payment.ExpiredAt.Time.Before(time.Now()) && payment.Proof != nil {
			payment.Status = constant.WaitingForPaymentConfirmation
			return payment
		}
		if payment.ExpiredAt.Time.Before(time.Now()) {
			payment.Status = constant.PaymentExpired
			return payment
		}
	}
	if payment.Proof != nil {
		payment.Status = constant.PaymentConfirmed
		return payment
	}
	payment.Status = constant.InvalidPayment
	return payment
}

func (u *paymentUsecaseImpl) sortingPaymentMapKey(payments map[uint]*entity.Payment) []uint {
	keys := make([]uint, 0, len(payments))
	for k := range payments {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	return keys
}
