package request

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPharmacyShipmentPriceQuery(ctx *gin.Context) map[uint]uint {
	pharmacies := ctx.QueryMap("pharmacy")
	result := make(map[uint]uint)

	for key, val := range pharmacies {
		intKey, err := strconv.Atoi(key)
		if err != nil {
			continue
		}

		intVal, err := strconv.Atoi(val)
		if err != nil {
			continue
		}

		result[uint(intKey)] = uint(intVal)
	}

	return result
}
