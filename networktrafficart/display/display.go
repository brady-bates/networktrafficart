package display

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"networktrafficart/networktrafficart/capture"
	"networktrafficart/networktrafficart/dotenv"
	"networktrafficart/networktrafficart/universe"
	"networktrafficart/networktrafficart/util"
	"time"
)

var circleImage *ebiten.Image

type Display struct {
	PacketEventIn chan capture.PacketEvent
	Universe      *universe.Universe
	screenWidth   int16
	screenHeight  int16
}

func NewDisplay(pe chan capture.PacketEvent, u *universe.Universe) *Display {
	env := dotenv.GetDotenv()
	circleImage = ebiten.NewImage(100, 100)
	vector.FillCircle(circleImage, 50, 50, 50, color.White, true)
	d := &Display{
		PacketEventIn: pe,
		Universe:      u,
		screenWidth:   800,
		screenHeight:  600,
	}
	go d.WatchPacketEventChannel(env.PacketEventWatcherAggressionCurve, env.PacketEventWatcherMaxDelayMicros)
	return d
}

func (d *Display) Update() error {
	d.Universe.Tick()
	return nil
}

func (d *Display) Draw(screen *ebiten.Image) {
	d.Universe.DrawParticles(screen, circleImage)
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

		p := d.Universe.CreateParticle(packetEvent, d.screenWidth, d.screenHeight)
		d.Universe.AddToParticles(p)

		//fmt.Printf("PacketEvent: %+v\n", packetEvent)
		//fmt.Printf("Particle: %+v\n", p)
		time.Sleep(micro)
	}
}
