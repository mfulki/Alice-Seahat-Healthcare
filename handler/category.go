package handler

import (
	"net/http"
	"strconv"

	"Alice-Seahat-Healthcare/seahat-be/apperror"
	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/dto/request"
	"Alice-Seahat-Healthcare/seahat-be/dto/response"
	"Alice-Seahat-Healthcare/seahat-be/usecase"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryUsecase usecase.CategoryUsecase
}

func NewCategoryHandler(categoryUsecase usecase.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{
		categoryUsecase: categoryUsecase,
	}
}

func (h *CategoryHandler) GetAllCategory(ctx *gin.Context) {
	collection := request.GetCollectionQuery(ctx)
	cats, err := h.categoryUsecase.GetAll(ctx, &collection)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       response.NewMultipleCategoryDto(cats),
		Pagination: response.NewPaginationDto(collection),
	})
}

func (h *CategoryHandler) CreateCategory(ctx *gin.Context) {
	body := new(request.CategoryRequest)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	cat, err := h.categoryUsecase.CreateCategory(ctx, body.Category())
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response.Body{
		Message: constant.DataCreatedMsg,
		Data:    response.NewCategoryDto(*cat),
	})
}

func (h *CategoryHandler) GetCategoryByID(ctx *gin.Context) {
	categoryID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || categoryID < 1 {
		ctx.Error(apperror.InvalidParam)
		return
	}

	cat, err := h.categoryUsecase.GetCategoryByID(ctx, uint(categoryID))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    response.NewCategoryDto(*cat),
	})
}

func (h *CategoryHandler) UpdateCategoryByID(ctx *gin.Context) {
	categoryID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || categoryID < 1 {
		ctx.Error(apperror.InvalidParam)
		return
	}

	body := new(request.CategoryRequest)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	cat := body.Category()
	cat.ID = uint(categoryID)

	err = h.categoryUsecase.UpdateCategory(ctx, cat)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataEditMsg,
	})
}

func (h *CategoryHandler) DeleteCategoryByID(ctx *gin.Context) {
	categoryID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || categoryID < 1 {
		ctx.Error(apperror.InvalidParam)
		return
	}

	err = h.categoryUsecase.DeleteCategory(ctx, uint(categoryID))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataDeletedMsg,
	})
}
