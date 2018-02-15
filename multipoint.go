package geom

// MultiPoint 是点的集合
type MultiPoint struct {
	geom1
}

// NewMultiPoint 函数返回一个新的、空的 MultiPoint
func NewMultiPoint(layout Layout) *MultiPoint {
	return NewMultiPointFlat(layout, nil)
}

// NewMultiPointFlat函数 创建一个新的，有坐标的MultiPoint
func NewMultiPointFlat(layout Layout, flatCoords []float64) *MultiPoint {
	mp := new(MultiPoint)
	mp.layout = layout
	mp.stride = layout.Stride()
	mp.flatCoords = flatCoords
	return mp
}
/**
*------------------------------
*				MultiPoint（多点）相关的方法
*---------------------------------
*/

// Area方法 返回0
func (mp *MultiPoint) Area() float64 {
	return 0
}

// Clone 一个深层拷贝的点
func (mp *MultiPoint) Clone() *MultiPoint {
	return deriveCloneMultiPoint(mp)
}

// Empty方法 在点集为空时，返回true
func (mp *MultiPoint) Empty() bool {
	return mp.NumPoints() == 0
}

// Length方法 返回0
func (mp *MultiPoint) Length() float64 {
	return 0
}

// MustSetCoords方法 设置点集的坐标，当有任何错误时均会抛出
func (mp *MultiPoint) MustSetCoords(coords []Coord) *MultiPoint {
	Must(mp.SetCoords(coords))
	return mp
}

// SetCoords方法 设置点集坐标
func (mp *MultiPoint) SetCoords(coords []Coord) (*MultiPoint, error) {
	if err := mp.setCoords(coords); err != nil {
		return nil, err
	}
	return mp, nil
}

// SetSRID方法 设置点集的坐标系参考
func (mp *MultiPoint) SetSRID(srid int) *MultiPoint {
	mp.srid = srid
	return mp
}

// NumPoints方法 返回点集中点的数目
func (mp *MultiPoint) NumPoints() int {
	return mp.NumCoords()
}

// Point方法 返回点击中指定索引的点
func (mp *MultiPoint) Point(i int) *Point {
	return NewPointFlat(mp.layout, mp.Coord(i))
}

// Push方法 向点集中添加新的点
func (mp *MultiPoint) Push(p *Point) error {
	if p.layout != mp.layout {
		return ErrLayoutMismatch{Got: p.layout, Want: mp.layout}
	}
	mp.flatCoords = append(mp.flatCoords, p.flatCoords...)
	return nil
}

// Swap方法 将本点集与传入的点集互相交换
func (mp *MultiPoint) Swap(mp2 *MultiPoint) {
	*mp, *mp2 = *mp2, *mp
}
