package simulation

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"networktrafficart/internal/capture"
	"networktrafficart/internal/util"
	"sync"
	"time"
)

type Simulation struct {
	Events            chan capture.PacketData
	Packets           []Packet
	mut               sync.RWMutex
	OffScreenDistance float32
	packetBuffer      chan Packet
}

func NewSimulation(e chan capture.PacketData) *Simulation {
	return &Simulation{
		Events:            e,
		Packets:           []Packet{},
		mut:               sync.RWMutex{},
		OffScreenDistance: 25,
		packetBuffer:      make(chan Packet, 50000),
	}
}

func (s *Simulation) Init(screenWidth, screenHeight, PacketBufferConsumerMaxDelayMicros int, PacketBufferConsumerAggressionCurve float64) {
	go s.WatchEventChannel(
		screenWidth,
		screenHeight,
	)
	go s.CreatePacketsFromBuffer(
		PacketBufferConsumerAggressionCurve,
		PacketBufferConsumerMaxDelayMicros,
	)
}

func (s *Simulation) Tick() {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.tickPackets()
}

func (s *Simulation) tickPackets() {
	var n int
	for _, p := range s.Packets {
		p.Y -= p.YDelta
		p.X += p.XSkew

		if p.Y >= -s.OffScreenDistance {
			s.Packets[n] = p
			n++
		} else {
			s.Packets[n] = Packet{}
		}
	}

	clear(s.Packets[n:])
	s.Packets = s.Packets[:n]
}

func (s *Simulation) DrawPackets(screen *ebiten.Image, circle *ebiten.Image) {
	s.mut.RLock()
	defer s.mut.RUnlock()
	opts := &ebiten.DrawImageOptions{}
	for _, p := range s.Packets {
		opts.GeoM.Reset()
		opts.ColorScale.Reset()

		scale := float64(p.Size / 50)

		opts.GeoM.Scale(scale, scale)
		opts.GeoM.Translate(float64(p.X), float64(p.Y))
		opts.ColorScale.ScaleWithColor(p.Color)

		screen.DrawImage(circle, opts)
	}
}

func (s *Simulation) AddToPackets(p Packet) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.Packets = append(s.Packets, p)
}

func (s *Simulation) WatchEventChannel(screenWidth, screenHeight int) {
	var event capture.PacketData
	for {
		select {
		case event = <-s.Events:
		}

		select {
		case s.packetBuffer <- NewPacketFromEvent(event, screenWidth, screenHeight):
		default:
			fmt.Println("Packet buffer is full")
		}
	}
}

func (s *Simulation) CreatePacketsFromBuffer(aggressionCurve float64, maxWatcherDelay int) {
	curve := util.ClampValue(aggressionCurve, 0.0, math.Inf(+1))
	capacity := float64(cap(s.packetBuffer))
	minDelay := 0.0
	maxDelay := float64(maxWatcherDelay)

	var packet Packet
	for {
		select {
		case packet = <-s.packetBuffer:
			count := float64(len(s.packetBuffer))
			fullness := count / (capacity * .6)
			modulationFactor := math.Pow(fullness, curve)
			modulatedDelay := maxDelay + modulationFactor*(minDelay-maxDelay)
			micro := time.Duration(modulatedDelay) * time.Microsecond

			s.AddToPackets(packet)

			time.Sleep(micro)
		}
	}
}

// TODO find a new visual representation of a packet
// Packet/trail, leaves a faded trail
// Arc of some kind, brief fade in and out - dot traveling along arc?
