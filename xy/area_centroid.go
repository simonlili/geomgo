// 定义了几何图形质心计算相关函数
package xy
import (
	"math"

	"github.com/chengxiaoer/geomGo"
	"github.com/chengxiaoer/geomGo/xy/internal"
)

// PolygonsCentroid函数 计算一个多边形的几何中心
//
// 算法：
// 基于通常的算法计算的质心作为一个区域的质心加权分解成三角形（可能有重叠）。
//
// 该算法已扩展到处理多个多边形的孔。
//
// See http://www.faqs.org/faqs/graphics/algorithms-faq/ for 具体细节关于这一基本算法.
//
// 该代码还扩展到处理退化（零面积）多边形
//
// 在这种情况下，将返回多边形中线段的质心。
func PolygonsCentroid(polygon *geom.Polygon, extraPolys ...*geom.Polygon) (centroid geom.Coord) {

	calc := NewAreaCentroidCalculator(polygon.Layout())
	calc.AddPolygon(polygon)
	for _, p := range extraPolys {
		calc.AddPolygon(p)
	}
	return calc.GetCentroid()

}

// MultiPolygonCentroid 计算区域几何（多边形集合）的形心。（multipolygon）
//
// 算法：
// 基于通常的算法计算的质心作为一个区域的质心加权分解成三角形（可能有重叠）。
//
// 该算法已扩展到处理多个多边形的孔。
//
// See http://www.faqs.org/faqs/graphics/algorithms-faq/ for 具体细节关于这一基本算法.
//
// 该代码还扩展到处理退化（零面积）多边形。
//
// 在这种情况下，将返回多边形中线段的质心。
func MultiPolygonCentroid(polygon *geom.MultiPolygon) (centroid geom.Coord) {

	calc := NewAreaCentroidCalculator(polygon.Layout())
	for i := 0; i < polygon.NumPolygons(); i++ {
		calc.AddPolygon(polygon.Polygon(i))
	}
	return calc.GetCentroid()

}

// AreaCentroidCalculator 是质心计算数据的数据结构。这类型无法使用其0的价值，
//它必须使用newareacentroid函数来创建
type AreaCentroidCalculator struct {
	layout        geom.Layout
	stride        int
	basePt        geom.Coord
	triangleCent3 geom.Coord // 三角形质心的临时变量
	areasum2      float64    // 局部地区和
	cg3           geom.Coord // 部分质心的总和

	centSum     geom.Coord // 线性质心计算的数据，如果需要时存在
	totalLength float64
}

// NewAreaCentroidCalculator函数 创建计算的新实例。
// 创建计算器后，可以添加多边形。
// GetCentroid方法 在任何点都可以得到当前的质心，每次添加多边形时，质心都会自然变化。
func NewAreaCentroidCalculator(layout geom.Layout) *AreaCentroidCalculator {
	return &AreaCentroidCalculator{
		layout:        layout,
		stride:        layout.Stride(),
		centSum:       geom.Coord(make([]float64, layout.Stride())),
		triangleCent3: geom.Coord(make([]float64, layout.Stride())),
		cg3:           geom.Coord(make([]float64, layout.Stride())),
	}
}


/**
*------------------------------
*				AreaCentroidCalculator（质心计算器）相关的方法
*---------------------------------
*/
// GetCentroid方法 获得当前计算的质心。返回一个0，如果没有几何已添加
func (calc *AreaCentroidCalculator) GetCentroid() geom.Coord {
	cent := geom.Coord(make([]float64, calc.stride))

	if calc.centSum == nil {
		return cent
	}

	if math.Abs(calc.areasum2) > 0.0 {
		cent[0] = calc.cg3[0] / 3 / calc.areasum2
		cent[1] = calc.cg3[1] / 3 / calc.areasum2
	} else {
		// 如果多边形退化，则计算线性质心。
		cent[0] = calc.centSum[0] / calc.totalLength
		cent[1] = calc.centSum[1] / calc.totalLength
	}
	return cent
}

