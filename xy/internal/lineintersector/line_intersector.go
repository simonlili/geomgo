package lineintersector

import (
	"math"

	"github.com/chengxiaoer/go-geom"
	"github.com/chengxiaoer/go-geom/xy/lineintersection"
)

// Strategy is 线交点的接口
type Strategy interface {
	computePointOnLineIntersection(data *lineIntersectorData, p, lineEndpoint1, lineEndpoint2 geom.Coord)
	computeLineOnLineIntersection(data *lineIntersectorData, line1End1, line1End2, line2End1, line2End2 geom.Coord)
}

// PointIntersectsLine函数 测试点是否在线上
func PointIntersectsLine(strategy Strategy, point, lineStart, lineEnd geom.Coord) (hasIntersection bool) {
	intersectorData := &lineIntersectorData{
		strategy:           strategy,
		inputLines:         [2][2]geom.Coord{{lineStart, lineEnd}, {}},
		intersectionPoints: [2]geom.Coord{{0, 0}, {0, 0}},
	}

	intersectorData.pa = intersectorData.intersectionPoints[0]
	intersectorData.pb = intersectorData.intersectionPoints[1]

	strategy.computePointOnLineIntersection(intersectorData, point, lineStart, lineEnd)

	return intersectorData.intersectionType != lineintersection.NoIntersection
}

// LineIntersectsLine函数 测试第一条直线(line1Start,line1End)与第二条直线(line2Start, line2End)是否相交。
// and 返回表示有相交类型、相交点的数据结构
//查看 lineintersection对象 了解更详细的解释结果
func LineIntersectsLine(strategy Strategy, line1Start, line1End, line2Start, line2End geom.Coord) lineintersection.Result {
	intersectorData := &lineIntersectorData{
		strategy:           strategy,
		inputLines:         [2][2]geom.Coord{{line2Start, line2End}, {line1Start, line1End}},
		intersectionPoints: [2]geom.Coord{{0, 0}, {0, 0}},
	}

	intersectorData.pa = intersectorData.intersectionPoints[0]
	intersectorData.pb = intersectorData.intersectionPoints[1]

	strategy.computeLineOnLineIntersection(intersectorData, line1Start, line1End, line2Start, line2End)

	var intersections []geom.Coord

	switch intersectorData.intersectionType {
	case lineintersection.NoIntersection:
		intersections = []geom.Coord{}
	case lineintersection.PointIntersection:
		intersections = intersectorData.intersectionPoints[:1]
	case lineintersection.CollinearIntersection:
		intersections = intersectorData.intersectionPoints[:2]
	}
	return lineintersection.NewResult(intersectorData.intersectionType, intersections)
}

// 计算期间的一个内部数据结构
type lineIntersectorData struct {
	// new Coordinate[2][2];
	inputLines [2][2]geom.Coord

	// 如果只有一个交点，然后0索引下的坐标将包含交叉点
	// 如果共线（线重叠），两个坐标表示重叠线的起始点和结束点。
	intersectionPoints [2]geom.Coord
	intersectionType   lineintersection.Type

	// 交叉线端点的索引，沿着相应的行顺序。
	isProper bool
	pa, pb   geom.Coord
	strategy Strategy
}

/**
 *  RParameter 计算 the parameter for the point p
 *  in the parameterized equation
 *  of the line from p1 to p2.
 *  This is equal to the 'distance' of p along p1-p2
 */
func rParameter(p1, p2, p geom.Coord) float64 {
	var r float64
	// compute maximum delta, for numerical stability
	// also handle case of p1-p2 being vertical or horizontal
	dx := math.Abs(p2[0] - p1[0])
	dy := math.Abs(p2[1] - p1[1])
	if dx > dy {
		r = (p[0] - p1[0]) / (p2[0] - p1[0])
	} else {
		r = (p[1] - p1[1]) / (p2[1] - p1[1])
	}
	return r
}
