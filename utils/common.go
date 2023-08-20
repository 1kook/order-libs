package utils

import "math"

func Round(num float64, n int) float64 {
	if n < 0 {
		n = 0
	}

	factor := math.Pow(10, float64(n))
	return math.Trunc(num*factor) / factor
}
