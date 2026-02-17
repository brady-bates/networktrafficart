package display

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"net"
	"networktrafficart/internal/geo"
	_map "networktrafficart/internal/map"
	"networktrafficart/internal/simulation"
	"networktrafficart/internal/util"
)

const (
	sw, sh = 1920, 1080
)

type Display struct {
	Simulation      *simulation.Simulation
	ScreenWidth     int
	ScreenHeight    int
	baseCircleImage *ebiten.Image
	screenBuffer    *ebiten.Image
	geoJsonData     _map.GeoJsonData
	geoService      geo.GeoService
	mapProjection   *ebiten.Image
}

func NewDisplay(s *simulation.Simulation, geoData _map.GeoJsonData, geoService geo.GeoService) *Display {
	circleImage := ebiten.NewImage(100, 100)
	vector.FillCircle(circleImage, 50, 50, 50, color.White, true)

	return &Display{
		Simulation:      s,
		ScreenWidth:     sw,
		ScreenHeight:    sh,
		baseCircleImage: circleImage,
		screenBuffer:    ebiten.NewImage(sw, sh),
		geoJsonData:     geoData,
		geoService:      geoService,
		mapProjection:   nil,
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
	if d.mapProjection == nil {
		d.mapProjection = _map.DrawMap(d.geoJsonData.MapBounds, d.geoJsonData.Features, ebiten.NewImage(d.ScreenWidth, d.ScreenHeight))
	}

	ip := net.ParseIP("***REMOVED***")
	city, err := d.geoService.GetCityFromIP(ip)
	if err != nil {
		log.Fatal(err)
	}

	circleImg := ebiten.NewImage(100, 100)
	vector.FillCircle(circleImg, 3, 3, 3, color.White, true)

	x, y := _map.Project(d.geoJsonData.MapBounds, city.Location.Longitude, city.Location.Latitude)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(x, y)

	d.screenBuffer.Fill(_map.OceanColor)
	d.screenBuffer.DrawImage(d.mapProjection, nil)

	screen.DrawImage(d.screenBuffer, nil)
	screen.DrawImage(circleImg, opts)
}

func (d *Display) Layout(w, h int) (int, int) {
	return d.ScreenWidth, d.ScreenHeight
}
