package usecase

import (
	"context"
	"errors"
	"time"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/database/transaction"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/libs/firebase"
	"Alice-Seahat-Healthcare/seahat-be/libs/mail"
	"Alice-Seahat-Healthcare/seahat-be/repository"
	"Alice-Seahat-Healthcare/seahat-be/utils"
)

type DoctorUsecase interface {
	Login(ctx context.Context, doctor entity.Doctor) (string, error)
	LoginOAuth(ctx context.Context, googleToken string) (string, error)
	Register(ctx context.Context, body entity.Doctor) (*entity.Doctor, error)
	RegisterOAuth(ctx context.Context, dr entity.Doctor, googleToken string) (*entity.Doctor, error)
	CreateToken(ctx context.Context, doctor entity.Doctor, tokenType string) (*entity.Token, error)
	ForgotPassword(ctx context.Context, email string) error
	Verification(ctx context.Context, password string, token string) error
	UpdateProfile(ctx context.Context, doctor entity.Doctor) error
	UpdatePassword(ctx context.Context, oldPwd string, newPwd string) error
	GetProfile(ctx context.Context) (*entity.Doctor, error)
	UpdateStatus(ctx context.Context, doctor entity.Doctor) error
	GetDoctorByID(ctx context.Context, id uint) (*entity.Doctor, error)
	GetAllDoctors(ctx context.Context, clc *entity.Collection) ([]entity.Doctor, error)
}

type doctorUsecaseImpl struct {
	doctorRepository repository.DoctorRepository
	tokenRepository  repository.TokenRepository
	transactor       transaction.Transactor
	mail             mail.MailDialer
	firebase         firebase.Firebase
}

func NewDoctorUsecase(
	doctorRepository repository.DoctorRepository,
	tokenRepository repository.TokenRepository,
	transactor transaction.Transactor,
	mail mail.MailDialer,
	firebase firebase.Firebase,
) *doctorUsecaseImpl {
	return &doctorUsecaseImpl{
		doctorRepository: doctorRepository,
		tokenRepository:  tokenRepository,
		transactor:       transactor,
		mail:             mail,
		firebase:         firebase,
	}
}

func (u *doctorUsecaseImpl) Login(ctx context.Context, body entity.Doctor) (string, error) {
	doctor, err := u.doctorRepository.SelectOneByEmail(ctx, body.Email)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return "", apperror.InvalidCredential
		}

		return "", err
	}

	if !utils.HashCompareDefault(*body.Password, doctor.Password) {
		return "", apperror.InvalidCredential
	}

	jwtData := map[string]any{
		"ID":    doctor.ID,
		"Email": doctor.Email,
	}

	token, err := utils.JwtGenerateDoctor(jwtData)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *doctorUsecaseImpl) LoginOAuth(ctx context.Context, googleToken string) (string, error) {
	googleUser, err := u.firebase.GetAuthIdentity(ctx, googleToken)
	if err != nil {
		return "", apperror.InvalidToken
	}

	doctor, err := u.doctorRepository.SelectOneByEmail(ctx, googleUser["Email"])
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return "", apperror.EmailOAuthNotFound
		}

		return "", err
	}

	if !doctor.IsOAuth {
		return "", apperror.InvalidToken
	}

	jwtData := map[string]any{
		"ID":    doctor.ID,
		"Email": doctor.Email,
	}

	token, err := utils.JwtGenerateDoctor(jwtData)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *doctorUsecaseImpl) Register(ctx context.Context, body entity.Doctor) (*entity.Doctor, error) {
	_, err := u.getDoctorByEmail(ctx, body.Email)
	if !errors.Is(err, apperror.ErrResourceNotFound) {
		if err == nil {
			return nil, apperror.EmailExist
		}

		return nil, err
	}

	body.PhotoURL = constant.DefaultPhotoURL

	doctor, err := u.doctorRepository.InsertOne(ctx, body)
	if err != nil {
		return nil, err
	}

	token, err := u.CreateToken(ctx, *doctor, constant.TokenTypeConfirm)
	if err != nil {
		return nil, err
	}

	go utils.SendEmailVerificationDoctor(u.mail, doctor.Email, token.Token)
	return doctor, nil
}

func (u *doctorUsecaseImpl) RegisterOAuth(ctx context.Context, dr entity.Doctor, googleToken string) (*entity.Doctor, error) {
	googleUser, err := u.firebase.GetAuthIdentity(ctx, googleToken)
	if err != nil {
		return nil, apperror.InvalidToken
	}

	_, err = u.getDoctorByEmail(ctx, googleUser["Email"])
	if !errors.Is(err, apperror.ErrResourceNotFound) {
		if err == nil {
			return nil, apperror.EmailExist
		}

		return nil, err
	}

	dr.Email = googleUser["Email"]
	dr.PhotoURL = googleUser["Picture"]
	dr.IsOAuth = true
	dr.IsVerified = true

	doctor, err := u.doctorRepository.InsertOne(ctx, dr)
	if err != nil {
		return nil, err
	}

	return doctor, nil
}

