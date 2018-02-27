package lineintersector

import (
	"github.com/chengxiaoer/geomGo"
	"github.com/chengxiaoer/geomGo/bigxy"
	"github.com/chengxiaoer/geomGo/xy/internal"
	"github.com/chengxiaoer/geomGo/xy/internal/centralendpoint"
	"github.com/chengxiaoer/geomGo/xy/internal/hcoords"
	"github.com/chengxiaoer/geomGo/xy/lineintersection"
	"github.com/chengxiaoer/geomGo/xy/orientation"
)

// RobustLineIntersector is 是一个不完整的实现相比于非稳健实施，但是
// 在极端情况下提供更一致的结果
type RobustLineIntersector struct {
}

func (intersector RobustLineIntersector) computePointOnLineIntersection(data *lineIntersectorData, point, lineStart, lineEnd geom.Coord) {
	data.isProper = false
	// 先检查一下，因为它比定位测试快。
	if internal.IsPointWithinLineBounds(point, lineStart, lineEnd) {
		if bigxy.OrientationIndex(lineStart, lineEnd, point) == orientation.Collinear && bigxy.OrientationIndex(lineEnd, lineStart, point) == orientation.Collinear {
			data.isProper = true
			if internal.Equal(point, 0, lineStart, 0) || internal.Equal(point, 0, lineEnd, 0) {
				data.isProper = false
			}
			data.intersectionType = lineintersection.PointIntersection
			return
		}
	}
	data.intersectionType = lineintersection.NoIntersection
}

func (intersector RobustLineIntersector) computeLineOnLineIntersection(data *lineIntersectorData, line1Start, line1End, line2Start, line2End geom.Coord) {
	data.isProper = false

	// 首先尝试一个快速测试，检查是否相交。
	if !internal.DoLinesOverlap(line1Start, line1End, line2Start, line2End) {
		data.intersectionType = lineintersection.NoIntersection
		return
	}

	// 对于每个结束端点，计算它所在的其他部分的哪一方。
	// 如果两个结束端点位于另一个段的同一侧，
	// 则线段不相交
	line2StartToLine1Orientation := bigxy.OrientationIndex(line1Start, line1End, line2Start)
	line2EndToLine1Orientation := bigxy.OrientationIndex(line1Start, line1End, line2End)

	if (line2StartToLine1Orientation > orientation.Collinear && line2EndToLine1Orientation > orientation.Collinear) || (line2StartToLine1Orientation < orientation.Collinear && line2EndToLine1Orientation < orientation.Collinear) {
		data.intersectionType = lineintersection.NoIntersection
		return
	}

	line1StartToLine2Orientation := bigxy.OrientationIndex(line2Start, line2End, line1Start)
	line1EndToLine2Orientation := bigxy.OrientationIndex(line2Start, line2End, line1End)

	if (line1StartToLine2Orientation > orientation.Collinear && line1EndToLine2Orientation > orientation.Collinear) || (line1StartToLine2Orientation < 0 && line1EndToLine2Orientation < 0) {
		data.intersectionType = lineintersection.NoIntersection
		return
	}

	collinear := line2StartToLine1Orientation == orientation.Collinear && line2EndToLine1Orientation == orientation.Collinear &&
		line1StartToLine2Orientation == orientation.Collinear && line1EndToLine2Orientation == orientation.Collinear

	if collinear {
		data.intersectionType = computeCollinearIntersection(data, line1Start, line1End, line2Start, line2End)
		return
	}

	/*
	 * At this point we know that there is a single intersection point
	 * (since the lines are not collinear).
	 */

	/*
	 *  Check if the intersection is an endpoint. If it is, copy the endpoint as
	 *  the intersection point. Copying the point rather than computing it
	 *  ensures the point has the exact value, which is important for
	 *  robustness. It is sufficient to simply check for an endpoint which is on
	 *  the other line, since at this point we know that the inputLines must
	 *  intersect.
	 */
	if line2StartToLine1Orientation == orientation.Collinear || line2EndToLine1Orientation == orientation.Collinear ||
		line1StartToLine2Orientation == orientation.Collinear || line1EndToLine2Orientation == orientation.Collinear {
		data.isProper = false

		/*
		 * Check for two equal endpoints.
		 * This is done explicitly rather than by the orientation tests
		 * below in order to improve robustness.
		 *
		 * [An example where the orientation tests fail to be consistent is
		 * the following (where the true intersection is at the shared endpoint
		 * POINT (19.850257749638203 46.29709338043669)
		 *
		 * LINESTRING ( 19.850257749638203 46.29709338043669, 20.31970698357233 46.76654261437082 )
		 * and
		 * LINESTRING ( -48.51001596420236 -22.063180333403878, 19.850257749638203 46.29709338043669 )
		 *
		 * which used to produce the INCORRECT result: (20.31970698357233, 46.76654261437082, NaN)
		 *
		 */
		if internal.Equal(line1Start, 0, line2Start, 0) || internal.Equal(line1Start, 0, line2End, 0) {
			copy(data.intersectionPoints[0], line1Start)
		} else if internal.Equal(line1End, 0, line2Start, 0) || internal.Equal(line1End, 0, line2End, 0) {
			copy(data.intersectionPoints[0], line1End)
		} else if line2StartToLine1Orientation == orientation.Collinear {
			// Now check to see if any endpoint lies on the interior of the other segment.
			copy(data.intersectionPoints[0], line2Start)
		} else if line2EndToLine1Orientation == orientation.Collinear {
			copy(data.intersectionPoints[0], line2End)
		} else if line1StartToLine2Orientation == orientation.Collinear {
			copy(data.intersectionPoints[0], line1Start)
		} else if line1EndToLine2Orientation == orientation.Collinear {
			copy(data.intersectionPoints[0], line1End)
		}
	} else {
		data.isProper = true
		data.intersectionPoints[0] = intersection(data, line1Start, line1End, line2Start, line2End)
	}

	data.intersectionType = lineintersection.PointIntersection
}

