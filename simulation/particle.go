package simulation

import (
	"encoding/binary"
	"image/color"
	"math"
	"math/rand"
	"net"
	"networktrafficart/capture"
	"networktrafficart/util"
)

const (
	offScreenSpawnDistance = 25.0
)

type Particle struct {
	X, Y   float32
	YDelta float32
	XSkew  float32
	Color  color.RGBA
	Size   float32
}

func NewParticle(pe capture.PacketEvent, screenWidth, screenHeight int) *Particle {
	rand0to1 := rand.Float32() - .5
	maxIPv4Bits := float32(math.MaxUint32)
	ip := binary.BigEndian.Uint32(pe.SrcIP)
	packetBits := float32(pe.Size)

	xSkewIntensity := float32(util.ClampValue(.4, 0.0, 1.0))

	xStart := (float32(ip) / maxIPv4Bits) * float32(screenWidth) // TODO fix for ipv6
	yStart := float32(screenHeight) + offScreenSpawnDistance
	ySpeed := float32(7.0)
	xSkew := rand0to1 * xSkewIntensity
	rgba := ipv4ToRGBA(pe.SrcIP)

	return &Particle{
		xStart,
		yStart,
		ySpeed,
		xSkew,
		rgba,
		packetBits / 75,
	}
}

func ipv4ToRGBA(src net.IP) color.RGBA {
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
