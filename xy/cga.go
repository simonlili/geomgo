// Package xy 包含了低维平面（xy)相关的函数。
// 数据可以是任意维度，但是每个坐标的前两个坐标必须是x，y坐标。所有其他坐标都将被忽略。
package xy

import (
	"fmt"
	"math"

	"github.com/chengxiaoer/geomGo"
	"github.com/chengxiaoer/geomGo/bigxy"
	"github.com/chengxiaoer/geomGo/xy/internal"
	"github.com/chengxiaoer/geomGo/xy/internal/lineintersector"
	"github.com/chengxiaoer/geomGo/xy/orientation"
)

// OrientationIndex函数 返回一个点关于一个特殊的向量指向的方向的索引
//
// vectorOrigin - 向量的起点
// vectorEnd - 向量的终点
// point - 计算方向的点
func OrientationIndex(vectorOrigin, vectorEnd, point geom.Coord) orientation.Type {
	return bigxy.OrientationIndex(vectorOrigin, vectorEnd, point)
}

// IsOnLine函数 检测一个点是否在由坐标数组构成的线段上，如果这个点是线段端点或者在线段上返回true
func IsOnLine(layout geom.Layout, point geom.Coord, lineSegmentCoordinates []float64) bool {

	stride := layout.Stride()
	if len(lineSegmentCoordinates) < (2 * stride) {
		panic(fmt.Sprintf("At least two coordinates are required in the lineSegmentsCoordinates array in 'algorithms.IsOnLine', was: %v", lineSegmentCoordinates))
	}
	strategy := lineintersector.RobustLineIntersector{}

	for i := stride; i < len(lineSegmentCoordinates); i += stride {
		segmentStart := lineSegmentCoordinates[i-stride : i-stride+2]
		segmentEnd := lineSegmentCoordinates[i : i+2]

		if lineintersector.PointIntersectsLine(strategy, point, geom.Coord(segmentStart), geom.Coord(segmentEnd)) {
			return true
		}
	}
	return false
}

// IsRingCounterClockwise函数 判断是否入的坐标是否能组成一个逆时针的线环。
//
// - 点的列表被假定为第一个和最后一个点相等。
// - 这将会处理坐标列表里面重复的点
// 此算法要求是能构成一个有效的环，无效的环可能导致结果错误.
//
// Param ring - 形成环的坐标数组。
// Returns 如果这个环是逆时针方向排列的，返回true
// Panics 如果传入的点数目小于3，会报错
func IsRingCounterClockwise(layout geom.Layout, ring []float64) bool {
	stride := layout.Stride()

	// # of ordinates without closing endpoint
	nOrds := len(ring) - stride
	// # of points without closing endpoint
	nPts := nOrds / stride
	// 完整性检查
	if nPts < 3 {
		panic("Ring has fewer than 3 points, so orientation cannot be determined")
	}

	// find highest point
	hiIndex := 0
	for i := stride; i <= len(ring)-stride; i += stride {
		if ring[i+1] > ring[hiIndex+1] {
			hiIndex = i
		}
	}

	// find distinct point before highest point
	iPrev := hiIndex
	for {
		iPrev = iPrev - stride
		if iPrev < 0 {
			iPrev = nOrds
		}

		if !internal.Equal(ring, iPrev, ring, hiIndex) || iPrev == hiIndex {
			break
		}
	}

	// find distinct point after highest point
	iNext := hiIndex
	for {
		iNext = (iNext + stride) % nOrds

		if !internal.Equal(ring, iNext, ring, hiIndex) || iNext == hiIndex {
			break
		}
	}

	// This check catches cases where the ring contains an A-B-A configuration
	// of points. This can happen if the ring does not contain 3 distinct points
	// (including the case where the input array has fewer than 4 elements), or
	// it contains coincident line segments.
	if internal.Equal(ring, iPrev, ring, hiIndex) || internal.Equal(ring, iNext, ring, hiIndex) || internal.Equal(ring, iPrev, ring, iNext) {
		return false
	}

	disc := bigxy.OrientationIndex(geom.Coord(ring[iPrev:iPrev+2]), geom.Coord(ring[hiIndex:hiIndex+2]), geom.Coord(ring[iNext:iNext+2]))

	// If disc is exactly 0, lines are collinear. There are two possible cases:
	// (1) the lines lie along the x axis in opposite directions (2) the lines
	// lie on top of one another
	//
	// (1) is handled by checking if next is left of prev ==> CCW (2) will never
	// happen if the ring is valid, so don't check for it (Might want to assert
	// this)
	isCCW := false
	if disc == 0 {
		// poly is CCW if prev x is right of next x
		isCCW = (ring[iPrev] > ring[iNext])
	} else {
		// if area is positive, points are ordered CCW
		isCCW = (disc > 0)
	}
	return isCCW
}

