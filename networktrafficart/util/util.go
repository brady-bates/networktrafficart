package util

import (
	"log"
	"math"
	"strconv"
)

func ClampValue(val, min, max float64) float64 {
	return math.Max(min, math.Min(val, max))
}

func IsTrueStr(s string) bool {
	return s == "true"
}

func ParseToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}
