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

func NewSimulation() *Simulation {
	return &Simulation{
		Particles:         []*Particle{},
		mut:               sync.RWMutex{},
		OffScreenDistance: 25,
	}
}

func (s *Simulation) TickSimulation() {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.tickParticles()
}

func (s *Simulation) tickParticles() {
	n := 0
	for _, p := range s.Particles {
		p.Y -= p.YDelta
		p.X += p.XSkew

		if p.Y >= -s.OffScreenDistance {
			s.Particles[n] = p
			n++
		} else {
			s.Particles[n] = nil
		}
	}

	clear(s.Particles[n:])
	s.Particles = s.Particles[:n]
}

func (s *Simulation) DrawParticles(screen *ebiten.Image, circle *ebiten.Image) {
	s.mut.RLock()
	defer s.mut.RUnlock()
	opts := &ebiten.DrawImageOptions{}
	for _, p := range s.Particles {
		opts.GeoM.Reset()
		opts.ColorScale.Reset()

		s := float64(p.Size / 50)

		opts.GeoM.Scale(s, s)
		opts.GeoM.Translate(float64(p.X), float64(p.Y))
		opts.ColorScale.ScaleWithColor(p.Color)

		screen.DrawImage(circle, opts)
	}
}

func (s *Simulation) AddToParticles(p *Particle) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.Particles = append(s.Particles, p)
}
