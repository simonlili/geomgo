// Package bigxy 包含平面（XY）数据的强大地理功能。
//计算是使用大浮点对象实现的，具有最高的精确度和健壮性。
//
// Note:要求所有坐标都必须有x和y坐标，在geom.Coord数组的第一、二位置上。
// 鉴于坐标可以是任何大小，除了X和Y是在这些计算中忽略了所有的数据。
package bigxy

import (
	"math"
	"math/big"

	"github.com/chengxiaoer/geomGo"
	"github.com/chengxiaoer/geomGo/xy/orientation"
)

// dpSafeEpsilon 该值是安全的，比big.Flaot的相对误差的最大精度数大。
var dpSafeEpsilon = 1e-15

// OrientationIndex函数 返回一个点的主要方向，相对于由 vectorOrigin-vectorEnd 组成的向量
//
// vectorOrigin - 矢量的起点
// vectorEnd - 矢量的终点
// point - 点计算方向
//
// Returns CounterClockwise 如果这点相对于向量是逆时针旋转的
// Returns Clockwise 如果这个点相对于向量是顺时针旋转的
// Returns Collinear 如果该点与向量共线
func OrientationIndex(vectorOrigin, vectorEnd, point geom.Coord) orientation.Type {
	// 定向指标快速滤波器
	// 避免在许多情况下使用慢扩展精度运算。
	index := orientationIndexFilter(vectorOrigin, vectorEnd, point)
	if index <= 1 {
		return index
	}

	var dx1, dy1, dx2, dy2 big.Float

	// 归一化坐标
	dx1.SetFloat64(vectorEnd[0]).Add(&dx1, big.NewFloat(-vectorOrigin[0]))
	dy1.SetFloat64(vectorEnd[1]).Add(&dy1, big.NewFloat(-vectorOrigin[1]))
	dx2.SetFloat64(point[0]).Add(&dx2, big.NewFloat(-vectorEnd[0]))
	dy2.SetFloat64(point[1]).Add(&dy2, big.NewFloat(-vectorEnd[1]))

	// 计算因子.  计算的性能主要体现在 dx1 上
	dx1.Mul(&dx1, &dy2)
	dy1.Mul(&dy1, &dx2)
	dx1.Sub(&dx1, &dy1)

	return orientationBasedOnSignForBig(dx1)
}

