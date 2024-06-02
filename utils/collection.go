package utils

import (
	"math"

	"Alice-Seahat-Healthcare/seahat-be/entity"
)

func HardPagination[T any](datas []T, collection *entity.Collection) []T {
	if collection.Page == 0 || collection.Limit == 0 {
		return datas
	}

	totalData := len(datas)
	firstIndex := (collection.Page - 1) * collection.Limit
	collection.TotalRecords = uint(totalData)

	if firstIndex > uint(totalData) {
		return []T{}
	}

	lastIndex := firstIndex + collection.Limit
	minLastIndex := math.Min(float64(totalData), float64(lastIndex))

	return datas[firstIndex:uint(minLastIndex)]
}
