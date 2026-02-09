package simulation

import (
	"github.com/hajimehoshi/ebiten/v2"
	"sync"
)

type Simulation struct {
	Particles         []*Particle
	mut               sync.RWMutex
	OffScreenDistance float32
}

func NewUniverse() *Simulation {
	return &Simulation{
		Particles:         []*Particle{},
		mut:               sync.RWMutex{},
		OffScreenDistance: 25,
	}
}

func (u *Simulation) TickSimulation() {
	u.mut.Lock()
	defer u.mut.Unlock()
	u.tickParticles()
}

func (u *Simulation) tickParticles() {
	n := 0
	for _, p := range u.Particles {
		p.Y -= p.YDelta
		p.X += p.XSkew

		if p.Y >= -u.OffScreenDistance {
			u.Particles[n] = p
			n++
		} else {
			u.Particles[n] = nil
		}
	}

	clear(u.Particles[n:])
	u.Particles = u.Particles[:n]
}

func (u *Simulation) DrawParticles(screen *ebiten.Image, circle *ebiten.Image) {
	u.mut.RLock()
	defer u.mut.RUnlock()
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

func (u *Simulation) AddToParticles(p *Particle) {
	u.mut.Lock()
	defer u.mut.Unlock()
	u.Particles = append(u.Particles, p)
}