// DistanceFromPointToLine 计算一个点到一个线段的距离
//
// Note: 非强健算法
func DistanceFromPointToLine(p, lineStart, lineEnd geom.Coord) float64 {
	// 如果线段的起点与终点相同，则计算该点到终点的距离
	if lineStart[0] == lineEnd[0] && lineStart[1] == lineEnd[1] {
		return internal.Distance2D(p, lineStart)
	}

	// otherwise use comp.graphics.algorithms Frequently Asked Questions method

	// (1) r = AC dot AB
	//         ---------
	//         ||AB||^2
	//
	// r has the following meaning:
	//   r=0 P = A
	//   r=1 P = B
	//   r<0 P is on the backward extension of AB
	//   r>1 P is on the forward extension of AB
	//   0<r<1 P is interior to AB

	len2 := (lineEnd[0]-lineStart[0])*(lineEnd[0]-lineStart[0]) + (lineEnd[1]-lineStart[1])*(lineEnd[1]-lineStart[1])
	r := ((p[0]-lineStart[0])*(lineEnd[0]-lineStart[0]) + (p[1]-lineStart[1])*(lineEnd[1]-lineStart[1])) / len2

	if r <= 0.0 {
		return internal.Distance2D(p, lineStart)
	}
	if r >= 1.0 {
		return internal.Distance2D(p, lineEnd)
	}

	// (2) s = (Ay-Cy)(Bx-Ax)-(Ax-Cx)(By-Ay)
	//         -----------------------------
	//                    L^2
	//
	// Then the distance from C to P = |s|*L.
	//
	// This is the same calculation as {@link #distancePointLinePerpendicular}.
	// Unrolled here for performance.
	s := ((lineStart[1]-p[1])*(lineEnd[0]-lineStart[0]) - (lineStart[0]-p[0])*(lineEnd[1]-lineStart[1])) / len2
	return math.Abs(s) * math.Sqrt(len2)
}

// PerpendicularDistanceFromPointToLine函数 计算从点P到直线的垂直距离。
// containing the points lineStart/lineEnd
func PerpendicularDistanceFromPointToLine(p, lineStart, lineEnd geom.Coord) float64 {
	// use comp.graphics.algorithms Frequently Asked Questions method
	/*
	 * (2) s = (Ay-Cy)(Bx-Ax)-(Ax-Cx)(By-Ay)
	 *         -----------------------------
	 *                    L^2
	 *
	 * Then the distance from C to P = |s|*L.
	 */
	len2 := (lineEnd[0]-lineStart[0])*(lineEnd[0]-lineStart[0]) + (lineEnd[1]-lineStart[1])*(lineEnd[1]-lineStart[1])
	s := ((lineStart[1]-p[1])*(lineEnd[0]-lineStart[0]) - (lineStart[0]-p[0])*(lineEnd[1]-lineStart[1])) / len2

	distance := math.Abs(s) * math.Sqrt(len2)
	return distance
}

// DistanceFromPointToLineString函数 计算点到线段序列的距离
//
// Param p - 一个点
// Param line - 由顶点定义的连续线段。
func DistanceFromPointToLineString(layout geom.Layout, p geom.Coord, line []float64) float64 {
	if len(line) < 2 {
		panic(fmt.Sprintf("Line array must contain at least one vertex: %v", line))
	}
	// 这处理长度= 1的情况。
	firstPoint := line[0:2]
	minDistance := internal.Distance2D(p, firstPoint)
	stride := layout.Stride()
	for i := 0; i < len(line)-stride; i += stride {
		point1 := geom.Coord(line[i : i+2])
		point2 := geom.Coord(line[i+stride : i+stride+2])
		dist := DistanceFromPointToLine(p, point1, point2)
		if dist < minDistance {
			minDistance = dist
		}
	}
	return minDistance
}

