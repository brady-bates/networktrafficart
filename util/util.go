package util

import (
	"errors"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
)

func ClampValue(val, min, max float64) float64 {
	return math.Max(min, math.Min(val, max))
}

func IsTrueStr(s string) bool {
	return strings.TrimSpace(s) == "true"
}

func ParseToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func ParseToFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func GenerateRandomIPv4() net.IP {
	o1 := byte(rand.Intn(256))
	o2 := byte(rand.Intn(256))
	o3 := byte(rand.Intn(256))
	o4 := byte(rand.Intn(256))

	return net.IPv4(o1, o2, o3, o4).To4()
}

func FileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
