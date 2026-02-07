package universe

import (
	"encoding/binary"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math/rand"
	"net"
	"networktrafficart/networktrafficart/capture"
	"networktrafficart/networktrafficart/universe/particle"
	"networktrafficart/networktrafficart/util"
)

const (
	ipv4Range = 4294967295.0
)

type Universe struct {
	Particles         []*particle.Particle
	OffscreenDistance float32
}

func NewUniverse() *Universe {
	return &Universe{
		Particles:         []*particle.Particle{},
		OffscreenDistance: 20,
	}
}

func (u *Universe) Tick() {
	for i := len(u.Particles) - 1; i >= 0; i-- {
		p := u.Particles[i]
		p.Y -= p.YDelta
		p.X += p.XSkew
		if p.Y < -u.OffscreenDistance {
			copy(u.Particles[i:], u.Particles[i+1:])
			u.Particles = u.Particles[:len(u.Particles)-1]
		}
	}
}

func (u *Universe) DrawParticles(screen *ebiten.Image, circle *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	for _, p := range u.Particles {
		opts.GeoM.Reset()
		opts.ColorScale.Reset()

		s := float64(p.Size / 50)

		opts.GeoM.Scale(s, s)
		opts.GeoM.Translate(float64(p.X), float64(p.Y))
		opts.ColorScale.ScaleWithColor(p.Color)

		screen.DrawImage(circle, opts)
	}
}

func (u *Universe) AddToParticles(p *particle.Particle) {
	u.Particles = append(u.Particles, p)
}

func (u *Universe) CreateParticle(pe capture.PacketEvent, screenWidth, screenHeight int16) *particle.Particle {
	ip := binary.BigEndian.Uint32(pe.SrcIP)
	xStart := (float32(ip) / ipv4Range) * float32(screenWidth) // TODO fix for ipv6
	xSkewIntensity := float32(util.ClampValue(.4, 0.0, 1.0))

	p := particle.NewParticle(
		xStart,
		float32(screenHeight)+u.OffscreenDistance,
		5,
		(rand.Float32()-.5)*xSkewIntensity,
		ipv4ToRGBA(pe.SrcIP),
		float32(pe.Size)/75,
	)

	return p
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