// DistanceFromLineToLine函数 计算两个线段间的距离
//
// Note: 不稳健
//
// param line1Start - 线段1的起点
// param line1End - 线段1的终点（不能与起点相同）
// param line2Start - 线段2的起点
// param line2End - 线段2的终点（不能与起点相同）
func DistanceFromLineToLine(line1Start, line1End, line2Start, line2End geom.Coord) float64 {
	// 检查线段长度是否为0
	if line1Start.Equal(geom.XY, line1End) {
		return DistanceFromPointToLine(line1Start, line2Start, line2End)
	}
	if line2Start.Equal(geom.XY, line2End) {
		return DistanceFromPointToLine(line2End, line1Start, line1End)
	}
	// Let AB == line1 where A == line1Start and B == line1End
	// Let CD == line2 where C == line2Start and D == line2End
	//
	// AB and CD are line segments
	// from comp.graphics.algo
	//
	// Solving the above for r and s yields
	//
	//     (Ay-Cy)(Dx-Cx)-(Ax-Cx)(Dy-Cy)
	// r = ----------------------------- (eqn 1)
	//     (Bx-Ax)(Dy-Cy)-(By-Ay)(Dx-Cx)
	//
	//     (Ay-Cy)(Bx-Ax)-(Ax-Cx)(By-Ay)
	// s = ----------------------------- (eqn 2)
	//     (Bx-Ax)(Dy-Cy)-(By-Ay)(Dx-Cx)
	//
	// Let P be the position vector of the
	// intersection point, then
	//   P=A+r(B-A) or
	//   Px=Ax+r(Bx-Ax)
	//   Py=Ay+r(By-Ay)
	// By examining the values of r & s, you can also determine some other limiting
	// conditions:
	//   If 0<=r<=1 & 0<=s<=1, intersection exists
	//      r<0 or r>1 or s<0 or s>1 line segments do not intersect
	//   If the denominator in eqn 1 is zero, AB & CD are parallel
	//   If the numerator in eqn 1 is also zero, AB & CD are collinear.

	noIntersection := false
	if !internal.DoLinesOverlap(line1Start, line1End, line2Start, line2End) {
		noIntersection = true
	} else {
		denom := (line1End[0]-line1Start[0])*(line2End[1]-line2Start[1]) - (line1End[1]-line1Start[1])*(line2End[0]-line2Start[0])

		if denom == 0 {
			noIntersection = true
		} else {
			rNum := (line1Start[1]-line2Start[1])*(line2End[0]-line2Start[0]) - (line1Start[0]-line2Start[0])*(line2End[1]-line2Start[1])
			sNum := (line1Start[1]-line2Start[1])*(line1End[0]-line1Start[0]) - (line1Start[0]-line2Start[0])*(line1End[1]-line1Start[1])

			s := sNum / denom
			r := rNum / denom

			if (r < 0) || (r > 1) || (s < 0) || (s > 1) {
				noIntersection = true
			}
		}
	}
	if noIntersection {
		return internal.Min(
			DistanceFromPointToLine(line1Start, line2Start, line2End),
			DistanceFromPointToLine(line1End, line2Start, line2End),
			DistanceFromPointToLine(line2Start, line1Start, line1End),
			DistanceFromPointToLine(line2End, line1Start, line1End))
	}
	// 线段相交
	return 0.0
}

// SignedArea函数 计算一个线环的面积。 computes the signed area for a ring. The signed area is positive if the
// 如果线环是顺时针旋转的则结果为正数，如果逆时针旋转结果为负数。
//如果环是退化的或平坦的，结果为0
func SignedArea(layout geom.Layout, ring []float64) float64 {
	stride := layout.Stride()
	if len(ring) < 3*stride {
		return 0.0
	}
	sum := 0.0
	// 基于Shoelace公式。
	// http://en.wikipedia.org/wiki/Shoelace_formula
	x0 := ring[0]
	lenMinusOnePoint := len(ring) - stride
	for i := stride; i < lenMinusOnePoint; i += stride {
		x := ring[i] - x0
		y1 := ring[i+stride+1]
		y2 := ring[i-stride+1]
		sum += x * (y2 - y1)
	}
	return sum / 2.0
}

// IsPointWithinLineBounds函数 计算点是否在以lineEndpoint1、lineEndpoint2为线段边界外面
func IsPointWithinLineBounds(p, lineEndpoint1, lineEndpoint2 geom.Coord) bool {
	return internal.IsPointWithinLineBounds(p, lineEndpoint1, lineEndpoint2)
}

// DoLinesOverlap函数 计算的线段的边界是否重叠
func DoLinesOverlap(line1End1, line1End2, line2End1, line2End2 geom.Coord) bool {
	return internal.DoLinesOverlap(line1End1, line1End2, line2End1, line2End2)
}

// Equal函数 检查点start1在坐标数组1中是否与点start2以坐标数组2构成的向量相等.
//只有x和y坐标进行比较，x被假定为第一坐标和y作为第二坐标，
//这是一种实用方法，只在性能很重要时使用，因为它降低了可读性。
func Equal(coords1 []float64, start1 int, coords2 []float64, start2 int) bool {
	return internal.Equal(coords1, start1, coords2, start2)
}

// Distance函数 计算两点间的距离
func Distance(c1, c2 geom.Coord) float64 {
	return internal.Distance2D(c1, c2)
}
