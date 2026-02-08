package universe

import (
	"github.com/hajimehoshi/ebiten/v2"
	"networktrafficart/universe/particle"
	"sync"
)

type Universe struct {
	Particles         []*particle.Particle
	mut               sync.RWMutex
	OffScreenDistance float32
}

func NewUniverse() *Universe {
	return &Universe{
		Particles:         []*particle.Particle{},
		mut:               sync.RWMutex{},
		OffScreenDistance: 25,
	}
}

func (u *Universe) Tick() {
	u.mut.Lock()
	defer u.mut.Unlock()
	u.tickParticles()
}

func (u *Universe) tickParticles() {
	pSlice := u.Particles
	for i := len(pSlice) - 1; i >= 0; i-- {
		p := pSlice[i]
		p.Y -= p.YDelta
		p.X += p.XSkew
		if p.Y < -u.OffScreenDistance {
			copy(pSlice[i:], pSlice[i+1:])
			pSlice = pSlice[:len(pSlice)-1]
		}
	}
}

func (u *Universe) DrawParticles(screen *ebiten.Image, circle *ebiten.Image) {
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

func (u *Universe) AddToParticles(p *particle.Particle) {
	u.mut.Lock()
	defer u.mut.Unlock()
	u.Particles = append(u.Particles, p)
}
