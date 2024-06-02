package request

import (
	"strconv"

	"Alice-Seahat-Healthcare/seahat-be/entity"

	"github.com/gin-gonic/gin"
)

func GetCollectionQuery(ctx *gin.Context) entity.Collection {
	page, limit := 0, 0

	if p, err := strconv.Atoi(ctx.Query("page")); err == nil && p >= 1 {
		page = p
	}

	if l, err := strconv.Atoi(ctx.Query("limit")); err == nil && l >= 1 {
		limit = l
	}

	return entity.Collection{
		Filter: ctx.QueryMap("filter"),
		Sort:   ctx.Query("sort"),
		Search: ctx.Query("search"),
		Page:   uint(page),
		Limit:  uint(limit),
	}
}
