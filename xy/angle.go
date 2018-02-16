package xy

import (
	"math"

	"github.com/chengxiaoer/go-geom"
	"github.com/chengxiaoer/go-geom/xy/orientation"
)

const piTimes2 = math.Pi * 2

// Angle函数 计算向量从po到p1的角度，此角度相当于x轴正方向，角度范围为[-180,180]
func Angle(p0, p1 geom.Coord) float64 {
	dx := p1[0] - p0[0]
	dy := p1[1] - p0[1]
	return math.Atan2(dy, dx)
}

// AngleFromOrigin函数 计算向量从（0，0）到p点的角度。此角度相对于x轴正方向,角度范围为（-180，180]
func AngleFromOrigin(p geom.Coord) float64 {
	return math.Atan2(p[1], p[0])
}

// IsAcute函数 测试一个角度是否为锐角
// Note: 对于非常接近90度的角度来说，不太精确。
func IsAcute(endpoint1, base, endpoint2 geom.Coord) bool {
	// relies on fact that A dot B is positive iff A ang B is acute
	dx0 := endpoint1[0] - base[0]
	dy0 := endpoint1[1] - base[1]
	dx1 := endpoint2[0] - base[0]
	dy1 := endpoint2[1] - base[1]
	dotprod := dx0*dx1 + dy0*dy1
	return dotprod > 0
}

// IsObtuse函数 测试一个角度是否为钝角
// Note: 当角度非常接近90度时，不太准确
func IsObtuse(endpoint1, base, endpoint2 geom.Coord) bool {
	// relies on fact that A dot B is negative iff A ang B is obtuse
	dx0 := endpoint1[0] - base[0]
	dy0 := endpoint1[1] - base[1]
	dx1 := endpoint2[0] - base[0]
	dy1 := endpoint2[1] - base[1]
	dotprod := dx0*dx1 + dy0*dy1
	return dotprod < 0
}

// AngleBetween函数 计算向量间的最小夹角
//计算的角度范围在（0，180]之间
//
// Param tip1 - 向量的顶点
// param tail - 每一个向量的尾部
// param tip2 - 每一个向量的顶点
func AngleBetween(tip1, tail, tip2 geom.Coord) float64 {
	a1 := Angle(tail, tip1)
	a2 := Angle(tail, tip2)

	return Diff(a1, a2)
}

// AngleBetweenOriented函数 计算计算两向量间的最小夹角（有两种结果）.
// 计算的结果范围为（-180，180]
// 一个正数结果对应了从v1向量到v2向量逆时针旋转
// 负数结果对应了从v1到v2顺时针旋转所成的角
// 0 表示两向量间不存在夹角（方向一致）
func AngleBetweenOriented(tip1, tail, tip2 geom.Coord) float64 {
	a1 := Angle(tail, tip1)
	a2 := Angle(tail, tip2)
	angDel := a2 - a1

	return Normalize(angDel)
}

// InteriorAngle函数 计算环的两个部分之间的内角。
// 以顺时针为正向,计算结果的范围为 [0, 2Pi]
func InteriorAngle(p0, p1, p2 geom.Coord) float64 {
	anglePrev := Angle(p1, p0)
	angleNext := Angle(p1, p2)
	return math.Abs(angleNext - anglePrev)
}

// AngleOrientation函数 一个角度是否必须顺时针或逆时针旋转另一个角度
func AngleOrientation(ang1, ang2 float64) orientation.Type {
	crossproduct := math.Sin(ang2 - ang1)

	switch {
	case crossproduct > 0:
		return orientation.CounterClockwise
	case crossproduct < 0:
		return orientation.Clockwise
	default:
		return orientation.Collinear
	}
}

// Normalize函数 计算一个角的归一化值，它是在（-180，180]范围内的等效角
func Normalize(angle float64) float64 {
	for angle > math.Pi {
		angle -= piTimes2
	}
	for angle <= -math.Pi {
		angle += piTimes2
	}
	return angle
}

// NormalizePositive函数 计算一个角的归一化值，它是在[0,360]范围内的等效角
// E.g.:
// * normalizePositive(0.0) = 0.0
// * normalizePositive(-PI) = PI
// * normalizePositive(-2PI) = 0.0
// * normalizePositive(-3PI) = PI
// * normalizePositive(-4PI) = 0
// * normalizePositive(PI) = PI
// * normalizePositive(2PI) = 0.0
// * normalizePositive(3PI) = PI
// * normalizePositive(4PI) = 0.0
func NormalizePositive(angle float64) float64 {
	if angle < 0.0 {
		for angle < 0.0 {
			angle += piTimes2
		}
		// in case round-off error bumps the value over
		// 舍去误差使值超过
		if angle >= piTimes2 {
			angle = 0.0
		}
	} else {
		for angle >= piTimes2 {
			angle -= piTimes2
		}
		// in case round-off error bumps the value under
		// 舍去误差使值小于
		if angle < 0.0 {
			angle = 0.0
		}
	}
	return angle
}

// Diff函数 计算非定向的两个向量的最小角度。
// 假设角被归一化到范围[-π，π]。
//结果将在[0,π]之间
// Param ang1 - the angle of one vector (in [-Pi, Pi] )
// Param ang2 - the angle of the other vector (in range [-Pi, Pi] )
func Diff(ang1, ang2 float64) float64 {
	var delAngle float64

	if ang1 < ang2 {
		delAngle = ang2 - ang1
	} else {
		delAngle = ang1 - ang2
	}

	if delAngle > math.Pi {
		delAngle = piTimes2 - delAngle
	}

	return delAngle
}
