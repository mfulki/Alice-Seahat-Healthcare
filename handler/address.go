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

type AddressHandler struct {
	addressUsecase usecase.AddressUsecase
}

func NewAddressHandler(addressUsecase usecase.AddressUsecase) *AddressHandler {
	return &AddressHandler{
		addressUsecase: addressUsecase,
	}
}

func (h *AddressHandler) GetAllProvinces(ctx *gin.Context) {
	provinces, err := h.addressUsecase.GetAllProvinces(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    response.NewMultipleProvinceDto(provinces),
	})
}

func (h *AddressHandler) GetAllCities(ctx *gin.Context) {
	provinceID := 0
	if id, err := strconv.Atoi(ctx.Query("province_id")); err == nil {
		provinceID = id
	}

	cities, err := h.addressUsecase.GetAllCitiesWithProvinceQuery(ctx, uint(provinceID))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    response.NewMultipleCityDto(cities),
	})
}

func (h *AddressHandler) GetAllSubdistrict(ctx *gin.Context) {
	cityID := 0
	if id, err := strconv.Atoi(ctx.Query("city_id")); err == nil {
		cityID = id
	}

	subdistricts, err := h.addressUsecase.GetAllSubdistrictsWithCityQuery(ctx, uint(cityID))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    response.NewMultipleSubdistrictDto(subdistricts),
	})
}

func (h *AddressHandler) GetAllAddress(ctx *gin.Context) {
	var isActive *bool
	if parsed, err := strconv.ParseBool(ctx.Query("is_active")); err == nil {
		isActive = &parsed
	}

	addrs, err := h.addressUsecase.GetAllUserAddress(ctx, isActive)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    response.NewMultipleAddressDto(addrs),
	})
}

func (h *AddressHandler) AddAddress(ctx *gin.Context) {
	body := new(request.AddressRequest)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	addr := body.Addr()
	if addr == nil {
		ctx.Error(apperror.MainMustActiveAddress)
		return
	}

	addedAddr, err := h.addressUsecase.AddAddress(ctx, *addr)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response.Body{
		Message: constant.DataCreatedMsg,
		Data:    response.NewAddressDto(*addedAddr),
	})
}

func (h *AddressHandler) GetAddressByID(ctx *gin.Context) {
	addressID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || addressID < 1 {
		ctx.Error(apperror.InvalidParam)
		return
	}

	addr, err := h.addressUsecase.GetAddressByID(ctx, uint(addressID))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    response.NewAddressDto(*addr),
	})
}

func (h *AddressHandler) UpdateAddress(ctx *gin.Context) {
	addressID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || addressID < 1 {
		ctx.Error(apperror.InvalidParam)
		return
	}

	body := new(request.AddressRequest)
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.Error(err)
		return
	}

	addr := body.Addr()
	if addr == nil {
		ctx.Error(apperror.MainMustActiveAddress)
		return
	}

	addr.ID = uint(addressID)
	err = h.addressUsecase.UpdateAddressByID(ctx, *addr)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataEditMsg,
	})
}

func (h *AddressHandler) DeleteAddress(ctx *gin.Context) {
	addressID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || addressID < 1 {
		ctx.Error(apperror.InvalidParam)
		return
	}

	err = h.addressUsecase.DeleteAddressByID(ctx, uint(addressID))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataDeletedMsg,
	})
}

func (h *AddressHandler) GetShipmentPriceByAddressID(ctx *gin.Context) {
	addressID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || addressID < 1 {
		ctx.Error(apperror.InvalidParam)
		return
	}

	pharmacyQuery := request.GetPharmacyShipmentPriceQuery(ctx)
	data, err := h.addressUsecase.GetShipmentPriceByPharmaciesID(ctx, uint(addressID), pharmacyQuery)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response.Body{
		Message: constant.DataRetrievedMsg,
		Data:    response.NewMultiplePharmacyWithShipmentPrice(data),
	})
}
