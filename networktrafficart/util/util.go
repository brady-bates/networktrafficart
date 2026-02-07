package util

import (
	"log"
	"math"
	"math/rand"
	"net"
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

func GenerateRandomIPv4() net.IP {
	o1 := byte(rand.Intn(256))
	o2 := byte(rand.Intn(256))
	o3 := byte(rand.Intn(256))
	o4 := byte(rand.Intn(256))

	return net.IPv4(o1, o2, o3, o4).To4()
}
