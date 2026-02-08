package particle

import (
	"encoding/binary"
	"image/color"
	"math/rand"
	"net"
	"networktrafficart/capture"
	"networktrafficart/capture/packetevent"
	"networktrafficart/util"
)

const (
	offScreenDistance = 25.0
)

type Particle struct {
	X, Y   float32
	YDelta float32
	XSkew  float32
	Color  color.RGBA
	Size   float32
}

//func TickParticle

func CreateFromPacketEvent(pe packetevent.PacketEvent, screenWidth, screenHeight int16) *Particle {
	ip := binary.BigEndian.Uint32(pe.SrcIP)
	xSkewIntensity := float32(util.ClampValue(.4, 0.0, 1.0))
	xStart := (float32(ip) / capture.IPv4Range) * float32(screenWidth) // TODO fix for ipv6

	return &Particle{
		xStart,
		float32(screenHeight) + offScreenDistance,
		5,
		(rand.Float32() - .5) * xSkewIntensity,
		ipv4ToRGBA(pe.SrcIP),
		float32(pe.Size) / 75,
	}
}

func ipv4ToRGBA(src net.IP) color.RGBA {
	r := src[0]
	g := src[1]
	b := src[2]

	// TODO Make sure the colors don't get mangled by this, want them to actually be representative
	brightness := (uint32(r)*299 + uint32(g)*587 + uint32(b)*114) / 1000

	if brightness < 120 {
		r += 100
		g += 100
		b += 100
	}

	return color.RGBA{R: r, G: g, B: b, A: 255}
}