func computeCollinearIntersection(data *lineIntersectorData, line1Start, line1End, line2Start, line2End geom.Coord) lineintersection.Type {
	line2StartWithinLine1Bounds := internal.IsPointWithinLineBounds(line2Start, line1Start, line1End)
	line2EndWithinLine1Bounds := internal.IsPointWithinLineBounds(line2End, line1Start, line1End)
	line1StartWithinLine2Bounds := internal.IsPointWithinLineBounds(line1Start, line2Start, line2End)
	line1EndWithinLine2Bounds := internal.IsPointWithinLineBounds(line1End, line2Start, line2End)

	if line1StartWithinLine2Bounds && line1EndWithinLine2Bounds {
		data.intersectionPoints[0] = line1Start
		data.intersectionPoints[1] = line1End
		return lineintersection.CollinearIntersection
	}

	if line2StartWithinLine1Bounds && line2EndWithinLine1Bounds {
		data.intersectionPoints[0] = line2Start
		data.intersectionPoints[1] = line2End
		return lineintersection.CollinearIntersection
	}

	if line2StartWithinLine1Bounds && line1StartWithinLine2Bounds {
		data.intersectionPoints[0] = line2Start
		data.intersectionPoints[1] = line1Start

		return isPointOrCollinearIntersection(data, line2Start, line1Start, line2EndWithinLine1Bounds, line1EndWithinLine2Bounds)
	}
	if line2StartWithinLine1Bounds && line1EndWithinLine2Bounds {
		data.intersectionPoints[0] = line2Start
		data.intersectionPoints[1] = line1End

		return isPointOrCollinearIntersection(data, line2Start, line1End, line2EndWithinLine1Bounds, line1StartWithinLine2Bounds)
	}

	if line2EndWithinLine1Bounds && line1StartWithinLine2Bounds {
		data.intersectionPoints[0] = line2End
		data.intersectionPoints[1] = line1Start

		return isPointOrCollinearIntersection(data, line2End, line1Start, line2StartWithinLine1Bounds, line1EndWithinLine2Bounds)
	}

	if line2EndWithinLine1Bounds && line1EndWithinLine2Bounds {
		data.intersectionPoints[0] = line2End
		data.intersectionPoints[1] = line1End

		return isPointOrCollinearIntersection(data, line2End, line1End, line2StartWithinLine1Bounds, line1StartWithinLine2Bounds)
	}

	return lineintersection.NoIntersection
}

func isPointOrCollinearIntersection(data *lineIntersectorData, lineStart, lineEnd geom.Coord, intersection1, intersection2 bool) lineintersection.Type {
	if internal.Equal(lineStart, 0, lineEnd, 0) && !intersection1 && !intersection2 {
		return lineintersection.PointIntersection
	}
	return lineintersection.CollinearIntersection
}

/**
 * 此方法计算交点的实际值。
 * 通过求交求取最大精度，
 * 坐标通过减去最小坐标值（绝对值）来进行标准化。
 *这有从计算中删除普通有效数字以保持更多精度的效果。
 */
func intersection(data *lineIntersectorData, line1Start, line1End, line2Start, line2End geom.Coord) geom.Coord {
	intPt := intersectionWithNormalization(line1Start, line1End, line2Start, line2End)

	/**
	 * 由于四舍五入，可能计算交点可能在输入的线段外。显然这是不一致的。
	 * 这段代码检查这种情况，并迫使一个更合理的答案。
	 */
	if !isInSegmentEnvelopes(data, intPt) {
		intPt = centralendpoint.GetIntersection(line1Start, line1End, line2Start, line2End)
	}

	// TODO Enable if we add a precision model
	//if precisionModel != null {
	//	precisionModel.makePrecise(intPt);
	//}

	return intPt
}

