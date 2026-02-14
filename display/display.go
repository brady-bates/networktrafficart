package display

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"networktrafficart/simulation"
	"networktrafficart/util"
)

const (
	sw, sh = 1920, 1080
)

var backgroundShader = []byte(fmt.Sprintf(`
package main

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
    y := position.y / %d
    topColor := vec3(0.01, 0.06, 0.12)
    bottomColor := vec3(0.12, 0.05, 0.01)
    t := y * y * (3.0 - 2.0 * y)
    finalRGB := mix(topColor, bottomColor, t)

    return vec4(finalRGB, 1.0)
}
`, sh))

type Display struct {
	Simulation       *simulation.Simulation
	ScreenWidth      int
	ScreenHeight     int
	baseCircleImage  *ebiten.Image
	screenBuffer     *ebiten.Image
	backgroundShader *ebiten.Shader
}

func NewDisplay(s *simulation.Simulation) *Display {
	circleImage := ebiten.NewImage(100, 100)
	vector.FillCircle(circleImage, 50, 50, 50, color.White, true)

	shader, err := ebiten.NewShader(backgroundShader)
	if err != nil {
		log.Fatalln(err)
	}
	return &Display{
		Simulation:       s,
		ScreenWidth:      sw,
		ScreenHeight:     sh,
		baseCircleImage:  circleImage,
		screenBuffer:     ebiten.NewImage(sw, sh),
		backgroundShader: shader,
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
	d.screenBuffer.DrawRectShader(screen.Bounds().Dx(), screen.Bounds().Dy(), d.backgroundShader, nil)
	d.Simulation.DrawParticles(d.screenBuffer, d.baseCircleImage)
	screen.DrawImage(d.screenBuffer, nil)
}

func (d *Display) Layout(w, h int) (int, int) {
	return d.ScreenWidth, d.ScreenHeight
}
