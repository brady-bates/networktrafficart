package simulation

import (
	"encoding/binary"
	"image/color"
	"math"
	"math/rand"
	"net"
	"networktrafficart/internal/capture"
	"networktrafficart/internal/util"
)

const (
	offScreenSpawnDistance = 25.0
	maxIPv4Bits            = math.MaxUint32
	speed                  = float32(6.0)
)

var (
	xSkewIntensity = float32(util.ClampValue(.4, 0.0, 1.0))
)

type Packet struct {
	X, Y   float32
	YDelta float32
	XSkew  float32
	Color  color.RGBA
}

func NewPacketFromEvent(e capture.PacketData, screenWidth, screenHeight int) Packet {
	rand0to1 := rand.Float32() - .5
	ip := binary.BigEndian.Uint32(e.SrcIP)
	ipRatio := float64(ip) / float64(maxIPv4Bits)

	xStart := float32(ipRatio) * float32(screenWidth)
	var yStart float32
	var ySpeed float32
	if e.IsIncoming {
		// Inbound - bottom to top
		yStart = float32(screenHeight) + offScreenSpawnDistance
		ySpeed = speed
	} else {
		// Outbound - top to bottom
		yStart = -offScreenSpawnDistance
		ySpeed = -speed
	}
	xSkew := rand0to1 * xSkewIntensity
	rgba := ipToRGBA(e.SrcIP)

	return Packet{
		xStart,
		yStart,
		ySpeed,
		xSkew,
		rgba,
	}
}

func ipToRGBA(src net.IP) color.RGBA {
	r := src[1]
	g := src[2]
	b := src[3]

	brightness := (uint32(r)*299 + uint32(g)*587 + uint32(b)*114) / 1000

	if brightness < 120 {
		r += 100
		g += 100
		b += 100
	}

	return color.RGBA{R: r, G: g, B: b, A: 255}
}
