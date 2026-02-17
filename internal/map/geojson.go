package _map

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"log"
	"os"
)

type GeoJsonData struct {
	Features  []*geojson.Feature
	MapBounds orb.Bound
}

func LoadGeoJSON(path string) GeoJsonData {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	fc, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		log.Fatal(err)
	}

	bound := orb.Bound{
		Min: orb.Point{180, 90},
		Max: orb.Point{-180, -90},
	}

	for _, f := range fc.Features {
		if f.Geometry == nil {
			continue
		}

		// Skip geometries below -60 Lat
		if f.Geometry.Bound().Max.Lat() < -60 {
			continue
		}

		switch geom := f.Geometry.(type) {
		case orb.Polygon:
			expandBounds(&bound, geom.Bound())
		case orb.MultiPolygon:
			expandBounds(&bound, geom.Bound())
		}
	}

	if bound.Min[1] < -85.05112 {
		bound.Min[1] = -85.05112
	}
	if bound.Max[1] > 85.05112 {
		bound.Max[1] = 85.05112
	}

	return GeoJsonData{fc.Features, bound}
}

func expandBounds(bound *orb.Bound, expandTo orb.Bound) {
	if expandTo.Min.Lon() < bound.Min.Lon() {
		bound.Min[0] = expandTo.Min.Lon()
	}
	if expandTo.Min.Lat() < bound.Min.Lat() {
		bound.Min[1] = expandTo.Min.Lat()
	}
	if expandTo.Max.Lon() > bound.Max.Lon() {
		bound.Max[0] = expandTo.Max.Lon()
	}
	if expandTo.Max.Lat() > bound.Max.Lat() {
		bound.Max[1] = expandTo.Max.Lat()
	}
}
