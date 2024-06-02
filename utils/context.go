package utils

import (
	"context"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/entity"

	"github.com/sirupsen/logrus"
)

func CtxGetUser(ctx context.Context) (*entity.User, bool) {
	userMap, ok := getDetailActor(ctx, constant.User)
	if !ok {
		return nil, false
	}

	return &entity.User{
		ID:    uint(userMap["ID"].(float64)),
		Email: userMap["Email"].(string),
	}, true
}

func CtxGetDoctor(ctx context.Context) (*entity.Doctor, bool) {
	doctorMap, ok := getDetailActor(ctx, constant.Doctor)
	if !ok {
		return nil, false
	}

	return &entity.Doctor{
		ID:    uint(doctorMap["ID"].(float64)),
		Email: doctorMap["Email"].(string),
	}, true
}

func CtxGetManager(ctx context.Context) (*entity.PharmacyManager, bool) {
	managerMap, ok := getDetailActor(ctx, constant.Manager)
	if !ok {
		return nil, false
	}

	return &entity.PharmacyManager{
		ID:    uint(managerMap["ID"].(float64)),
		Email: managerMap["Email"].(string),
	}, true
}

func CtxGetAdmin(ctx context.Context) (*entity.Admin, bool) {
	adminMap, ok := getDetailActor(ctx, constant.Admin)
	if !ok {
		return nil, false
	}

	return &entity.Admin{
		ID:    uint(adminMap["ID"].(float64)),
		Email: adminMap["Email"].(string),
	}, true
}

func getDetailActor(ctx context.Context, actor string) (map[string]any, bool) {
	val := ctx.Value(constant.UserContext)

	act, ok := val.(map[string]any)
	if !ok {
		logrus.Error(apperror.ErrAssertingAny)
		return nil, false
	}

	if act["role"] != actor {
		return nil, false
	}

	return act, true
}