// AddPolygon方法 向计算器中添加多边形
func (calc *AreaCentroidCalculator) AddPolygon(polygon *geom.Polygon) {

	calc.setBasePoint(polygon.Coord(0))

	calc.addShell(polygon.LinearRing(0).FlatCoords())
	for i := 1; i < polygon.NumLinearRings(); i++ {
		calc.addHole(polygon.LinearRing(i).FlatCoords())
	}
}

func (calc *AreaCentroidCalculator) setBasePoint(basePt geom.Coord) {
	if calc.basePt == nil {
		calc.basePt = basePt
	}
}

func (calc *AreaCentroidCalculator) addShell(pts []float64) {
	stride := calc.stride

	isPositiveArea := !IsRingCounterClockwise(calc.layout, pts)
	p1 := geom.Coord{0, 0}
	p2 := geom.Coord{0, 0}

	for i := 0; i < len(pts)-stride; i += stride {
		p1[0] = pts[i]
		p1[1] = pts[i+1]
		p2[0] = pts[i+stride]
		p2[1] = pts[i+stride+1]
		calc.addTriangle(calc.basePt, p1, p2, isPositiveArea)
	}
	calc.addLinearSegments(pts)
}
func (calc *AreaCentroidCalculator) addHole(pts []float64) {
	stride := calc.stride

	isPositiveArea := IsRingCounterClockwise(calc.layout, pts)
	p1 := geom.Coord{0, 0}
	p2 := geom.Coord{0, 0}

	for i := 0; i < len(pts)-stride; i += stride {
		p1[0] = pts[i]
		p1[1] = pts[i+1]
		p2[0] = pts[i+stride]
		p2[1] = pts[i+stride+1]
		calc.addTriangle(calc.basePt, p1, p2, isPositiveArea)
	}
	calc.addLinearSegments(pts)
}

func (calc *AreaCentroidCalculator) addTriangle(p0, p1, p2 geom.Coord, isPositiveArea bool) {
	sign := float64(1.0)
	if isPositiveArea {
		sign = -1.0
	}
	centroid3(p0, p1, p2, calc.triangleCent3)
	area2 := area2(p0, p1, p2)
	calc.cg3[0] += sign * area2 * calc.triangleCent3[0]
	calc.cg3[1] += sign * area2 * calc.triangleCent3[1]
	calc.areasum2 += sign * area2
}

//centroid3函数 返回三角形p1-p2-p3质心的三倍
// 3的系数留在允许除法直到以后避免。
func centroid3(p1, p2, p3, c geom.Coord) {
	c[0] = p1[0] + p2[0] + p3[0]
	c[1] = p1[1] + p2[1] + p3[1]
}

// Returns twice the signed area of the triangle p1-p2-p3,
// positive if a,b,c are oriented ccw, and negative if cw.
func area2(p1, p2, p3 geom.Coord) float64 {
	return (p2[0]-p1[0])*(p3[1]-p1[1]) - (p3[0]-p1[0])*(p2[1]-p1[1])
}

//addLinearSegments方法 添加由线性坐标系的坐标阵列定义的线性段。
// 这是在多边形具有零面积的情况下进行的，在这种情况下，计算线性质心。
//
// Param pts - 一个坐标数组
func (calc *AreaCentroidCalculator) addLinearSegments(pts []float64) {
	stride := calc.stride
	for i := 0; i < len(pts)-stride; i += stride {
		segmentLen := internal.Distance2D(geom.Coord(pts[i:i+2]), pts[i+stride:i+stride+2])
		calc.totalLength += segmentLen

		midx := (pts[i] + pts[i+stride]) / 2
		calc.centSum[0] += segmentLen * midx
		midy := (pts[i+1] + pts[i+stride+1]) / 2
		calc.centSum[1] += segmentLen * midy
	}
}
