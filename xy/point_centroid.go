package xy

import "github.com/chengxiaoer/geomGo"

// PointsCentroid函数 计算参数传入的点的质心
//
// 算法实现：所有点的平均值
func PointsCentroid(point *geom.Point, extra ...*geom.Point) geom.Coord {
	calc := NewPointCentroidCalculator()
	calc.AddCoord(geom.Coord(point.FlatCoords()))

	for _, p := range extra {
		calc.AddCoord(geom.Coord(p.FlatCoords()))
	}

	return calc.GetCentroid()
}

// MultiPointCentroid函数 计算的点的集合的质心
//
// 算法实现：集合中所有点的平均值
func MultiPointCentroid(point *geom.MultiPoint) geom.Coord {
	calc := NewPointCentroidCalculator()
	coords := point.FlatCoords()
	stride := point.Layout().Stride()
	for i := 0; i < len(coords); i += stride {
		calc.AddCoord(geom.Coord(coords[i : i+stride]))
	}

	return calc.GetCentroid()
}

// PointsCentroidFlat函数 计算点数组中的点的质心
// 布局仅用于确定如何查找每个坐标，x-y坐标每个点必须的参数
// 算法实现: 所有点的平均值
func PointsCentroidFlat(layout geom.Layout, pointData []float64) geom.Coord {
	calc := NewPointCentroidCalculator()

	coord := geom.Coord{0, 0}
	stride := layout.Stride()
	arrayLen := len(pointData)
	for i := 0; i < arrayLen; i += stride {
		coord[0] = pointData[i]
		coord[1] = pointData[i+1]
		calc.AddCoord(coord)
	}

	return calc.GetCentroid()
}

// PointCentroidCalculator结构 点质心计算的数据组织结构。
// 该结构不能使用0值来进行初始化，必须使用 NewPointCentroid 函数来进行创建
type PointCentroidCalculator struct {
	ptCount int
	centSum geom.Coord
}

// NewPointCentroidCalculator函数 创建点的计算器结构/对象
// 计算器对象创建后可以继续添加坐标或点
//使用 GetCentedrid 方法可以获取最新的计算结果
func NewPointCentroidCalculator() PointCentroidCalculator {
	return PointCentroidCalculator{centSum: geom.Coord{0, 0}}
}

/**
*--------------------------------------------------------
*				PointCentroidCalculator（点质心计算器）相关的方法
*-----------------------------------------------------------
*/

// AddPoint方法 向计算器中添加点
func (calc *PointCentroidCalculator) AddPoint(point *geom.Point) {
	calc.AddCoord(geom.Coord(point.FlatCoords()))
}

// AddCoord方法 向计算器中添加点坐标
func (calc *PointCentroidCalculator) AddCoord(point geom.Coord) {
	calc.ptCount++
	calc.centSum[0] += point[0]
	calc.centSum[1] += point[1]
}

// GetCentroid方法 获取最新的质心计算结果. 如果计算器中没有点则返回0
func (calc *PointCentroidCalculator) GetCentroid() geom.Coord {
	cent := geom.Coord{0, 0}
	cent[0] = calc.centSum[0] / float64(calc.ptCount)
	cent[1] = calc.centSum[1] / float64(calc.ptCount)
	return cent
}
