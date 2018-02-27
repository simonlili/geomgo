package lineintersection

import "github.com/chengxiaoer/geomGo"

// Type 枚举了两条线的位置关系
type Type int

const (
	// NoIntersection 表示不相交
	NoIntersection Type = iota
	// PointIntersection 表示相交于某一点
	PointIntersection
	// CollinearIntersection 表示两条线相互重叠
	CollinearIntersection
)

var labels = [3]string{"NoIntersection", "PointIntersection", "CollinearIntersection"}

func (t Type) String() string {
	return labels[t]
}

// Result LineIntersectsLine函数的返回结果 .
// 它包含交叉点（S），并指示有交集的类型（或没有交集）。
type Result struct {
	intersectionType Type
	intersection     []geom.Coord
}

// NewResult函数 创建一个结果对象
func NewResult(intersectionType Type, intersection []geom.Coord) Result {
	return Result{
		intersectionType: intersectionType,
		intersection:     intersection}
}
/**
*------------------------------
*			Result（结果集）相关的方法
*---------------------------------
*/
// HasIntersection方法 如果交叉返回true
func (i *Result) HasIntersection() bool {
	return i.intersectionType != NoIntersection
}

// Type方法 返回两条线的交叉类型
func (i *Result) Type() Type {
	return i.intersectionType
}

// Intersection 返回一个交点数组
// 如果 type 是 PointIntersection ，返回一个点（第一个点）
// 如果 type 是  CollinearIntersection ，返回两个点，一个是线的起点，一个是线的终点
func (i *Result) Intersection() []geom.Coord {
	return i.intersection
}
