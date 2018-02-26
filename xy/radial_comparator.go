package xy

import (
	"sort"

	"github.com/chengxiaoer/go-geom"
	"github.com/chengxiaoer/go-geom/bigxy"
	"github.com/chengxiaoer/go-geom/sorting"
	"github.com/chengxiaoer/go-geom/xy/orientation"
)

// NewRadialSorting 创建一个实现的排序.Interface which will sort the wrapped coordinate array
// radially around the focal point.  The comparison is based on the angle and distance
// from the focal point.
// First the angle is checked.
// Counter clockwise indicates a greater value and clockwise indicates a lesser value
// If co-linear then the coordinate nearer to the focalPoint is considered less.
func NewRadialSorting(layout geom.Layout, coordData []float64, focalPoint geom.Coord) sort.Interface {
	isLess := func(v1, v2 []float64) bool {
		orient := bigxy.OrientationIndex(focalPoint, v1, v2)

		if orient == orientation.CounterClockwise {
			return false
		}
		if orient == orientation.Clockwise {
			return true
		}

		dxp := v1[0] - focalPoint[0]
		dyp := v1[1] - focalPoint[1]
		dxq := v2[0] - focalPoint[0]
		dyq := v2[1] - focalPoint[1]

		// 点是共线 - 检查距离
		op := dxp*dxp + dyp*dyp
		oq := dxq*dxq + dyq*dyq
		return op < oq
	}
	return sorting.NewFlatCoordSorting(layout, coordData, isLess)
}
