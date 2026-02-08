package universe

import (
	"github.com/hajimehoshi/ebiten/v2"
	"networktrafficart/universe/particle"
)

type Universe struct {
	Particles         []*particle.Particle
	OffscreenDistance float32
}

func NewUniverse() *Universe {
	return &Universe{
		Particles:         []*particle.Particle{},
		OffscreenDistance: 25,
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
