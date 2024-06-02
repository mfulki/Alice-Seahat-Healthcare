package usecase

import (
	"context"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/entity"
	"Alice-Seahat-Healthcare/seahat-be/repository"
	"Alice-Seahat-Healthcare/seahat-be/utils"
)

type AdminReportUsecase interface {
	GetAdminDrugReport(ctx context.Context, clc *entity.Collection) ([]entity.AdminReportByDrug, error)
	GetAdminCategoryReport(ctx context.Context, clc *entity.Collection) ([]entity.AdminReportByCategory, error)
	GetManagerDrugReport(ctx context.Context, clc *entity.Collection) ([]entity.AdminReportByDrug, error)
	GetManagerCategoryReport(ctx context.Context, clc *entity.Collection) ([]entity.AdminReportByCategory, error)
}

type adminReportUsecaseImpl struct {
	adminReportRepository repository.AdminReportRepository
}

func NewAdminReportUsecase(
	adminReportRepository repository.AdminReportRepository,
) *adminReportUsecaseImpl {
	return &adminReportUsecaseImpl{
		adminReportRepository: adminReportRepository,
	}
}

func (u *adminReportUsecaseImpl) GetAdminDrugReport(ctx context.Context, clc *entity.Collection) ([]entity.AdminReportByDrug, error) {
	return u.adminReportRepository.GetDrugReportByCurrentMonth(ctx, 0, clc)
}

func (u *adminReportUsecaseImpl) GetAdminCategoryReport(ctx context.Context, clc *entity.Collection) ([]entity.AdminReportByCategory, error) {
	return u.adminReportRepository.GetCategoryReportByCurrentMonth(ctx, 0, clc)
}

func (u *adminReportUsecaseImpl) GetManagerDrugReport(ctx context.Context, clc *entity.Collection) ([]entity.AdminReportByDrug, error) {
	managerCtx, ok := utils.CtxGetManager(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}

	return u.adminReportRepository.GetDrugReportByCurrentMonth(ctx, managerCtx.ID, clc)
}

func (u *adminReportUsecaseImpl) GetManagerCategoryReport(ctx context.Context, clc *entity.Collection) ([]entity.AdminReportByCategory, error) {
	managerCtx, ok := utils.CtxGetManager(ctx)
	if !ok {
		return nil, apperror.ErrInternalServer
	}

	return u.adminReportRepository.GetCategoryReportByCurrentMonth(ctx, managerCtx.ID, clc)
}
