# go-geom

[![Build Status](https://travis-ci.org/chengxiaoer/go-geom.svg?branch=master)](https://travis-ci.org/chengxiaoer/go-geom)
[![GoDoc](https://godoc.org/github.com/chengxiaoer/go-geom?status.svg)](https://godoc.org/github.com/chengxiaoer/go-geom)
[![Go Report Card](https://goreportcard.com/badge/github.com/chengxiaoer/go-geom)](https://goreportcard.com/badge/github.com/chengxiaoer/go-geom)

Package geom implements efficient geometry types for geospatial applications.

## Key features

 * OpenGeo Consortium-style geometries.
 * Support for 2D and 3D geometries, measures (time and/or distance), and
   unlimited extra dimensions.
 * Encoding and decoding of common geometry formats (GeoJSON, KML, WKB, and
   others) including [`sql.Scanner`](https://godoc.org/database/sql#Scanner)
   and [`driver.Value`](https://godoc.org/database/sql/driver#Value) interface
   implementations for easy database integration.
 * [2D](https://godoc.org/github.com/chengxiaoer/go-geom/xy) and
   [3D](https://godoc.org/github.com/chengxiaoer/go-geom/xyz) topology functions.
 * Efficient, cache-friendly [internal representation](INTERNALS.md).

## Detailed features

### Geometry types

 * [Point](https://godoc.org/github.com/chengxiaoer/go-geom#Point)
 * [LineString](https://godoc.org/github.com/chengxiaoer/go-geom#LineString)
 * [Polygon](https://godoc.org/github.com/chengxiaoer/go-geom#Polygon)
 * [MultiPoint](https://godoc.org/github.com/chengxiaoer/go-geom#MultiPoint)
 * [MultiLineString](https://godoc.org/github.com/chengxiaoer/go-geom#MultiLineString)
 * [MultiPolygon](https://godoc.org/github.com/chengxiaoer/go-geom#MultiPolygon)
 * [GeometryCollection](https://godoc.org/github.com/chengxiaoer/go-geom#GeometryCollection)

### Encoding and decoding

 * [GeoJSON](https://godoc.org/github.com/chengxiaoer/go-geom/encoding/geojson)
 * [IGC](https://godoc.org/github.com/chengxiaoer/go-geom/encoding/igc) (decoding only)
 * [KML](https://godoc.org/github.com/chengxiaoer/go-geom/encoding/kml) (encoding only)
 * [WKB](https://godoc.org/github.com/chengxiaoer/go-geom/encoding/wkb)
 * [EWKB](https://godoc.org/github.com/chengxiaoer/go-geom/encoding/ewkb)
 * [WKT](https://godoc.org/github.com/chengxiaoer/go-geom/encoding/wkt) (encoding only)
 * [WKB Hex](https://godoc.org/github.com/chengxiaoer/go-geom/encoding/wkbhex)
 * [EWKB Hex](https://godoc.org/github.com/chengxiaoer/go-geom/encoding/ewkbhex)

### Geometry functions

 * [XY](https://godoc.org/github.com/chengxiaoer/go-geom/xy) 2D geometry functions
 * [XYZ](https://godoc.org/github.com/chengxiaoer/go-geom/xyz) 3D geometry functions

## Related libraries

 * [github.com/chengxiaoer/go-gpx](https://github.com/chengxiaoer/go-gpx) GPX encoding and decoding
 * [github.com/chengxiaoer/go-kml](https://github.com/chengxiaoer/go-kml) KML encoding
 * [github.com/chengxiaoer/go-polyline](https://github.com/chengxiaoer/go-polyline) Google Maps Polyline encoding and decoding
 * [github.com/chengxiaoer/go-vali](https://github.com/chengxiaoer/go-vali) IGC validation

[License](LICENSE)
