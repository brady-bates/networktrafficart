package _map

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"log"
	"os"
)

func LoadGeoJSON(path string) ([]*geojson.Feature, orb.Bound) {
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
			expandBound(&bound, geom.Bound())
		case orb.MultiPolygon:
			expandBound(&bound, geom.Bound())
		}
	}

	if bound.Min[1] < -85.05112 {
		bound.Min[1] = -85.05112
	}
	if bound.Max[1] > 85.05112 {
		bound.Max[1] = 85.05112
	}

	return fc.Features, bound
}

func expandBound(bound *orb.Bound, b orb.Bound) {
	if b.Min.Lon() < bound.Min.Lon() {
		bound.Min[0] = b.Min.Lon()
	}
	if b.Min.Lat() < bound.Min.Lat() {
		bound.Min[1] = b.Min.Lat()
	}
	if b.Max.Lon() > bound.Max.Lon() {
		bound.Max[0] = b.Max.Lon()
	}
	if b.Max.Lat() > bound.Max.Lat() {
		bound.Max[1] = b.Max.Lat()
	}
}
