package geom

import (
	"math"
)

// A Bounds 表示多维边界框
type Bounds struct {
	layout Layout
	min    Coord
	max    Coord
}

// NewBounds函数 创建一个边界
func NewBounds(layout Layout) *Bounds {
	stride := layout.Stride()
	min, max := make(Coord, stride), make(Coord, stride)
	for i := 0; i < stride; i++ {
		min[i], max[i] = math.Inf(1), math.Inf(-1)
	}
	return &Bounds{
		layout: layout,
		min:    min,
		max:    max,
	}
}

/**
*------------------------------
*				Bounds（边界）相关的方法
*---------------------------------
*/

// Clone方法 深度拷贝一个Bounds
func (b *Bounds) Clone() *Bounds {
	return deriveCloneBounds(b)
}

// Extend extends b to include geometry g.
func (b *Bounds) Extend(g T) *Bounds {
	b.extendLayout(g.Layout())
	if b.layout == XYZM && g.Layout() == XYM {
		return b.extendXYZMFlatCoordsWithXYM(g.FlatCoords(), 0, len(g.FlatCoords()))
	}
	return b.extendFlatCoords(g.FlatCoords(), 0, len(g.FlatCoords()), g.Stride())
}

// IsEmpty方法 检测边界是否为空
func (b *Bounds) IsEmpty() bool {
	for i, stride := 0, b.layout.Stride(); i < stride; i++ {
		if b.max[i] < b.min[i] {
			return true
		}
	}
	return false
}

// Layout方法 返回边界的视图布局
func (b *Bounds) Layout() Layout {
	return b.layout
}

// Max方法 获取维数中的最大值
func (b *Bounds) Max(dim int) float64 {
	return b.max[dim]
}

// Min方法 获取维数中的最小值
func (b *Bounds) Min(dim int) float64 {
	return b.min[dim]
}

// Overlaps方法 检测本对象是否覆盖传入的边界
func (b *Bounds) Overlaps(layout Layout, b2 *Bounds) bool {
	for i, stride := 0, layout.Stride(); i < stride; i++ {
		if b.min[i] > b2.max[i] || b.max[i] < b2.min[i] {
			return false
		}
	}
	return true
}

// Set方法 设置最小值和最大值.参数必须是一个偶数值
//第一部分为最小值
//第二部分为最大值
func (b *Bounds) Set(args ...float64) *Bounds {
	if len(args)&1 != 0 {
		panic("geom: even number of arguments required")
	}
	stride := len(args) / 2
	b.extendStride(stride)
	for i := 0; i < stride; i++ {
		b.min[i], b.max[i] = args[i], args[i+stride]
	}
	return b
}

// SetCoords方法 设置边界的最大坐标范围，与最小坐标范围
func (b *Bounds) SetCoords(min, max Coord) *Bounds {
	b.min = Coord(make([]float64, b.layout.Stride()))
	b.max = Coord(make([]float64, b.layout.Stride()))
	for i := 0; i < b.layout.Stride(); i++ {
		b.min[i] = math.Min(min[i], max[i])
		b.max[i] = math.Max(min[i], max[i])
	}
	return b
}

// OverlapsPoint方法 点是否在边界框上（点在边界的边界内或边界上）
func (b *Bounds) OverlapsPoint(layout Layout, point Coord) bool {
	for i, stride := 0, layout.Stride(); i < stride; i++ {
		if b.min[i] > point[i] || b.max[i] < point[i] {
			return false
		}
	}
	return true
}

func (b *Bounds) extendFlatCoords(flatCoords []float64, offset, end, stride int) *Bounds {
	b.extendStride(stride)
	for i := offset; i < end; i += stride {
		for j := 0; j < stride; j++ {
			b.min[j] = math.Min(b.min[j], flatCoords[i+j])
			b.max[j] = math.Max(b.max[j], flatCoords[i+j])
		}
	}
	return b
}

func (b *Bounds) extendLayout(layout Layout) {
	switch {
	case b.layout == XYZ && layout == XYM:
		b.min = append(b.min, math.Inf(1))
		b.max = append(b.max, math.Inf(-1))
		b.layout = XYZM
	case b.layout == XYM && (layout == XYZ || layout == XYZM):
		b.min = append(b.min[:2], math.Inf(1), b.min[2])
		b.max = append(b.max[:2], math.Inf(-1), b.max[2])
		b.layout = XYZM
	case b.layout < layout:
		b.extendStride(layout.Stride())
		b.layout = layout
	}
}

func (b *Bounds) extendStride(stride int) {
	for s := b.layout.Stride(); s < stride; s++ {
		b.min = append(b.min, math.Inf(1))
		b.max = append(b.max, math.Inf(-1))
	}
}

func (b *Bounds) extendXYZMFlatCoordsWithXYM(flatCoords []float64, offset, end int) *Bounds {
	for i := offset; i < end; i += 3 {
		b.min[0] = math.Min(b.min[0], flatCoords[i+0])
		b.max[0] = math.Max(b.max[0], flatCoords[i+0])
		b.min[1] = math.Min(b.min[1], flatCoords[i+1])
		b.max[1] = math.Max(b.max[1], flatCoords[i+1])
		b.min[3] = math.Min(b.min[3], flatCoords[i+2])
		b.max[3] = math.Max(b.max[3], flatCoords[i+2])
	}
	return b
}