func (u *doctorUsecaseImpl) CreateToken(ctx context.Context, doctor entity.Doctor, tokenType string) (*entity.Token, error) {
	if err := u.tokenRepository.DeleteByDoctorID(ctx, doctor.ID, tokenType); err != nil {
		return nil, err
	}

	tokenGenerated, err := utils.RandomString(constant.TokenLength)
	if err != nil {
		return nil, err
	}

	token, err := u.tokenRepository.InsertOne(ctx, entity.Token{
		DoctorID:  &doctor.ID,
		Type:      tokenType,
		Token:     tokenGenerated,
		ExpiredAt: time.Now().Add(constant.TokenDuration),
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (u *doctorUsecaseImpl) ForgotPassword(ctx context.Context, email string) error {
	doctor, err := u.doctorRepository.SelectOneByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return apperror.InvalidEmail
		}

		return err
	}

	if doctor.IsOAuth {
		return apperror.EmailCantResetPassword
	}

	token, err := u.CreateToken(ctx, *doctor, constant.TokenTypeReset)
	if err != nil {
		return err
	}

	go utils.SendEmailForgotToken(u.mail, doctor.Email, token.Token)
	return nil
}

func (u *doctorUsecaseImpl) Verification(ctx context.Context, password string, token string) error {
	t, err := u.tokenRepository.SelectOneByToken(ctx, entity.Token{
		Token: token,
		Type:  constant.TokenTypeConfirm,
	})

	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return apperror.InvalidToken
		}

		return err
	}

	if t.DoctorID == nil {
		return apperror.InvalidToken
	}

	hashPwd, err := utils.HashPasswordDefault(&password)
	if err != nil {
		return err
	}

	if err := u.doctorRepository.UpdatePassword(ctx, *t.DoctorID, *hashPwd); err != nil {
		return err
	}

	if err := u.tokenRepository.DeleteByToken(ctx, token); err != nil {
		return err
	}

	return nil
}

func (u *doctorUsecaseImpl) UpdateProfile(ctx context.Context, doctor entity.Doctor) error {
	doctorCtx, ok := utils.CtxGetDoctor(ctx)
	if !ok {
		return apperror.ErrInternalServer
	}

	doctor.ID = doctorCtx.ID
	err := u.doctorRepository.UpdatePersonalByID(ctx, doctor)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return apperror.ResourceNotFound
		}

		return err
	}

	return nil
}

func (u *doctorUsecaseImpl) UpdatePassword(ctx context.Context, oldPwd string, newPwd string) error {
	doctorCtx, ok := utils.CtxGetDoctor(ctx)
	if !ok {
		return apperror.ErrInternalServer
	}

	doctor, err := u.doctorRepository.SelectOneByEmail(ctx, doctorCtx.Email)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return apperror.ResourceNotFound
		}

		return err
	}

	if !utils.HashCompareDefault(oldPwd, doctor.Password) {
		return apperror.InvalidPassword
	}

	hashPwd, err := utils.HashPasswordDefault(&newPwd)
	if err != nil {
		return err
	}

	err = u.doctorRepository.UpdatePassword(ctx, doctor.ID, *hashPwd)
	if err != nil {
		return err
	}

	return nil
}

func (u *doctorUsecaseImpl) GetProfile(ctx context.Context) (*entity.Doctor, error) {
	doctorCtx, ok := utils.CtxGetDoctor(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}

	doctor, err := u.doctorRepository.SelectOneByEmail(ctx, doctorCtx.Email)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.ResourceNotFound
		}

		return nil, err
	}

	return doctor, nil
}

func (u *doctorUsecaseImpl) GetAllDoctors(ctx context.Context, clc *entity.Collection) ([]entity.Doctor, error) {
	doctors, err := u.doctorRepository.SelectAll(ctx, clc)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.ResourceNotFound
		}

		return nil, err
	}

	return doctors, nil
}

func (u *doctorUsecaseImpl) UpdateStatus(ctx context.Context, doctor entity.Doctor) error {
	doctorCtx, ok := utils.CtxGetDoctor(ctx)
	if !ok {
		return apperror.ErrInternalServer
	}

	doctor.ID = doctorCtx.ID
	if err := u.doctorRepository.UpdateStatus(ctx, doctor); err != nil {
		return err
	}

	return nil
}

func (u *doctorUsecaseImpl) getDoctorByEmail(ctx context.Context, email string) (*entity.Doctor, error) {
	d, err := u.doctorRepository.SelectOneByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if d.IsVerified {
		return d, nil
	}

	_, err = u.tokenRepository.SelectOneByDoctorID(ctx, entity.Token{
		DoctorID: &d.ID,
		Type:     constant.TokenTypeConfirm,
	})

	if err == nil {
		return d, nil
	}

	if errors.Is(err, apperror.ErrResourceNotFound) {
		if err := u.doctorRepository.DeleteByID(ctx, d.ID); err != nil {
			return nil, err
		}

		return nil, err
	}

	return nil, err
}

func (u *doctorUsecaseImpl) GetDoctorByID(ctx context.Context, id uint) (*entity.Doctor, error) {
	doctor, err := u.doctorRepository.SelectOneByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.ResourceNotFound
		}

		return nil, err
	}

	return doctor, nil

}
