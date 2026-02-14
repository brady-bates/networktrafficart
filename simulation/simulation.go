package simulation

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"networktrafficart/capture"
	"networktrafficart/util"
	"sync"
	"time"
)

type Simulation struct {
	EventIn           chan capture.Event
	Particles         []Particle
	mut               sync.RWMutex
	OffScreenDistance float32
	particleBuffer    chan Particle
}

func NewSimulation(e chan capture.Event) *Simulation {
	size := 50000
	return &Simulation{
		EventIn:           e,
		Particles:         []Particle{},
		mut:               sync.RWMutex{},
		OffScreenDistance: 25,
		particleBuffer:    make(chan Particle, size),
	}
}

// TODO add background noise in the lulls, make it low key and natural
func (s *Simulation) Init(screenWidth, screenHeight, ParticleBufferConsumerMaxDelayMicros int, ParticleBufferConsumerAggressionCurve float64) {
	go s.WatchEventChannel(
		screenWidth,
		screenHeight,
	)
	go s.CreateParticlesFromBuffer(
		ParticleBufferConsumerAggressionCurve,
		ParticleBufferConsumerMaxDelayMicros,
	)
}

func (s *Simulation) Tick() {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.tickParticles()
}

func (s *Simulation) tickParticles() {
	var n int
	for _, p := range s.Particles {
		p.Y -= p.YDelta
		p.X += p.XSkew

		if p.Y >= -s.OffScreenDistance {
			s.Particles[n] = p
			n++
		} else {
			s.Particles[n] = Particle{}
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

		scale := float64(p.Size / 50)

		opts.GeoM.Scale(scale, scale)
		opts.GeoM.Translate(float64(p.X), float64(p.Y))
		opts.ColorScale.ScaleWithColor(p.Color)

		screen.DrawImage(circle, opts)
	}
}

func (s *Simulation) AddToParticles(p Particle) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.Particles = append(s.Particles, p)
}

func (s *Simulation) WatchEventChannel(screenWidth, screenHeight int) {
	var event capture.Event
	for {
		select {
		case event = <-s.EventIn:
		}

		select {
		case s.particleBuffer <- NewParticleFromEvent(event, screenWidth, screenHeight):
		default:
			fmt.Println("Particle buffer is full")
		}
	}
}

func (s *Simulation) CreateParticlesFromBuffer(aggressionCurve float64, maxWatcherDelay int) {
	curve := util.ClampValue(aggressionCurve, 0.0, math.Inf(+1))
	capacity := float64(cap(s.particleBuffer))
	minDelay := 0.0
	maxDelay := float64(maxWatcherDelay)

	var particle Particle
	for {
		select {
		case particle = <-s.particleBuffer:
			count := float64(len(s.particleBuffer))
			fullness := count / (capacity * .6)
			modulationFactor := math.Pow(fullness, curve)
			modulatedDelay := maxDelay + modulationFactor*(minDelay-maxDelay)
			micro := time.Duration(modulatedDelay) * time.Microsecond

			s.AddToParticles(particle)

			time.Sleep(micro)
		}
	}
}
