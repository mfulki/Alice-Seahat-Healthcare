package usecase

import (
	"context"
	"errors"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/libs/mail"
	"Alice-Seahat-Healthcare/seahat-be/repository"
	"Alice-Seahat-Healthcare/seahat-be/utils"
)

type PartnerUsecase interface {
	GetAll(ctx context.Context, clc *entity.Collection) ([]entity.Partner, error)
	CreatePartner(ctx context.Context, p entity.Partner) (*entity.Partner, error)
	GetPartnerByID(ctx context.Context, id uint) (*entity.Partner, error)
	UpdatePartnerByID(ctx context.Context, partner entity.Partner) error
}

type partnerUsecaseImpl struct {
	pharmacyManagerRepository repository.PharmacyManagerRepository
	partnerRepository         repository.PartnerRepository
	transactor                transaction.Transactor
	mail                      mail.MailDialer
}

func NewPartnerUsecase(
	pharmacyManagerRepository repository.PharmacyManagerRepository,
	partnerRepository repository.PartnerRepository,
	transactor transaction.Transactor,
	mail mail.MailDialer,
) *partnerUsecaseImpl {
	return &partnerUsecaseImpl{
		pharmacyManagerRepository: pharmacyManagerRepository,
		partnerRepository:         partnerRepository,
		transactor:                transactor,
		mail:                      mail,
	}
}

func (u *partnerUsecaseImpl) GetAll(ctx context.Context, clc *entity.Collection) ([]entity.Partner, error) {
	return u.partnerRepository.GetAll(ctx, clc)
}

func (u *partnerUsecaseImpl) CreatePartner(ctx context.Context, p entity.Partner) (*entity.Partner, error) {
	email := p.PharmacyManager.Email
	_, err := u.pharmacyManagerRepository.SelectOneByEmail(ctx, email)
	if !errors.Is(err, apperror.ErrResourceNotFound) {
		if err == nil {
			return nil, apperror.EmailExist
		}

		return nil, err
	}

	pwdGenerated, err := utils.RandomString(constant.GeneratedPasswordLength)
	if err != nil {
		return nil, err
	}

	hashPwd, err := utils.HashPasswordDefault(&pwdGenerated)
	if err != nil {
		return nil, err
	}

	p.PharmacyManager.Password = *hashPwd
	manager, err := u.pharmacyManagerRepository.InsertOne(ctx, p.PharmacyManager)
	if err != nil {
		return nil, err
	}

	p.PharmacyManagerID = manager.ID
	partner, err := u.partnerRepository.InsertOne(ctx, p)
	if err != nil {
		return nil, err
	}

	partner.PharmacyManager = *manager
	go utils.SendEmailAddPartner(u.mail, *partner, pwdGenerated)

	return partner, nil
}

func (u *partnerUsecaseImpl) GetPartnerByID(ctx context.Context, id uint) (*entity.Partner, error) {
	p, err := u.partnerRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.ResourceNotFound
		}

		return nil, err
	}

	return p, nil
}

func (u *partnerUsecaseImpl) UpdatePartnerByID(ctx context.Context, partner entity.Partner) error {
	p, err := u.partnerRepository.UpdateByID(ctx, partner)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return apperror.ResourceNotFound
		}

		return err
	}

	partner.PharmacyManager.ID = p.PharmacyManagerID

	err = u.pharmacyManagerRepository.UpdateByID(ctx, partner.PharmacyManager)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return apperror.ResourceNotFound
		}

		return err
	}

	return nil
}
