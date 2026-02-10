package simulation

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"networktrafficart/capture"
	"networktrafficart/util"
	"sync"
)

type Simulation struct {
	PacketEventIn     chan capture.PacketEvent
	Particles         []*Particle
	mut               sync.RWMutex
	OffScreenDistance float32
}

func NewSimulation(pe chan capture.PacketEvent) *Simulation {
	return &Simulation{
		PacketEventIn:     pe,
		Particles:         []*Particle{},
		mut:               sync.RWMutex{},
		OffScreenDistance: 25,
	}
}

func (s *Simulation) Tick() {
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

func (s *Simulation) WatchPacketEventChannel(aggressionCurve float64, maxWatcherDelay int, screenWidth, screenHeight int) {
	fmt.Println(ebiten.WindowSize())
	curve := util.ClampValue(aggressionCurve, 0.0, math.Inf(+1))
	capacity := float64(cap(s.PacketEventIn))

	minDelay := 0.0
	maxDelay := float64(maxWatcherDelay)

	var packetEvent capture.PacketEvent
	for {
		select {
		case packetEvent = <-s.PacketEventIn:
		}

		dlen := float64(len(s.PacketEventIn))

		fullness := dlen / capacity
		mod := math.Pow(fullness, curve)

		modulatedDelay := maxDelay + mod*(minDelay-maxDelay)
		//micro := time.Duration(modulatedDelay) * time.Microsecond

		fmt.Printf("micros: %f fullness: %.8f mod: %.2f \n", modulatedDelay, fullness, mod)

		p := NewParticle(packetEvent, screenWidth, screenHeight)
		s.AddToParticles(p)

		//time.Sleep(micro)
	}
}
