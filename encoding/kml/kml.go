// Package kml 实现 KML 的编码.
package kml

import (
	"github.com/chengxiaoer/geomGo"
	"github.com/chengxiaoer/go-kml"
)

// Encode函数  编码任意几何图形.
func Encode(g geom.T) (kml.Element, error) {
	switch g := g.(type) {
	case *geom.Point:
		return EncodePoint(g), nil
	case *geom.LineString:
		return EncodeLineString(g), nil
	case *geom.LinearRing:
		return EncodeLinearRing(g), nil
	case *geom.MultiLineString:
		return EncodeMultiLineString(g), nil
	case *geom.MultiPoint:
		return EncodeMultiPoint(g), nil
	case *geom.MultiPolygon:
		return EncodeMultiPolygon(g), nil
	case *geom.Polygon:
		return EncodePolygon(g), nil
	case *geom.GeometryCollection:
		return EncodeGeometryCollection(g)
	default:
		return nil, geom.ErrUnsupportedType{Value: g}
	}
}

// EncodeLineString函数  编码一个 LineString.
func EncodeLineString(ls *geom.LineString) kml.Element {
	flatCoords := ls.FlatCoords()
	return kml.LineString(kml.CoordinatesFlat(flatCoords, 0, len(flatCoords), ls.Stride(), dim(ls.Layout())))
}

// EncodeLinearRing函数  编码一个 LinearRing.
func EncodeLinearRing(lr *geom.LinearRing) kml.Element {
	flatCoords := lr.FlatCoords()
	return kml.LinearRing(kml.CoordinatesFlat(flatCoords, 0, len(flatCoords), lr.Stride(), dim(lr.Layout())))
}

// EncodeMultiLineString 编码一个MultiLineString.
func EncodeMultiLineString(mls *geom.MultiLineString) kml.Element {
	lineStrings := make([]kml.Element, mls.NumLineStrings())
	flatCoords := mls.FlatCoords()
	ends := mls.Ends()
	stride := mls.Stride()
	d := dim(mls.Layout())
	offset := 0
	for i, end := range ends {
		lineStrings[i] = kml.LineString(kml.CoordinatesFlat(flatCoords, offset, end, stride, d))
		offset = end
	}
	return kml.MultiGeometry(lineStrings...)
}

// EncodeMultiPoint 编码一个MultiPoint.
func EncodeMultiPoint(mp *geom.MultiPoint) kml.Element {
	points := make([]kml.Element, mp.NumPoints())
	flatCoords := mp.FlatCoords()
	stride := mp.Stride()
	d := dim(mp.Layout())
	for i, offset, end := 0, 0, len(flatCoords); offset < end; i++ {
		points[i] = kml.Point(kml.CoordinatesFlat(flatCoords, offset, offset+stride, stride, d))
		offset += stride
	}
	return kml.MultiGeometry(points...)
}

// EncodeMultiPolygon  编码一个MultiPolygon.
func EncodeMultiPolygon(mp *geom.MultiPolygon) kml.Element {
	polygons := make([]kml.Element, mp.NumPolygons())
	flatCoords := mp.FlatCoords()
	endss := mp.Endss()
	stride := mp.Stride()
	d := dim(mp.Layout())
	offset := 0
	for i, ends := range endss {
		boundaries := make([]kml.Element, len(ends))
		for j, end := range ends {
			linearRing := kml.LinearRing(kml.CoordinatesFlat(flatCoords, offset, end, stride, d))
			if j == 0 {
				boundaries[j] = kml.OuterBoundaryIs(linearRing)
			} else {
				boundaries[j] = kml.InnerBoundaryIs(linearRing)
			}
			offset = end
		}
		polygons[i] = kml.Polygon(boundaries...)
	}
	return kml.MultiGeometry(polygons...)
}

// EncodePoint 编码一个 Point.
func EncodePoint(p *geom.Point) kml.Element {
	flatCoords := p.FlatCoords()
	return kml.Point(kml.CoordinatesFlat(flatCoords, 0, len(flatCoords), p.Stride(), dim(p.Layout())))
}

// EncodePolygon 编码一个 Polygon.
func EncodePolygon(p *geom.Polygon) kml.Element {
	boundaries := make([]kml.Element, p.NumLinearRings())
	stride := p.Stride()
	flatCoords := p.FlatCoords()
	d := dim(p.Layout())
	offset := 0
	for i, end := range p.Ends() {
		linearRing := kml.LinearRing(kml.CoordinatesFlat(flatCoords, offset, end, stride, d))
		if i == 0 {
			boundaries[i] = kml.OuterBoundaryIs(linearRing)
		} else {
			boundaries[i] = kml.InnerBoundaryIs(linearRing)
		}
		offset = end
	}
	return kml.Polygon(boundaries...)
}

// EncodeGeometryCollection函数 创建一个GeometryCollection.
func EncodeGeometryCollection(g *geom.GeometryCollection) (kml.Element, error) {
	geometries := make([]kml.Element, g.NumGeoms())
	for i, g := range g.Geoms() {
		var err error
		geometries[i], err = Encode(g)
		if err != nil {
			return nil, err
		}
	}
	return kml.MultiGeometry(geometries...), nil
}

func dim(l geom.Layout) int {
	switch l {
	case geom.XY, geom.XYM:
		return 2
	default:
		return 3
	}
}
