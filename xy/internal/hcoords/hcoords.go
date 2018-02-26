package hcoords

import (
	"fmt"
	"math"

	"github.com/chengxiaoer/go-geom"
)

// GetIntersection 计算 齐次坐标下两线段间的（近似）交点。
//
// 注意，该算法的数值不稳定的； i.e. 它可能产生位于线段之外的交点。
//为了提高计算的精度，在将这些点传入函数之前，应该对输入点进行标准化。
func GetIntersection(line1End1, line1End2, line2End1, line2End2 geom.Coord) (geom.Coord, error) {
	// unrolled computation
	line1Xdiff := line1End1[1] - line1End2[1]
	line1Ydiff := line1End2[0] - line1End1[0]
	line1W := line1End1[0]*line1End2[1] - line1End2[0]*line1End1[1]

	line2X := line2End1[1] - line2End2[1]
	line2Y := line2End2[0] - line2End1[0]
	line2W := line2End1[0]*line2End2[1] - line2End2[0]*line2End1[1]

	x := line1Ydiff*line2W - line2Y*line1W
	y := line2X*line1W - line1Xdiff*line2W
	w := line1Xdiff*line2Y - line2X*line1Ydiff

	xIntersection := x / w
	yIntersection := y / w

	if math.IsNaN(xIntersection) || math.IsNaN(yIntersection) {
		return nil, fmt.Errorf("intersection cannot be calculated using the h-coords implementation")
	}

	if math.IsInf(xIntersection, 0) || math.IsInf(yIntersection, 0) {
		return nil, fmt.Errorf("intersection cannot be calculated using the h-coords implementation")
	}

	return geom.Coord{xIntersection, yIntersection}, nil
}
