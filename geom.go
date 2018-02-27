// Package geom 定义了有关地理信息有关的几何类型
package geom

//go:生成目标

import (
	"errors"
	"fmt"
	"math"
)

// 一个图层可以被维数来描述。N>4 也是一个有效的图层,维度使用x,y,z,m来描述。m是在经典维度描述中附加
// 的一个值。m可以用描述以时间等其他属性。当前不支持描述一维
type Layout int

const (
	// NoLayout是未知类型
	NoLayout Layout = iota
	// XY是2D图层 (X and Y)
	XY
	// XYZ是3D图层 (X, Y, and Z)
	XYZ
	// XYM是在2D图层的基础上附加一个M值
	XYM
	// XYZM是在3D图层的基础上附加一个M值
	XYZM
)

// ErrLayoutMismatch将会被返回，但图层的几何类型不正确时
// 不能合并
type ErrLayoutMismatch struct {
	Got  Layout
	Want Layout
}

func (e ErrLayoutMismatch) Error() string {
	return fmt.Sprintf("geom: layout mismatch, got %s, want %s", e.Got, e.Want)
}

// ErrStrideMismatch将会被返回当视图的维数与预期不符时
type ErrStrideMismatch struct {
	Got  int
	Want int
}

func (e ErrStrideMismatch) Error() string {
	return fmt.Sprintf("geom: stride mismatch, got %d, want %d", e.Got, e.Want)
}

// 当请求的图层类型不支持时会返回ErrUnsupportedLayout
type ErrUnsupportedLayout Layout

func (e ErrUnsupportedLayout) Error() string {
	return fmt.Sprintf("geom: unsupported layout %s", Layout(e))
}

// 当请求的类型不支持时将会返回 ErrUnsupportedType
type ErrUnsupportedType struct {
	Value interface{}
}

func (e ErrUnsupportedType) Error() string {
	return fmt.Sprintf("geom: unsupported type %T", e.Value)
}

// 一个Coord 表示一个坐标
type Coord []float64  //坐标

/**
*------------------------------
*				Coord（坐标）相关的方法
*---------------------------------
*/
// Clone 深度拷贝Coord
func (c Coord) Clone() Coord {
	return deriveCloneCoord(c)
}

// X 函数返回一个坐标的x坐标.
// Coord的第一个值为x坐标
func (c Coord) X() float64 {
	return c[0]
}

// Y 函数返回一个坐标的y坐标
// Coord的第二个值为y坐标
func (c Coord) Y() float64 {
	return c[1]
}

// Set 函数复制一个坐标点
func (c Coord) Set(other Coord) {
	copy(c, other)
}

// Equal 函数比较点的坐标与其他点的坐标是否相同
// 点的坐标维数必须钧相同
func (c Coord) Equal(layout Layout, other Coord) bool {

	numOrds := len(c)

	if layout.Stride() < numOrds {
		numOrds = layout.Stride()
	}

	if (len(c) < layout.Stride() || len(other) < layout.Stride()) && len(c) != len(other) {
		return false
	}

	for i := 0; i < numOrds; i++ {
		if math.IsNaN(c[i]) || math.IsNaN(other[i]) {
			if !math.IsNaN(c[i]) || !math.IsNaN(other[i]) {
				return false
			}
		} else if c[i] != other[i] {
			return false
		}
	}

	return true
}

// T 是所有几何类型都实现了的泛型接口
type T interface {
	Layout() Layout
	Stride() int
	Bounds() *Bounds
	FlatCoords() []float64
	Ends() []int
	Endss() [][]int
	SRID() int
	TransformXY(func(float64, float64) (float64, float64)) error
	TransformXYZ(func(float64, float64, float64) (float64, float64, float64)) error
}
/**
*------------------------------
*				Layout（视图）相关的方法
*---------------------------------
*/

// MIndex 函数返回视图中附加值M的索引，不存在m附加值时返回-1
func (l Layout) MIndex() int {
	switch l {
	case NoLayout, XY, XYZ:
		return -1
	case XYM:
		return 2
	case XYZM:
		return 3
	default:
		return 3
	}
}

// Stride 函数返回定义的视图的维数
func (l Layout) Stride() int {
	switch l {
	case NoLayout:
		return 0
	case XY:
		return 2
	case XYZ:
		return 3
	case XYM:
		return 3
	case XYZM:
		return 4
	default:
		return int(l)
	}
}

// String 函数返回视图的类型字符串
func (l Layout) String() string {
	switch l {
	case NoLayout:
		return "NoLayout"
	case XY:
		return "XY"
	case XYZ:
		return "XYZ"
	case XYM:
		return "XYM"
	case XYZM:
		return "XYZM"
	default:
		return fmt.Sprintf("Layout(%d)", int(l))
	}
}

// ZIndex 函数返回视图中 z 分量的索引，如果不存在时返回-1
func (l Layout) ZIndex() int {
	switch l {
	case NoLayout, XY, XYM:
		return -1
	default:
		return 2
	}
}

// Must 函数 当err 不为nil时报错，否则返回泛型接口T
func Must(g T, err error) T {
	if err != nil {
		panic(err)
	}
	return g
}

var (
	errIncorrectEnd         = errors.New("geom: incorrect end")
	errLengthStrideMismatch = errors.New("geom: length/stride mismatch")
	errMisalignedEnd        = errors.New("geom: misaligned end")
	errNonEmptyEnds         = errors.New("geom: non-empty ends")
	errNonEmptyEndss        = errors.New("geom: non-empty endss")
	errNonEmptyFlatCoords   = errors.New("geom: non-empty flatCoords")
	errOutOfOrderEnd        = errors.New("geom: out-of-order end")
	errStrideLayoutMismatch = errors.New("geom: stride/layout mismatch")
)
