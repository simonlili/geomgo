package centralendpoint

import (
	"math"

	"github.com/chengxiaoer/go-geom"
	"github.com/chengxiaoer/go-geom/xy/internal"
)

// GetIntersection 计算 通过线段端点最中间的两个线段的近似交点。
//
// 当线段几乎平行并在端点处相交时，这种方法是有效的。
// 对于一个线段的端点位于另一个线段的内部，这也是正确的。
// 取最中心端点确保计算的交点位于线段的范围中
//
// Also, 通过返回一个输入点，这将可以减少线段碎片。
// Intended to be used as a last resort for  computing ill-conditioned intersection situations which cause other methods to fail.
func GetIntersection(line1End1, line1End2, line2End1, line2End2 geom.Coord) geom.Coord {
	intersector := centralEndpointIntersector{
		line1End1: line1End1,
		line1End2: line1End2,
		line2End1: line2End1,
		line2End2: line2End2}
	intersector.compute()
	return intersector.intersectionPoint
}

type centralEndpointIntersector struct {
	line1End1, line1End2, line2End1, line2End2, intersectionPoint geom.Coord
}

func (intersector *centralEndpointIntersector) compute() {
	pts := [4]geom.Coord{intersector.line1End1, intersector.line1End2, intersector.line2End1, intersector.line2End2}
	centroid := average(pts)
	intersector.intersectionPoint = findNearestPoint(centroid, pts)
}

func average(pts [4]geom.Coord) geom.Coord {
	n := float64(len(pts))
	avg := geom.Coord{0, 0}

	for i := 0; i < len(pts); i++ {
		avg[0] += pts[i][0]
		avg[1] += pts[i][1]
	}
	if n > 0 {
		avg[0] = avg[0] / n
		avg[1] = avg[1] / n
	}
	return avg
}

func findNearestPoint(p geom.Coord, pts [4]geom.Coord) geom.Coord {
	minDist := math.MaxFloat64
	result := geom.Coord{}
	for i := 0; i < len(pts); i++ {
		dist := internal.Distance2D(p, pts[i])

		// 始终初始化结果
		if i == 0 || dist < minDist {
			minDist = dist
			result = pts[i]
		}
	}
	return result
}
