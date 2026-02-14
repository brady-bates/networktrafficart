package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"image/color"
	"log"
	"math"
	_map "networktrafficart/map"
)

const (
	screenWidth  = 1920
	screenHeight = 1080
)

var (
	white = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	blue  = color.RGBA{R: 0, G: 89, B: 179, A: 255}

	landmassColor = color.RGBA{R: 128, G: 192, B: 128, A: 255} // green
	borderColor   = color.RGBA{R: 51, G: 102, B: 51, A: 255}   // darker green
)

type Game struct {
	features      []*geojson.Feature
	mapBounds     orb.Bound
	mapProjection *ebiten.Image
}

func main() {
	features, mapBounds := _map.LoadGeoJSON("assets/map/map.geojson")
	g := NewGame(features, mapBounds)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}

func DrawMapFromFeatures(bounds orb.Bound, features []*geojson.Feature, screenBuffer *ebiten.Image) *ebiten.Image {
	for _, f := range features {
		switch geom := f.Geometry.(type) {
		case orb.Polygon:
			DrawCountry(bounds, screenBuffer, geom)
		case orb.MultiPolygon:
			for _, poly := range geom {
				DrawCountry(bounds, screenBuffer, poly)
			}
		}
	}
	return screenBuffer
}

func NewGame(features []*geojson.Feature, mapBounds orb.Bound) *Game {
	return &Game{
		features:      features,
		mapBounds:     mapBounds,
		mapProjection: nil,
	}
}

// Mercator projection
func Project(bounds orb.Bound, lon, lat float64) (float32, float32) {
	minLon := bounds.Min.Lon()
	minLat := bounds.Min.Lat()
	maxLon := bounds.Max.Lon()
	maxLat := bounds.Max.Lat()

	translateToMercator := func(latDeg float64) float64 {
		r := latDeg * math.Pi / 180
		return math.Log(math.Tan(math.Pi/4 + r/2))
	}

	mercMin := translateToMercator(minLat)
	mercMax := translateToMercator(maxLat)
	mercN := translateToMercator(lat)

	padding := 40.0
	w := float64(screenWidth) - padding*2
	h := float64(screenHeight) - padding*2

	x := padding + (lon-minLon)/(maxLon-minLon)*w
	y := padding + (mercMax-mercN)/(mercMax-mercMin)*h

	return float32(x), float32(y)
}

func DrawCountry(bounds orb.Bound, screen *ebiten.Image, poly orb.Polygon) {
	for _, ring := range poly {
		// skip undrawable polygons
		if len(ring) < 2 {
			continue
		}

		var path vector.Path

		x0, y0 := Project(bounds, ring[0].Lon(), ring[0].Lat())
		path.MoveTo(x0, y0)

		for _, pt := range ring[1:] {
			x, y := Project(bounds, pt.Lon(), pt.Lat())
			path.LineTo(x, y)
		}
		path.Close()

		fillDraw := &vector.DrawPathOptions{
			AntiAlias: true,
		}
		fillDraw.ColorScale.ScaleWithColor(landmassColor)
		vector.FillPath(screen, &path, &vector.FillOptions{}, fillDraw)

		strokeDraw := &vector.DrawPathOptions{
			AntiAlias: true,
		}
		strokeDraw.ColorScale.ScaleWithColor(borderColor)
		vector.StrokePath(screen, &path, &vector.StrokeOptions{
			Width: 1.0,
		}, strokeDraw)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.mapProjection == nil {
		g.initMapProjection()
	}

	screen.Fill(blue)
	screen.DrawImage(g.mapProjection, nil)
}

func (g *Game) Update() error { return nil }

func (g *Game) Layout(w, h int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) initMapProjection() {
	g.mapProjection = DrawMapFromFeatures(g.mapBounds, g.features, ebiten.NewImage(screenWidth, screenHeight))
}

// TODO project US on its own and larger due to having a large portion of the traffic?
