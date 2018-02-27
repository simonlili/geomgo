package wkb

import (
	"bytes"
	"database/sql/driver"
	"fmt"

	"github.com/chengxiaoer/geomGo"
	"github.com/chengxiaoer/geomGo/encoding/wkbcommon"
)

// ErrExpectedByteSlice函数  将返回当 需要一个 []byte时.
type ErrExpectedByteSlice struct {
	Value interface{}
}

func (e ErrExpectedByteSlice) Error() string {
	return fmt.Sprintf("wkb: want []byte, got %T", e.Value)
}

// A Point is 是一个WKB编码的 Point,实现了 sql.Scanner和 driver.Valuer接口.
type Point struct {
	*geom.Point
}

// A LineString 是一个WKB编码的 LineString,实现了 sql.Scanner和 driver.Valuer接口.
type LineString struct {
	*geom.LineString
}

// A Polygon 是一个WKB编码的 Polygon,实现了 sql.Scanner和 driver.Valuer接口.
type Polygon struct {
	*geom.Polygon
}

// A MultiPoint 是一个WKB编码的 MultiPoint,实现了 sql.Scanner和 driver.Valuer接口.
type MultiPoint struct {
	*geom.MultiPoint
}

// A MultiLineString 是一个WKB编码的 MultiLineString,实现了 sql.Scanner和 driver.Valuer接口.
type MultiLineString struct {
	*geom.MultiLineString
}

// A MultiPolygon 是一个WKB编码的 MultiPolygon,实现了 sql.Scanner和 driver.Valuer接口.
type MultiPolygon struct {
	*geom.MultiPolygon
}

// A GeometryCollection 是一个WKB编码的 GeometryCollection,实现了 sql.Scanner和 driver.Valuer接口.
type GeometryCollection struct {
	*geom.GeometryCollection
}

// Scan方法 从 []byte 中扫描（遍历）.
func (p *Point) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrExpectedByteSlice{Value: src}
	}
	got, err := Unmarshal(b)
	if err != nil {
		return err
	}
	p1, ok := got.(*geom.Point)
	if !ok {
		return wkbcommon.ErrUnexpectedType{Got: p1, Want: p}
	}
	p.Point = p1
	return nil
}

// Value方法 返回 p对象的 WKB 编码
func (p *Point) Value() (driver.Value, error) {
	return value(p.Point)
}

// Scan scans from a []byte.
func (ls *LineString) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrExpectedByteSlice{Value: src}
	}
	got, err := Unmarshal(b)
	if err != nil {
		return err
	}
	ls1, ok := got.(*geom.LineString)
	if !ok {
		return wkbcommon.ErrUnexpectedType{Got: ls1, Want: ls}
	}
	ls.LineString = ls1
	return nil
}

// Value returns the WKB encoding of ls.
func (ls *LineString) Value() (driver.Value, error) {
	return value(ls.LineString)
}

// Scan scans from a []byte.
func (p *Polygon) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrExpectedByteSlice{Value: src}
	}
	got, err := Unmarshal(b)
	if err != nil {
		return err
	}
	p1, ok := got.(*geom.Polygon)
	if !ok {
		return wkbcommon.ErrUnexpectedType{Got: p1, Want: p}
	}
	p.Polygon = p1
	return nil
}

// Value returns the WKB encoding of p.
func (p *Polygon) Value() (driver.Value, error) {
	return value(p.Polygon)
}

// Scan scans from a []byte.
func (mp *MultiPoint) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrExpectedByteSlice{Value: src}
	}
	got, err := Unmarshal(b)
	if err != nil {
		return err
	}
	mp1, ok := got.(*geom.MultiPoint)
	if !ok {
		return wkbcommon.ErrUnexpectedType{Got: mp1, Want: mp}
	}
	mp.MultiPoint = mp1
	return nil
}

// Value returns the WKB encoding of mp.
func (mp *MultiPoint) Value() (driver.Value, error) {
	return value(mp.MultiPoint)
}

// Scan scans from a []byte.
func (mls *MultiLineString) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrExpectedByteSlice{Value: src}
	}
	got, err := Unmarshal(b)
	if err != nil {
		return err
	}
	mls1, ok := got.(*geom.MultiLineString)
	if !ok {
		return wkbcommon.ErrUnexpectedType{Got: mls1, Want: mls}
	}
	mls.MultiLineString = mls1
	return nil
}

// Value returns the WKB encoding of mls.
func (mls *MultiLineString) Value() (driver.Value, error) {
	return value(mls.MultiLineString)
}

// Scan scans from a []byte.
func (mp *MultiPolygon) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrExpectedByteSlice{Value: src}
	}
	got, err := Unmarshal(b)
	if err != nil {
		return err
	}
	mp1, ok := got.(*geom.MultiPolygon)
	if !ok {
		return wkbcommon.ErrUnexpectedType{Got: mp1, Want: mp}
	}
	mp.MultiPolygon = mp1
	return nil
}

// Value returns the WKB encoding of mp.
func (mp *MultiPolygon) Value() (driver.Value, error) {
	return value(mp.MultiPolygon)
}

// Scan scans from a []byte.
func (gc *GeometryCollection) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrExpectedByteSlice{Value: src}
	}
	got, err := Unmarshal(b)
	if err != nil {
		return err
	}
	gc1, ok := got.(*geom.GeometryCollection)
	if !ok {
		return wkbcommon.ErrUnexpectedType{Got: gc1, Want: gc}
	}
	gc.GeometryCollection = gc1
	return nil
}

// Value returns the WKB encoding of gc.
func (gc *GeometryCollection) Value() (driver.Value, error) {
	return value(gc.GeometryCollection)
}

func value(g geom.T) (driver.Value, error) {
	b := &bytes.Buffer{}
	if err := Write(b, NDR, g); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
