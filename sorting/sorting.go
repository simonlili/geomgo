package sorting

import (
	"github.com/chengxiaoer/go-geom"
)

// FlatCoord 是一个排序接口的实现，这将通过排序函数改变坐标的顺序
//
// Note: this data struct cannot be used with its 0 values.  it must be constructed using NewFlatCoordSorting
type FlatCoord struct {
	isLess IsLess
	coords []float64
	layout geom.Layout
	stride int
}

// IsLess 被FlatCoord用于去对坐标数组进行排序
// returns true 如果v1小于 v2
type IsLess func(v1, v2 []float64) bool

// IsLess2D is 一个比较器，比较基于的X和Y坐标的大小。
//
// 首先比较 x 坐标
// 如果 x 坐标相同，然后比较 y 坐标
func IsLess2D(v1, v2 []float64) bool {
	if v1[0] < v2[0] {
		return true
	}
	if v1[0] > v2[0] {
		return false
	}
	if v1[1] < v2[1] {
		return true
	}

	return false
}

// NewFlatCoordSorting2D函数 创建一个基于排序接口实现的 2D比较器。
func NewFlatCoordSorting2D(layout geom.Layout, coordData []float64) FlatCoord {
	return NewFlatCoordSorting(layout, coordData, IsLess2D)
}

// NewFlatCoordSorting函数 创建一个基于比较函数的排序接口实现
func NewFlatCoordSorting(layout geom.Layout, coordData []float64, comparator IsLess) FlatCoord {
	return FlatCoord{
		isLess: comparator,
		coords: coordData,
		layout: layout,
		stride: layout.Stride(),
	}
}

func (s FlatCoord) Len() int {
	return len(s.coords) / s.stride
}
func (s FlatCoord) Swap(i, j int) {
	for k := 0; k < s.stride; k++ {
		s.coords[i*s.stride+k], s.coords[j*s.stride+k] = s.coords[j*s.stride+k], s.coords[i*s.stride+k]
	}
}
func (s FlatCoord) Less(i, j int) bool {
	is, js := i*s.stride, j*s.stride
	return s.isLess(s.coords[is:is+s.stride], s.coords[js:js+s.stride])
}
