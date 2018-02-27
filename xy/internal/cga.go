package internal

import (
	"math"

	"github.com/chengxiaoer/geomGo"
)

// IsPointWithinLineBounds函数 计算点p是否位于点lineEndpoint1、lineEndpoint2组成的线性边界之外
func IsPointWithinLineBounds(p, lineEndpoint1, lineEndpoint2 geom.Coord) bool {
	minx := math.Min(lineEndpoint1[0], lineEndpoint2[0])
	maxx := math.Max(lineEndpoint1[0], lineEndpoint2[0])
	miny := math.Min(lineEndpoint1[1], lineEndpoint2[1])
	maxy := math.Max(lineEndpoint1[1], lineEndpoint2[1])
	return minx <= p[0] && maxx >= p[0] && miny <= p[1] && maxy >= p[1]
}

// DoLinesOverlap函数 计算两条线的边界是否重叠
func DoLinesOverlap(line1End1, line1End2, line2End1, line2End2 geom.Coord) bool {

	min1x := math.Min(line1End1[0], line1End2[0])
	max1x := math.Max(line1End1[0], line1End2[0])
	min1y := math.Min(line1End1[1], line1End2[1])
	max1y := math.Max(line1End1[1], line1End2[1])

	min2x := math.Min(line2End1[0], line2End2[0])
	max2x := math.Max(line2End1[0], line2End2[0])
	min2y := math.Min(line2End1[1], line2End2[1])
	max2y := math.Max(line2End1[1], line2End2[1])

	if min1x > max2x || max1x < min2x {
		return false
	}
	if min1y > max2y || max1y < min2y {
		return false
	}
	return true
}

// Equal函数 检查 在coords1数组中的第start1个坐标与在coords2数组中第start2个坐标是否相同
// 仅仅 x 和 y 被比较，同时x默认为第一坐标，y为第二个坐标
// 这是一种实用方法，只在性能至关重要时使用，因为它降低了可读性。
func Equal(coords1 []float64, start1 int, coords2 []float64, start2 int) bool {
	if coords1[start1] != coords2[start2] {
		return false
	}

	if coords1[start1+1] != coords2[start2+1] {
		return false
	}

	return true
}

// Distance2D 计算两点在2D视图下的距离
func Distance2D(c1, c2 geom.Coord) float64 {
	dx := c1[0] - c2[0]
	dy := c1[1] - c2[1]

	return math.Sqrt(dx*dx + dy*dy)
}
