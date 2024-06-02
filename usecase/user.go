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

type UserUsecase interface {
	Login(ctx context.Context, user entity.User) (string, error)
	LoginOAuth(ctx context.Context, googleToken string) (string, error)
	Register(ctx context.Context, body entity.User) (*entity.User, error)
	RegisterOAuth(ctx context.Context, usr entity.User, googleToken string) (*entity.User, error)
	CreateToken(ctx context.Context, user entity.User, tokenType string) (*entity.Token, error)
	Verification(ctx context.Context, password string, token string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, password string, token string) error
	UpdateProfile(ctx context.Context, user entity.User) error
	UpdatePassword(ctx context.Context, oldPwd string, newPwd string) error
	GetProfile(ctx context.Context) (*entity.User, error)
	ResendVerification(ctx context.Context) error
}

type userUsecaseImpl struct {
	userRepository   repository.UserRepository
	doctorRepository repository.DoctorRepository
	tokenRepository  repository.TokenRepository
	transactor       transaction.Transactor
	mail             mail.MailDialer
	firebase         firebase.Firebase
}

func NewUserUsecase(
	userRepository repository.UserRepository,
	doctorRepository repository.DoctorRepository,
	tokenRepository repository.TokenRepository,
	transactor transaction.Transactor,
	mail mail.MailDialer,
	firebase firebase.Firebase,
) *userUsecaseImpl {
	return &userUsecaseImpl{
		userRepository:   userRepository,
		doctorRepository: doctorRepository,
		tokenRepository:  tokenRepository,
		transactor:       transactor,
		mail:             mail,
		firebase:         firebase,
	}
}

func (u *userUsecaseImpl) Login(ctx context.Context, body entity.User) (string, error) {
	user, err := u.userRepository.SelectOneByEmail(ctx, body.Email)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return "", apperror.InvalidCredential
		}

		return "", err
	}

	if !utils.HashCompareDefault(*body.Password, user.Password) {
		return "", apperror.InvalidCredential
	}

	jwtData := map[string]any{
		"ID":    user.ID,
		"Email": user.Email,
	}

	token, err := utils.JwtGenerateUser(jwtData)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *userUsecaseImpl) LoginOAuth(ctx context.Context, googleToken string) (string, error) {
	googleUser, err := u.firebase.GetAuthIdentity(ctx, googleToken)
	if err != nil {
		return "", apperror.InvalidToken
	}

	user, err := u.userRepository.SelectOneByEmail(ctx, googleUser["Email"])
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return "", apperror.EmailOAuthNotFound
		}

		return "", err
	}

	if !user.IsOAuth {
		return "", apperror.InvalidToken
	}

	jwtData := map[string]any{
		"ID":    user.ID,
		"Email": user.Email,
	}

	token, err := utils.JwtGenerateUser(jwtData)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *userUsecaseImpl) Register(ctx context.Context, body entity.User) (*entity.User, error) {
	_, err := u.userRepository.SelectOneByEmail(ctx, body.Email)
	if !errors.Is(err, apperror.ErrResourceNotFound) {
		if err == nil {
			return nil, apperror.EmailExist
		}

		return nil, err
	}

	body.PhotoURL = constant.DefaultPhotoURL
	body.Password, err = utils.HashPasswordDefault(body.Password)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepository.InsertOne(ctx, body)
	if err != nil {
		return nil, err
	}

	token, err := u.CreateToken(ctx, *user, constant.TokenTypeConfirm)
	if err != nil {
		return nil, err
	}

	go utils.SendEmailVerification(u.mail, user.Email, token.Token)

	return user, nil
}

