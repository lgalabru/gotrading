package core

import "math"

func Trunc8(m float64) float64 {
	return Trunc(m, 8)
}

func Trunc(m float64, d int) float64 {
	return math.Trunc(m*math.Pow(10, 8)) / math.Pow(10, 8)
}
