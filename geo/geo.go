package geo

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
)

type GeoService struct {
	geoipDB *geoip2.Reader
}

func NewGeoService(geoIPDBPath string) GeoService {
	return GeoService{
		geoipDB: openGeoIPDB(geoIPDBPath),
	}
}

func openGeoIPDB(dbpath string) *geoip2.Reader {
	geodb, err := geoip2.Open(dbpath)
	if err != nil {
		log.Fatal(err)
	}
	return geodb
}

func (g *GeoService) GetCityFromIP(ip net.IP) (*geoip2.City, error) {
	city, err := g.geoipDB.City(ip)
	if err != nil {
		return nil, err
	}

	if city == nil {
		return nil, fmt.Errorf("city not found for %s", ip.String())
	}

	return city, nil
}
