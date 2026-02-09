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
	"time"
)

type Display struct {
	PacketEventIn   chan packetevent.PacketEvent
	Universe        *universe.Universe
	screenWidth     int
	screenHeight    int
	baseCircleImage *ebiten.Image
}

func NewDisplay(pe chan packetevent.PacketEvent, u *universe.Universe) *Display {
	circleImage := ebiten.NewImage(100, 100)
	vector.FillCircle(circleImage, 50, 50, 50, color.White, true)
	return &Display{
		PacketEventIn:   pe,
		Universe:        u,
		screenWidth:     1920,
		screenHeight:    1080,
		baseCircleImage: circleImage,
	}
}

func (d *Display) Update() error {
	if ebiten.IsWindowBeingClosed() {
		ebiten.SetWindowClosingHandled(true)
		util.GetShutDownCtx().Cancel()

		return ebiten.Termination
	}

	d.Universe.Tick()
	return nil
}

func (d *Display) Draw(screen *ebiten.Image) {
	d.Universe.DrawParticles(screen, d.baseCircleImage)
}

func (d *Display) Layout(w, h int) (int, int) {
	return d.screenWidth, d.screenHeight
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

		d.Universe.AddToParticles(particle.CreateFromPacketEvent(packetEvent, d.screenWidth, d.screenHeight))

		time.Sleep(micro)
	}
}
