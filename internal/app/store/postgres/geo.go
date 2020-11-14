package postgres

import (
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

const (
	earthRadiusM = 6371000
)

func searchCellUnion(lon, lat float64, radius, storageLevel int) s2.CellUnion {
	latlng := s2.LatLngFromDegrees(lat, lon)
	centerPoint := s2.PointFromLatLng(latlng)

	centerAngle := float64(radius) / earthRadiusM
	cap := s2.CapFromCenterAngle(centerPoint, s1.Angle(centerAngle))

	rc := s2.RegionCoverer{MaxLevel: storageLevel, MinLevel: storageLevel}
	return rc.Covering(cap)
}
