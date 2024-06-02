package handler

import (
	"net/http"

	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/dto/request"
	"Alice-Seahat-Healthcare/seahat-be/dto/response"
	"Alice-Seahat-Healthcare/seahat-be/usecase"

	"github.com/gin-gonic/gin"
)

type AdminReportHandler struct {
	adminReportUsecase usecase.AdminReportUsecase
}

func NewAdminReportHandler(adminReportUsecase usecase.AdminReportUsecase) *AdminReportHandler {
	return &AdminReportHandler{
		adminReportUsecase: adminReportUsecase,
	}
}

func (h *AdminReportHandler) GetAdminDrugReport(ctx *gin.Context) {
	collection := request.GetCollectionQuery(ctx)
	adminReports, err := h.adminReportUsecase.GetAdminDrugReport(ctx, &collection)
	if err != nil {
		ctx.Error(err)
		return
	}

	adminReportsDTO := []response.AdminDrugReportDTO{}
	for i := 0; i < len(adminReports); i++ {
		adminReportsDTO = append(adminReportsDTO, response.NewAdminDrugReportDTO(adminReports[i]))
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       adminReportsDTO,
		Pagination: response.NewPaginationDto(collection),
	})
}

func (h *AdminReportHandler) GetAdminCategoryReport(ctx *gin.Context) {
	collection := request.GetCollectionQuery(ctx)
	adminReports, err := h.adminReportUsecase.GetAdminCategoryReport(ctx, &collection)
	if err != nil {
		ctx.Error(err)
		return
	}

	adminReportsDTO := []response.AdminCategoryReportDTO{}
	for i := 0; i < len(adminReports); i++ {
		adminReportsDTO = append(adminReportsDTO, response.NewAdminCategoryReportDTO(adminReports[i]))
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       adminReportsDTO,
		Pagination: response.NewPaginationDto(collection),
	})
}

func (h *AdminReportHandler) GetManagerDrugReport(ctx *gin.Context) {
	collection := request.GetCollectionQuery(ctx)
	adminReports, err := h.adminReportUsecase.GetManagerDrugReport(ctx, &collection)
	if err != nil {
		ctx.Error(err)
		return
	}

	adminReportsDTO := []response.AdminDrugReportDTO{}
	for i := 0; i < len(adminReports); i++ {
		adminReportsDTO = append(adminReportsDTO, response.NewAdminDrugReportDTO(adminReports[i]))
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       adminReportsDTO,
		Pagination: response.NewPaginationDto(collection),
	})
}

func (h *AdminReportHandler) GetManagerCategoryReport(ctx *gin.Context) {
	collection := request.GetCollectionQuery(ctx)
	adminReports, err := h.adminReportUsecase.GetManagerCategoryReport(ctx, &collection)
	if err != nil {
		ctx.Error(err)
		return
	}

	adminReportsDTO := []response.AdminCategoryReportDTO{}
	for i := 0; i < len(adminReports); i++ {
		adminReportsDTO = append(adminReportsDTO, response.NewAdminCategoryReportDTO(adminReports[i]))
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       adminReportsDTO,
		Pagination: response.NewPaginationDto(collection),
	})
}
