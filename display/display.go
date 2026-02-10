package display

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"networktrafficart/simulation"
	"networktrafficart/util"
)

type Display struct {
	Simulation      *simulation.Simulation
	ScreenWidth     int
	ScreenHeight    int
	baseCircleImage *ebiten.Image
	screenBuffer    *ebiten.Image
}

func NewDisplay(s *simulation.Simulation) *Display {
	sw, sh := 1920, 1080
	circleImage := ebiten.NewImage(100, 100)
	vector.FillCircle(circleImage, 50, 50, 50, color.White, true)
	return &Display{
		Simulation:      s,
		ScreenWidth:     sw,
		ScreenHeight:    sh,
		baseCircleImage: circleImage,
		screenBuffer:    ebiten.NewImage(sw, sh),
	}
}

func (d *Display) Update() error {
	//fmt.Printf("fps: %f tps: %f\n", ebiten.ActualFPS(), ebiten.ActualTPS())

	if ebiten.IsWindowBeingClosed() {
		ebiten.SetWindowClosingHandled(true)
		util.GetShutDownCtx().Cancel()

		return ebiten.Termination
	}

	d.Simulation.Tick()
	return nil
}

func (d *Display) Draw(screen *ebiten.Image) {
	d.screenBuffer.Clear()
	d.Simulation.DrawParticles(d.screenBuffer, d.baseCircleImage)
	screen.DrawImage(d.screenBuffer, nil)
}

func (d *Display) Layout(w, h int) (int, int) {
	return d.ScreenWidth, d.ScreenHeight
}
