package handler

import (
	"net/http"

	"Alice-Seahat-Healthcare/seahat-be/constant"
	"Alice-Seahat-Healthcare/seahat-be/dto/request"
	"Alice-Seahat-Healthcare/seahat-be/dto/response"
	"Alice-Seahat-Healthcare/seahat-be/usecase"

	"github.com/gin-gonic/gin"
)

type DrugHandler struct {
	drugUsecase usecase.DrugUsecase
}

func NewDrugHandler(drugUsecase usecase.DrugUsecase) *DrugHandler {
	return &DrugHandler{
		drugUsecase: drugUsecase,
	}
}

func (h *DrugHandler) GetAll(ctx *gin.Context) {
	collection := request.GetCollectionQuery(ctx)

	drugs, err := h.drugUsecase.GetAll(ctx, &collection)
	if err != nil {
		ctx.Error(err)
		return
	}

	res := make([]*response.DrugDto, 0)
	for _, drug := range drugs {
		res = append(res, response.NewDrugDto(*drug))
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message:    constant.DataRetrievedMsg,
		Data:       res,
		Pagination: response.NewPaginationDto(collection),
	})
}

func (h *DrugHandler) GetOneById(ctx *gin.Context) {
	req := new(request.GetIdUri)
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.Error(err)
		return
	}

	drug, err := h.drugUsecase.GetOneById(ctx, req.Drugs())
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    response.NewDrugDto(*drug),
	})
}

func (h *DrugHandler) InsertOne(ctx *gin.Context) {
	req := new(request.DrugDTO)
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	drugReq := req.ConvertIntoEntityDrugs()
	drug, err := h.drugUsecase.InsertOne(ctx, *drugReq)
	if err != nil {
		ctx.Error(err)
		return

	}

	ctx.JSON(http.StatusCreated, response.Body{
		Message: constant.DataCreatedMsg,
		Data:    response.NewDrugDto(*drug),
	})

}

func (h *DrugHandler) UpdateOne(ctx *gin.Context) {
	var req request.DrugDTO
	var drugIdUri request.GetIdUri

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	if err := ctx.ShouldBindUri(&drugIdUri); err != nil {
		ctx.Error(err)
		return
	}

	drugReq := req.ConvertMultipleReqIntoEntityDrugs(drugIdUri)
	drug, err := h.drugUsecase.UpdateOne(ctx, *drugReq)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataEditMsg,
		Data:    response.NewDrugDto(*drug),
	})
}