func (u *userUsecaseImpl) RegisterOAuth(ctx context.Context, usr entity.User, googleToken string) (*entity.User, error) {
	googleUser, err := u.firebase.GetAuthIdentity(ctx, googleToken)
	if err != nil {
		return nil, apperror.InvalidToken
	}

	_, err = u.userRepository.SelectOneByEmail(ctx, googleUser["Email"])
	if !errors.Is(err, apperror.ErrResourceNotFound) {
		if err == nil {
			return nil, apperror.EmailExist
		}

		return nil, err
	}

	usr.Email = googleUser["Email"]
	usr.PhotoURL = googleUser["Picture"]
	usr.IsOAuth = true
	usr.IsVerified = true

	user, err := u.userRepository.InsertOne(ctx, usr)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecaseImpl) CreateToken(ctx context.Context, user entity.User, tokenType string) (*entity.Token, error) {
	if err := u.tokenRepository.DeleteByUserID(ctx, user.ID, tokenType); err != nil {
		return nil, err
	}

	tokenGenerated, err := utils.RandomString(constant.TokenLength)
	if err != nil {
		return nil, err
	}

	token, err := u.tokenRepository.InsertOne(ctx, entity.Token{
		UserID:    &user.ID,
		Type:      tokenType,
		Token:     tokenGenerated,
		ExpiredAt: time.Now().Add(constant.TokenDuration),
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (u *userUsecaseImpl) Verification(ctx context.Context, password string, token string) error {
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

	if t.UserID == nil {
		return apperror.InvalidToken
	}

	user, err := u.userRepository.SelectOneByID(ctx, *t.UserID)
	if err != nil {
		return err
	}

	if !utils.HashCompareDefault(password, user.Password) {
		return apperror.InvalidPassToken
	}

	user.IsVerified = true

	if _, err := u.userRepository.UpdateOne(ctx, *user); err != nil {
		return err
	}

	if err := u.tokenRepository.DeleteByToken(ctx, token); err != nil {
		return err
	}

	return nil
}

func (u *userUsecaseImpl) ForgotPassword(ctx context.Context, email string) error {
	user, err := u.userRepository.SelectOneByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return apperror.InvalidEmail
		}

		return err
	}

	if user.IsOAuth {
		return apperror.EmailCantResetPassword
	}

	token, err := u.CreateToken(ctx, *user, constant.TokenTypeReset)
	if err != nil {
		return err
	}

	go utils.SendEmailForgotToken(u.mail, user.Email, token.Token)
	return nil
}

func (u *userUsecaseImpl) ResetPassword(ctx context.Context, password string, token string) error {
	t, err := u.tokenRepository.SelectOneByToken(ctx, entity.Token{
		Token: token,
		Type:  constant.TokenTypeReset,
	})

	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return apperror.InvalidToken
		}

		return err
	}

	var id = t.UserID
	var repository repository.ActorRepository = u.userRepository

	if t.UserID == nil {
		id = t.DoctorID
		repository = u.doctorRepository
	}

	hashedPwd, err := utils.HashPasswordDefault(&password)
	if err != nil {
		return err
	}

	if err := repository.UpdatePassword(ctx, *id, *hashedPwd); err != nil {
		return err
	}

	if err := u.tokenRepository.DeleteByToken(ctx, token); err != nil {
		return err
	}

	return nil
}

func (u *userUsecaseImpl) UpdateProfile(ctx context.Context, user entity.User) error {
	userCtx, ok := utils.CtxGetUser(ctx)
	if !ok {
		return apperror.ErrInternalServer
	}

	user.ID = userCtx.ID
	err := u.userRepository.UpdatePersonalByID(ctx, user)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return apperror.ResourceNotFound
		}

		return err
	}

	return nil
}

func (u *userUsecaseImpl) UpdatePassword(ctx context.Context, oldPwd string, newPwd string) error {
	userCtx, ok := utils.CtxGetUser(ctx)
	if !ok {
		return apperror.ErrInternalServer
	}

	user, err := u.userRepository.SelectOneByID(ctx, userCtx.ID)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return apperror.ResourceNotFound
		}

		return err
	}

	if !utils.HashCompareDefault(oldPwd, user.Password) {
		return apperror.InvalidPassword
	}

	hashPwd, err := utils.HashPasswordDefault(&newPwd)
	if err != nil {
		return err
	}

	err = u.userRepository.UpdatePassword(ctx, user.ID, *hashPwd)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUsecaseImpl) GetProfile(ctx context.Context) (*entity.User, error) {
	userCtx, ok := utils.CtxGetUser(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}

	user, err := u.userRepository.SelectOneByID(ctx, userCtx.ID)
	if err != nil {
		if errors.Is(err, apperror.ErrResourceNotFound) {
			return nil, apperror.ResourceNotFound
		}

		return nil, err
	}

	return user, nil
}

func (u *userUsecaseImpl) ResendVerification(ctx context.Context) error {
	userCtx, ok := utils.CtxGetUser(ctx)
	if !ok {
		return apperror.ErrInternalServer
	}

	user, err := u.userRepository.SelectOneByEmail(ctx, userCtx.Email)
	if err != nil {
		return err
	}

	if user.IsVerified {
		return apperror.HasBeenVerified
	}

	token, err := u.CreateToken(ctx, *user, constant.TokenTypeConfirm)
	if err != nil {
		return err
	}

	go utils.SendEmailVerification(u.mail, userCtx.Email, token.Token)

	return nil
}