func intersectionWithNormalization(line1Start, line1End, line2Start, line2End geom.Coord) geom.Coord {
	var line1End1Norm, line1End2Norm, line2End1Norm, line2End2Norm geom.Coord = geom.Coord{0, 0}, geom.Coord{0, 0}, geom.Coord{0, 0}, geom.Coord{0, 0}
	copy(line1End1Norm, line1Start)
	copy(line1End2Norm, line1End)
	copy(line2End1Norm, line2Start)
	copy(line2End2Norm, line2End)

	normPt := geom.Coord{0, 0}
	normalizeToEnvCentre(line1End1Norm, line1End2Norm, line2End1Norm, line2End2Norm, normPt)

	intPt := safeHCoordinateIntersection(line1End1Norm, line1End2Norm, line2End1Norm, line2End2Norm)

	intPt[0] += normPt[0]
	intPt[1] += normPt[1]

	return intPt
}

/**
 * 使用齐次坐标计算线段相交。
 * 舍入误差会导致原始计算失败，（通常是由于线段近似平行）。
 * 如果出现这种情况，则计算一个合理的近似值。
 */
func safeHCoordinateIntersection(line1Start, line1End, line2Start, line2End geom.Coord) geom.Coord {
	intPt, err := hcoords.GetIntersection(line1Start, line1End, line2Start, line2End)
	if err != nil {
		return centralendpoint.GetIntersection(line1Start, line1End, line2Start, line2End)
	}
	return intPt
}

/*
 * 测试一个点是否位于输入的线段中。
 * 在测试过程中，正确计算的交点应当返回 <code>true</code>
 *
 * 由于此测试仅用于调试目的，所以没有尝试优化包含测试。
 *
 * 如果输入点位于输入的线段上，则返回true。
 */
func isInSegmentEnvelopes(data *lineIntersectorData, intersectionPoint geom.Coord) bool {
	intersection1 := internal.IsPointWithinLineBounds(intersectionPoint, data.inputLines[0][0], data.inputLines[0][1])
	intersection2 := internal.IsPointWithinLineBounds(intersectionPoint, data.inputLines[1][0], data.inputLines[1][1])

	return intersection1 && intersection2
}

/**
 * 将所提供的坐标规范化，使其相交点位于原点。
 */
func normalizeToEnvCentre(line1Start, line1End, line2Start, line2End, normPt geom.Coord) {
	// Note: All these "max" checks are inlined for performance.
	// It would be visually cleaner to do that but requires more function calls

	line1MinX := line1End[0]
	if line1Start[0] < line1End[0] {
		line1MinX = line1Start[0]
	}

	line1MinY := line1End[1]
	if line1Start[1] < line1End[1] {
		line1MinY = line1Start[1]
	}
	line1MaxX := line1End[0]
	if line1Start[0] > line1End[0] {
		line1MaxX = line1Start[0]
	}
	line1MaxY := line1End[1]
	if line1Start[1] > line1End[1] {
		line1MaxY = line1Start[1]
	}

	line2MinX := line2End[0]
	if line2Start[0] < line2End[0] {
		line2MinX = line2Start[0]
	}
	line2MinY := line2End[1]
	if line2Start[1] < line2End[1] {
		line2MinY = line2Start[1]
	}
	line2MaxX := line2End[0]
	if line2Start[0] > line2End[0] {
		line2MaxX = line2Start[0]
	}
	line2MaxY := line2End[1]
	if line2Start[1] > line2End[1] {
		line2MaxY = line2Start[1]
	}

	intMinX := line2MinX
	if line1MinX > line2MinX {
		intMinX = line1MinX
	}
	intMaxX := line2MaxX
	if line1MaxX < line2MaxX {
		intMaxX = line1MaxX
	}
	intMinY := line2MinY
	if line1MinY > line2MinY {
		intMinY = line1MinY
	}
	intMaxY := line2MaxY
	if line1MaxY < line2MaxY {
		intMaxY = line1MaxY
	}

	intMidX := (intMinX + intMaxX) / 2.0
	intMidY := (intMinY + intMaxY) / 2.0
	normPt[0] = intMidX
	normPt[1] = intMidY

	line1Start[0] -= normPt[0]
	line1Start[1] -= normPt[1]
	line1End[0] -= normPt[0]
	line1End[1] -= normPt[1]
	line2Start[0] -= normPt[0]
	line2Start[1] -= normPt[1]
	line2End[0] -= normPt[0]
	line2End[1] -= normPt[1]
}
