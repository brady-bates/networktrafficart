package display

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"networktrafficart/capture/packetevent"
	"networktrafficart/universe"
	"networktrafficart/universe/particle"
	"networktrafficart/util"
	"networktrafficart/util/shutdown"
	"time"
)

type Display struct {
	PacketEventIn   chan packetevent.PacketEvent
	Universe        *universe.Universe
	screenWidth     int16
	screenHeight    int16
	baseCircleImage *ebiten.Image
}

func NewDisplay(pe chan packetevent.PacketEvent, u *universe.Universe, curve float64, delayMicros int) *Display {
	circleImage := ebiten.NewImage(100, 100)
	vector.FillCircle(circleImage, 50, 50, 50, color.White, true)
	d := &Display{
		PacketEventIn:   pe,
		Universe:        u,
		screenWidth:     800,
		screenHeight:    600,
		baseCircleImage: circleImage,
	}
	go d.WatchPacketEventChannel(curve, delayMicros)
	return d
}

func (d *Display) Update() error {
	if ebiten.IsWindowBeingClosed() {
		ebiten.SetWindowClosingHandled(true)
		shutdown.GetShutDownCtx().Cancel()

		return ebiten.Termination
	}

	d.Universe.Tick()
	return nil
}

func (d *Display) Draw(screen *ebiten.Image) {
	d.Universe.DrawParticles(screen, d.baseCircleImage)
}

func (d *Display) Layout(w, h int) (int, int) {
	d.screenWidth = int16(w)
	d.screenHeight = int16(h)
	return w, h
}

// WatchPacketEventChannel
// Pulls out of channel and adds to the displays universe
func (d *Display) WatchPacketEventChannel(aggressionCurve float64, maxWatcherDelay int) {
	curve := util.ClampValue(aggressionCurve, 0.0, math.Inf(+1))
	capacity := float64(cap(d.PacketEventIn))

	minDelay := 0.0
	maxDelay := float64(maxWatcherDelay)

	vals := struct{ Cap, Curve, Min, Max float64 }{capacity, curve, minDelay, maxDelay}
	fmt.Printf("WatchPacketEventChannel init values: %+v\n", vals)

	for packetEvent := range d.PacketEventIn {
		dlen := float64(len(d.PacketEventIn))

		fullness := dlen / capacity
		mod := math.Pow(fullness, curve)

		modulatedDelay := maxDelay + mod*(minDelay-maxDelay)
		micro := time.Duration(modulatedDelay) * time.Microsecond

		p := particle.CreateParticle(packetEvent, d.screenWidth, d.screenHeight, d.Universe.OffscreenDistance)
		d.Universe.AddToParticles(p) // TODO batch add?

		//fmt.Printf("PacketEvent: %+v\n", packetEvent)
		//fmt.Printf("Particle: %+v\n", p)
		time.Sleep(micro)
	}
}
