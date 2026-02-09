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
}

func NewDisplay(s *simulation.Simulation) *Display {
	circleImage := ebiten.NewImage(100, 100)
	vector.FillCircle(circleImage, 50, 50, 50, color.White, true)
	return &Display{
		Simulation:      s,
		ScreenWidth:     1920,
		ScreenHeight:    1080,
		baseCircleImage: circleImage,
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
	d.Simulation.DrawParticles(screen, d.baseCircleImage)
}

func (d *Display) Layout(w, h int) (int, int) {
	return d.ScreenWidth, d.ScreenHeight
}
