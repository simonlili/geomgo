package xy

import (
	"fmt"

	"github.com/chengxiaoer/geomGo"
)

// Centroid函数 计算几何体的质心。、
//根据几何学的拓扑结构，质心可能在几何之外。
func Centroid(geometry geom.T) (centroid geom.Coord, err error) {
	switch t := geometry.(type) {
	case *geom.Point:
		centroid = PointsCentroid(t)
	case *geom.MultiPoint:
		centroid = MultiPointCentroid(t)
	case *geom.LineString:
		centroid = LinesCentroid(t)
	case *geom.LinearRing:
		centroid = LinearRingsCentroid(t)
	case *geom.MultiLineString:
		centroid = MultiLineCentroid(t)
	case *geom.Polygon:
		centroid = PolygonsCentroid(t)
	case *geom.MultiPolygon:
		centroid = MultiPolygonCentroid(t)
	default:
		err = fmt.Errorf("%v is not a supported type for centroid calculation", t)
	}

	return centroid, err
}
