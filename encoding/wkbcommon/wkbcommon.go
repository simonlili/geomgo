// Package wkbcommon 包含了 WKB and EWKB 编码相关的公共代码.
package wkbcommon

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Byte order IDs.
const (
	XDRID = 0
	NDRID = 1
)

// Byte orders.
var (
	XDR = binary.BigEndian
	NDR = binary.LittleEndian
)

// An ErrUnknownByteOrder 将返回当一个位置的 byte顺序是非法的时.
type ErrUnknownByteOrder byte

func (e ErrUnknownByteOrder) Error() string {
	return fmt.Sprintf("wkb: unknown byte order: %b", byte(e))
}

// An ErrUnsupportedByteOrder 将返回当遇到不支持的 byte order
type ErrUnsupportedByteOrder struct{}

func (e ErrUnsupportedByteOrder) Error() string {
	return "wkb: unsupported byte order"
}

// A Type is a WKB code.
type Type uint32

// An ErrUnknownType 将返回当遇到未知类型时
type ErrUnknownType Type

func (e ErrUnknownType) Error() string {
	return fmt.Sprintf("wkb: unknown type: %d", uint(e))
}

// An ErrUnsupportedType 将返回当遇到不支持的类型时.
type ErrUnsupportedType Type

func (e ErrUnsupportedType) Error() string {
	return fmt.Sprintf("wkb: unsupported type: %d", uint(e))
}

// An ErrUnexpectedType 将返回当遇到不符合要求的类型时..
type ErrUnexpectedType struct {
	Got  interface{}
	Want interface{}
}

func (e ErrUnexpectedType) Error() string {
	return fmt.Sprintf("wkb: got %T, want %T", e.Got, e.Want)
}

// MaxGeometryElements 是在不同级别解码的元素的最大数目.其主要目的是防止错误的输入造成过度的内存分配。
// (担心被用作拒绝服务攻击。).
// FIXME 这个应当是局部的，不是全局的
// FIXME 考虑每个几何图形的极限，而不是每一级极限。
var MaxGeometryElements = [4]uint32{
	0,
	1 << 20, // 没有 LineString, LinearRing, or MultiPoint 可以包含 超过 1048576个坐标
	1 << 15, // 没有 MultiLineString or Polygon 可以包含超过 32768 个LineStrings or LinearRings
	1 << 10, // 没有 MultiPolygon 可以包含超过 1024 个 Polygons
}

// An ErrGeometryTooLarge 将返回当几何图形过大.
type ErrGeometryTooLarge struct {
	Level int
	N     uint32
	Limit uint32
}

func (e ErrGeometryTooLarge) Error() string {
	return fmt.Sprintf("wkb: number of elements at level %d (%d) exceeds %d", e.Level, e.N, e.Limit)
}

// Geometry type IDs.
const (
	PointID              = 1
	LineStringID         = 2
	PolygonID            = 3
	MultiPointID         = 4
	MultiLineStringID    = 5
	MultiPolygonID       = 6
	GeometryCollectionID = 7
	PolyhedralSurfaceID  = 15
	TINID                = 16
	TriangleID           = 17
)

// ReadFlatCoords0函数 读取平面坐标 0.
func ReadFlatCoords0(r io.Reader, byteOrder binary.ByteOrder, stride int) ([]float64, error) {
	coord := make([]float64, stride)
	if err := ReadFloatArray(r, byteOrder, coord); err != nil {
		return nil, err
	}
	return coord, nil
}

// ReadFlatCoords1函数 读取平面坐标 1.
func ReadFlatCoords1(r io.Reader, byteOrder binary.ByteOrder, stride int) ([]float64, error) {
	n, err := ReadUInt32(r, byteOrder)
	if err != nil {
		return nil, err
	}
	if n > MaxGeometryElements[1] {
		return nil, ErrGeometryTooLarge{Level: 1, N: n, Limit: MaxGeometryElements[1]}
	}
	flatCoords := make([]float64, int(n)*stride)
	if err := ReadFloatArray(r, byteOrder, flatCoords); err != nil {
		return nil, err
	}
	return flatCoords, nil
}

// ReadFlatCoords2函数 读取平面坐标 2.
func ReadFlatCoords2(r io.Reader, byteOrder binary.ByteOrder, stride int) ([]float64, []int, error) {
	n, err := ReadUInt32(r, byteOrder)
	if err != nil {
		return nil, nil, err
	}
	if n > MaxGeometryElements[2] {
		return nil, nil, ErrGeometryTooLarge{Level: 2, N: n, Limit: MaxGeometryElements[2]}
	}
	var flatCoordss []float64
	var ends []int
	for i := 0; i < int(n); i++ {
		flatCoords, err := ReadFlatCoords1(r, byteOrder, stride)
		if err != nil {
			return nil, nil, err
		}
		flatCoordss = append(flatCoordss, flatCoords...)
		ends = append(ends, len(flatCoordss))
	}
	return flatCoordss, ends, nil
}

// WriteFlatCoords0函数 写入平面坐标 0 .
func WriteFlatCoords0(w io.Writer, byteOrder binary.ByteOrder, coord []float64) error {
	return WriteFloatArray(w, byteOrder, coord)
}

// WriteFlatCoords1函数 写入平面坐标 1 .
func WriteFlatCoords1(w io.Writer, byteOrder binary.ByteOrder, coords []float64, stride int) error {
	if err := WriteUInt32(w, byteOrder, uint32(len(coords)/stride)); err != nil {
		return err
	}
	return WriteFloatArray(w, byteOrder, coords)
}

// WriteFlatCoords2函数 写入平面坐标 2 .
func WriteFlatCoords2(w io.Writer, byteOrder binary.ByteOrder, flatCoords []float64, ends []int, stride int) error {
	if err := WriteUInt32(w, byteOrder, uint32(len(ends))); err != nil {
		return err
	}
	offset := 0
	for _, end := range ends {
		if err := WriteFlatCoords1(w, byteOrder, flatCoords[offset:end], stride); err != nil {
			return err
		}
		offset = end
	}
	return nil
}
