package _map

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/paulmach/orb"
	"image/color"
	"math"
)

var (
	landmassColor = color.RGBA{R: 25, G: 40, B: 60, A: 255} // dark slate blue
	borderColor   = color.RGBA{R: 40, G: 60, B: 90, A: 255} // slightly lighter, very subtle
	OceanColor    = color.RGBA{R: 8, G: 15, B: 35, A: 255}  // near black navy
)

func DrawMap(data GeoJsonData, screenBuffer *ebiten.Image) *ebiten.Image {
	for _, f := range data.Features {
		switch geom := f.Geometry.(type) {
		case orb.Polygon:
			drawPolygon(data.MapBounds, screenBuffer, geom)
		case orb.MultiPolygon:
			for _, poly := range geom {
				drawPolygon(data.MapBounds, screenBuffer, poly)
			}
		}
	}
	return screenBuffer
}

// Mercator projection
func Project(bounds orb.Bound, lon, lat float64) (float64, float64) {
	minLon := bounds.Min.Lon()
	minLat := bounds.Min.Lat()
	maxLon := bounds.Max.Lon()
	maxLat := bounds.Max.Lat()

	mercMin := translateToMercator(minLat)
	mercMax := translateToMercator(maxLat)
	mercN := translateToMercator(lat)

	padding := 40.0

	w, h := ebiten.WindowSize()
	width := float64(w) - padding*2
	height := float64(h) - padding*2

	x := padding + (lon-minLon)/(maxLon-minLon)*width
	y := padding + (mercMax-mercN)/(mercMax-mercMin)*height

	return x, y
}

func translateToMercator(latDeg float64) float64 {
	r := latDeg * math.Pi / 180
	return math.Log(math.Tan(math.Pi/4 + r/2))
}

func drawPolygon(bounds orb.Bound, screen *ebiten.Image, poly orb.Polygon) {
	for _, ring := range poly {
		// skip undrawable polygons
		if len(ring) < 2 {
			continue
		}

		var path vector.Path

		x0, y0 := Project(bounds, ring[0].Lon(), ring[0].Lat())
		path.MoveTo(float32(x0), float32(y0))

		for _, pt := range ring[1:] {
			x, y := Project(bounds, pt.Lon(), pt.Lat())
			path.LineTo(float32(x), float32(y))
		}
		path.Close()

		fillOpts := &vector.DrawPathOptions{
			AntiAlias: true,
		}
		fillOpts.ColorScale.ScaleWithColor(landmassColor)
		vector.FillPath(screen, &path, &vector.FillOptions{}, fillOpts)

		borderOpts := &vector.DrawPathOptions{
			AntiAlias: true,
		}
		borderOpts.ColorScale.ScaleWithColor(borderColor)
		vector.StrokePath(screen, &path, &vector.StrokeOptions{
			Width: 1.0,
		}, borderOpts)
	}
}
