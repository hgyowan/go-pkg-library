package variable

import "math"

func GetSafeValue[TYPE any](targetValue *TYPE, defaultValue TYPE) TYPE {
	if targetValue != nil {
		return *targetValue
	}

	return defaultValue
}

func ConvertToPointer[TYPE any](defaultValue TYPE) *TYPE {
	return &defaultValue
}

func GetSafeNaNValue(value float64) float64 {
	if math.IsNaN(value) {
		return 0
	}

	return value
}

func ToSortableNumber(name string) int64 {
	var sum int64
	for _, r := range name {
		sum = sum*1000 + int64(r)
	}
	return sum
}
