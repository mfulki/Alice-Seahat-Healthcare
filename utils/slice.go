package utils

import "Alice-Seahat-Healthcare/seahat-be/constant"

type sliceType interface {
	string | int
}

func SliceIsContain[V sliceType](slice []V, target V) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}

	return false
}

func Num2OperationDay(num uint) []bool {
	if num == 0 {
		return []bool{}
	}

	days := make([]bool, constant.TotalDays)

	for i := constant.TotalDays - 1; i >= 0; i-- {
		days[i] = num&1 == 1
		num = num >> 1
	}

	return days
}

func OperationDay2Num(days []bool) uint {
	daysLength := len(days)
	if daysLength != constant.TotalDays {
		return 0
	}

	num := 0
	for index, day := range days {
		numDay := 0
		if day {
			numDay = 1
		}

		num = num | numDay
		if index != daysLength-1 {
			num = num << 1
		}
	}

	return uint(num)
}
