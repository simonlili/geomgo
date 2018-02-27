package xy

import (
	"github.com/chengxiaoer/geomGo"
	"github.com/chengxiaoer/geomGo/xy/internal"
)

// LinesCentroid函数 计算参数传入的所有线状要素的质心
//
// Algorithm: 计算各线段的中点段长度加权平均。
func LinesCentroid(line *geom.LineString, extraLines ...*geom.LineString) (centroid geom.Coord) {
	calculator := NewLineCentroidCalculator(line.Layout())
	calculator.AddLine(line)

	for _, l := range extraLines {
		calculator.AddLine(l)
	}

	return calculator.GetCentroid()
}

// LinearRingsCentroid函数 计算参数传入的所有线环的质心
//
// Algorithm: 计算各线段的中点段长度加权平均。
func LinearRingsCentroid(line *geom.LinearRing, extraLines ...*geom.LinearRing) (centroid geom.Coord) {
	calculator := NewLineCentroidCalculator(line.Layout())
	calculator.AddLinearRing(line)

	for _, l := range extraLines {
		calculator.AddLinearRing(l)
	}

	return calculator.GetCentroid()
}

// MultiLineCentroid函数 计算MultiLineString的质心
//
// Algorithm: 计算线段长度加权的所有线段的平均值。
func MultiLineCentroid(line *geom.MultiLineString) (centroid geom.Coord) {
	calculator := NewLineCentroidCalculator(line.Layout())
	start := 0
	for _, end := range line.Ends() {
		calculator.addLine(line.FlatCoords(), start, end)
		start = end
	}

	return calculator.GetCentroid()
}

// LineCentroidCalculator结构 是质心计算的数据结构
//  该结构没有默认零值,必须使用NewLineCentroid函数创建
type LineCentroidCalculator struct {
	layout      geom.Layout
	stride      int
	centSum     geom.Coord
	totalLength float64
}

// NewLineCentroidCalculator 创建计算器的新实例。
// 计算器创建后多边形、线要素或线性环可以添加
// GetCentroid方法 可以用于在任何点获得当前质心。
// 每次添加几何体时，质心都会发生自然变化。
func NewLineCentroidCalculator(layout geom.Layout) *LineCentroidCalculator {
	return &LineCentroidCalculator{
		layout:  layout,
		stride:  layout.Stride(),
		centSum: geom.Coord(make([]float64, layout.Stride())),
	}
}

/**
*------------------------------
*				LineCentroidCalculator（线性质心计算器）相关的方法
*---------------------------------
*/

// GetCentroid 获取质心，如果没有几何类型加入则返回0
func (calc *LineCentroidCalculator) GetCentroid() geom.Coord {
	cent := geom.Coord(make([]float64, calc.layout.Stride()))
	cent[0] = calc.centSum[0] / calc.totalLength
	cent[1] = calc.centSum[1] / calc.totalLength
	return cent
}

// AddPolygon方法 向计算器中添加多边形。
func (calc *LineCentroidCalculator) AddPolygon(polygon *geom.Polygon) *LineCentroidCalculator {
	for i := 0; i < polygon.NumLinearRings(); i++ {
		calc.AddLinearRing(polygon.LinearRing(i))
	}

	return calc
}

// AddLine方法 向计算器中添加线段
func (calc *LineCentroidCalculator) AddLine(line *geom.LineString) *LineCentroidCalculator {
	coords := line.FlatCoords()
	calc.addLine(coords, 0, len(coords))
	return calc
}

// AddLinearRing方法 向计算器中添加线环
func (calc *LineCentroidCalculator) AddLinearRing(line *geom.LinearRing) *LineCentroidCalculator {
	coords := line.FlatCoords()
	calc.addLine(coords, 0, len(coords))
	return calc
}

func (calc *LineCentroidCalculator) addLine(line []float64, startLine, endLine int) {
	lineMinusLastPoint := endLine - calc.stride
	for i := startLine; i < lineMinusLastPoint; i += calc.stride {
		segmentLen := internal.Distance2D(geom.Coord(line[i:i+2]), geom.Coord(line[i+calc.stride:i+calc.stride+2]))
		calc.totalLength += segmentLen

		midx := (line[i] + line[i+calc.stride]) / 2
		calc.centSum[0] += segmentLen * midx
		midy := (line[i+1] + line[i+calc.stride+1]) / 2
		calc.centSum[1] += segmentLen * midy
	}
}
