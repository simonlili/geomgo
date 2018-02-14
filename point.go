package geom

// 一个Point代表一个点
type Point struct {
	geom0
}

// NewPoint 函数根据视图分配一个新点，全部为0值
func NewPoint(l Layout) *Point {
	return NewPointFlat(l, make([]float64, l.Stride()))
}

// NewPointFlat 函数根据视图 l 和传入的点坐标数组分配一个点
func NewPointFlat(l Layout, flatCoords []float64) *Point {
	p := new(Point)
	p.layout = l
	p.stride = l.Stride()
	p.flatCoords = flatCoords
	return p
}
/**
*------------------------------
*				Point（点）相关的方法
*---------------------------------
*/

// Area 函数返回点的面积，只为零
func (p *Point) Area() float64 {
	return 0
}

// Clone 函数返回一个点的拷贝，这个不是别名
func (p *Point) Clone() *Point {
	return deriveClonePoint(p)
}

// Empty 函数返回false
func (p *Point) Empty() bool {
	return false
}

// Length 函数返回点的长度，只为零
func (p *Point) Length() float64 {
	return 0
}

// MustSetCoords方法 设置点的坐标，但是如果有任何错误都会抛出
func (p *Point) MustSetCoords(coords Coord) *Point {
	Must(p.SetCoords(coords))
	return p
}

// SetCoords 函数传入一个坐标Coord,设置点坐标
func (p *Point) SetCoords(coords Coord) (*Point, error) {
	if err := p.setCoords(coords); err != nil {
		return nil, err
	}
	return p, nil
}

// SetSRID 设置点的SRID参考
func (p *Point) SetSRID(srid int) *Point {
	p.srid = srid
	return p
}

// Swap 函数交换两个点
func (p *Point) Swap(p2 *Point) {
	*p, *p2 = *p2, *p
}

// X 函数返回点的x坐标
func (p *Point) X() float64 {
	return p.flatCoords[0]
}

// Y 函数返回点的y坐标
func (p *Point) Y() float64 {
	return p.flatCoords[1]
}

// Z 返回点的Z坐标，如果点所属的视图类型不具有z轴返回0
func (p *Point) Z() float64 {
	zIndex := p.layout.ZIndex()
	if zIndex == -1 {
		return 0
	}
	return p.flatCoords[zIndex]
}

// M 函数返回点的m属性值，如果点所属的的视图类型不具有 m 轴返回0
func (p *Point) M() float64 {
	mIndex := p.layout.MIndex()
	if mIndex == -1 {
		return 0
	}
	return p.flatCoords[mIndex]
}
