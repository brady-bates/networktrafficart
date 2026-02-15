package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"networktrafficart/geo"
	_map "networktrafficart/map"
	"networktrafficart/util"
)

const (
	ScreenWidth  = 1920
	ScreenHeight = 1080
)

var (
	white = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	blue  = color.RGBA{R: 0, G: 89, B: 179, A: 255}
)

type Game struct {
	geoJsonData   _map.GeoJsonData
	mapProjection *ebiten.Image
	geoService    geo.GeoService
}

func main() {
	geoData := _map.LoadGeoJSON("assets/map/map.geojson")
	geoService := geo.NewGeoService("assets/geolitedb/GeoLite2-City.mmdb")
	g := NewGame(geoData, geoService)
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}

func NewGame(geoData _map.GeoJsonData, geoService geo.GeoService) *Game {
	return &Game{
		geoJsonData:   geoData,
		mapProjection: nil,
		geoService:    geoService,
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.mapProjection == nil {
		g.mapProjection = _map.DrawMap(g.geoJsonData.MapBounds, g.geoJsonData.Features, ebiten.NewImage(ScreenWidth, ScreenHeight))
	}

	ip := util.GenerateRandomIPv4()
	city, err := g.geoService.GetCityFromIP(ip)
	if err != nil {
		log.Fatal(err)
	}

	circleImg := ebiten.NewImage(100, 100)
	vector.FillCircle(circleImg, 3, 3, 3, color.White, true)

	x, y := _map.Project(g.geoJsonData.MapBounds, city.Location.Longitude, city.Location.Latitude)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(x, y)

	screen.Fill(blue)
	screen.DrawImage(g.mapProjection, nil)

	screen.DrawImage(circleImg, opts)
	//log.Printf("IP: %s Lat: %.3f Long: %.3f\n X: %.3f Y: %.3f", ip, city.Location.Latitude, city.Location.Longitude, x, y)
}

func (g *Game) Update() error { return nil }

func (g *Game) Layout(w, h int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// TODO project US on its own and larger due to having a large portion of the traffic?
