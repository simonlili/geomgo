# geomGo

[![Build Status](https://travis-ci.org/chengxiaoer/geomGo.svg?branch=master)](https://travis-ci.org/chengxiaoer/geomGo)
[![GoDoc](https://godoc.org/github.com/chengxiaoer/geomGo?status.svg)](https://godoc.org/github.com/chengxiaoer/geomGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/chengxiaoer/geomGo)](https://goreportcard.com/badge/github.com/chengxiaoer/geomGo)

Package geom implements efficient geometry types for geospatial applications.

## Key features

 * OpenGeo Consortium-style geometries.
 * Support for 2D and 3D geometries, measures (time and/or distance), and
   unlimited extra dimensions.
 * Encoding and decoding of common geometry formats (GeoJSON, KML, WKB, and
   others) including [`sql.Scanner`](https://godoc.org/database/sql#Scanner)
   and [`driver.Value`](https://godoc.org/database/sql/driver#Value) interface
   implementations for easy database integration.
 * [2D](https://godoc.org/github.com/chengxiaoer/geomGo/xy) and
   [3D](https://godoc.org/github.com/chengxiaoer/geomGo/xyz) topology functions.
 * Efficient, cache-friendly [internal representation](INTERNALS.md).

## Detailed features

### Geometry types

 * [Point](https://godoc.org/github.com/chengxiaoer/geomGo#Point)
 * [LineString](https://godoc.org/github.com/chengxiaoer/geomGo#LineString)
 * [Polygon](https://godoc.org/github.com/chengxiaoer/geomGo#Polygon)
 * [MultiPoint](https://godoc.org/github.com/chengxiaoer/geomGo#MultiPoint)
 * [MultiLineString](https://godoc.org/github.com/chengxiaoer/geomGo#MultiLineString)
 * [MultiPolygon](https://godoc.org/github.com/chengxiaoer/geomGo#MultiPolygon)
 * [GeometryCollection](https://godoc.org/github.com/chengxiaoer/geomGo#GeometryCollection)

### Encoding and decoding

 * [GeoJSON](https://godoc.org/github.com/chengxiaoer/geomGo/encoding/geojson)
 * [IGC](https://godoc.org/github.com/chengxiaoer/geomGo/encoding/igc) (decoding only)
 * [KML](https://godoc.org/github.com/chengxiaoer/geomGo/encoding/kml) (encoding only)
 * [WKB](https://godoc.org/github.com/chengxiaoer/geomGo/encoding/wkb)
 * [EWKB](https://godoc.org/github.com/chengxiaoer/geomGo/encoding/ewkb)
 * [WKT](https://godoc.org/github.com/chengxiaoer/geomGo/encoding/wkt) (encoding only)
 * [WKB Hex](https://godoc.org/github.com/chengxiaoer/geomGo/encoding/wkbhex)
 * [EWKB Hex](https://godoc.org/github.com/chengxiaoer/geomGo/encoding/ewkbhex)

### Geometry functions

 * [XY](https://godoc.org/github.com/chengxiaoer/geomGo/xy) 2D geometry functions
 * [XYZ](https://godoc.org/github.com/chengxiaoer/geomGo/xyz) 3D geometry functions

## Related libraries

 * [github.com/chengxiaoer/go-gpx](https://github.com/chengxiaoer/go-gpx) GPX encoding and decoding
 * [github.com/chengxiaoer/go-kml](https://github.com/chengxiaoer/go-kml) KML encoding
 * [github.com/chengxiaoer/go-polyline](https://github.com/chengxiaoer/go-polyline) Google Maps Polyline encoding and decoding
 * [github.com/chengxiaoer/go-vali](https://github.com/chengxiaoer/go-vali) IGC validation

[License](LICENSE)
