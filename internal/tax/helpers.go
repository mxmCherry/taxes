package tax

import (
	"math"
	"time"
)

func QuarterOf(month time.Month) int {
	return int(month-1)/3 + 1
}

func Round(x float64, decimals int) float64 {
	if decimals <= 0 {
		return x
	}

	exp := math.Pow10(decimals)
	return math.Round(x*exp) / exp
}