// Intersection函数 使用数学计算两条直线的交点。大浮点运算.
// Line被认为是无限长的线（可以双向延展）。  For example, (0,0), (1, 0) and (2, 1) (2, 2) 将相较于 (2, 0)
// 目前未处理平行线的情况。
func Intersection(line1Start, line1End, line2Start, line2End geom.Coord) geom.Coord {
	var denom1, denom2, denom, tmp1, tmp2 big.Float

	denom1.SetFloat64(line2End[1]).Sub(&denom1, tmp2.SetFloat64(line2Start[1])).Mul(&denom1, tmp1.SetFloat64(line1End[0]).Sub(&tmp1, tmp2.SetFloat64(line1Start[0])))
	denom2.SetFloat64(line2End[0]).Sub(&denom2, tmp2.SetFloat64(line2Start[0])).Mul(&denom2, tmp1.SetFloat64(line1End[1]).Sub(&tmp1, tmp2.SetFloat64(line1Start[1])))
	denom.Sub(&denom1, &denom2)

	// Cases:
	// - denom is 0 如果直线相互平行
	// - 如果 fracp的值在0和1之间，交点位于线段的p上
	// - 如果fracQ的值在0和1之间，交点位于线段Q上

	// 重用以前的变量以获得性能
	numx1 := &denom1
	numx2 := &denom2
	var numx big.Float

	numx1.SetFloat64(line2End[0]).Sub(numx1, tmp2.SetFloat64(line2Start[0])).Mul(numx1, tmp1.SetFloat64(line1Start[1]).Sub(&tmp1, tmp2.SetFloat64(line2Start[1])))
	numx2.SetFloat64(line2End[1]).Sub(numx2, tmp2.SetFloat64(line2Start[1])).Mul(numx2, tmp1.SetFloat64(line1Start[0]).Sub(&tmp1, tmp2.SetFloat64(line2Start[0])))
	numx.Sub(numx1, numx2)

	fracP, _ := numx.Quo(&numx, &denom).Float64()

	x, _ := numx1.SetFloat64(line1Start[0]).Add(numx1, tmp2.SetFloat64(line1End[0])).Sub(numx1, tmp2.SetFloat64(line1Start[0])).Mul(numx1, tmp1.SetFloat64(fracP)).Float64()

	// 重用以前的变量以获得性能
	numy1 := &denom1
	numy2 := &denom2
	var numy big.Float

	numy1.SetFloat64(line1End[0]).Sub(numy1, tmp2.SetFloat64(line1Start[0])).Mul(numy1, tmp1.SetFloat64(line1Start[1]).Sub(&tmp1, tmp2.SetFloat64(line2Start[1])))
	numy2.SetFloat64(line1End[1]).Sub(numy2, tmp2.SetFloat64(line1Start[1])).Mul(numy2, tmp1.SetFloat64(line1Start[0]).Sub(&tmp1, tmp2.SetFloat64(line2Start[0])))
	numy.Sub(numy1, numy2)

	fracQ, _ := numy.Quo(&numy, &denom).Float64()

	tmp2.SetFloat64(line1End[1]).Sub(&tmp2, tmp1.SetFloat64(line1Start[1]))

	if tmp2.IsInf() && fracQ == 0 || tmp1.SetFloat64(0).Cmp(&tmp2) == 0 && math.IsInf(fracQ, 0) {
		// can't perform calculation
		return geom.Coord{math.Inf(1), math.Inf(1)}
	}

	y, _ := numx1.SetFloat64(line1Start[1]).Add(numx1, tmp2.Mul(&tmp2, tmp1.SetFloat64(fracQ))).Float64()

	return geom.Coord{x, y}
}

/////////////////  实现 /////////////////////////////////

// 一种计算三坐标方位指数的过滤器。
//
// 如果可以，使用标准DP算法安全地计算方向, 这个goroutine返回方向索引。
// 其他情况, 一个值 i > 1 将被返回
// 在这种情况下，方向指数必须用其他更可靠的方法计算。
// 该过滤器是快速计算，所以可以用来避免使用较慢的强健方法，除非他们真的需要，从而提供更好的平均性能。
//
// Uses an approach due to Jonathan Shewchuk, which is in the public domain（这是一个公共领域）.
// Jonathan Shewchuk---分治法中三角形的几何信息和拓扑信息的操作
// Return 方向索引，如果它被安全的计算
// Return i > 1，如果这个这个方向索引不是安全计算的
func orientationIndexFilter(vectorOrigin, vectorEnd, point geom.Coord) orientation.Type {
	var detsum float64

	detleft := (vectorOrigin[0] - point[0]) * (vectorEnd[1] - point[1])
	detright := (vectorOrigin[1] - point[1]) * (vectorEnd[0] - point[0])
	det := detleft - detright

	if detleft > 0.0 {
		if detright <= 0.0 {
			return orientationBasedOnSign(det)
		}

		detsum = detleft + detright
	} else if detleft < 0.0 {
		if detright >= 0.0 {
			return orientationBasedOnSign(det)
		}
		detsum = -detleft - detright
	} else {
		return orientationBasedOnSign(det)
	}

	errbound := dpSafeEpsilon * detsum
	if (det >= errbound) || (-det >= errbound) {
		return orientationBasedOnSign(det)
	}

	return 2
}

func orientationBasedOnSign(x float64) orientation.Type {
	if x > 0 {
		return orientation.CounterClockwise
	}
	if x < 0 {
		return orientation.Clockwise
	}
	return orientation.Collinear
}
func orientationBasedOnSignForBig(x big.Float) orientation.Type {
	if x.IsInf() {
		return orientation.Collinear
	}
	switch x.Sign() {
	case -1:
		return orientation.Clockwise
	case 0:
		return orientation.Collinear
	default:
		return orientation.CounterClockwise
	}
}
