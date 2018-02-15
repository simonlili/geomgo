package geom

// LinearRing 是线环对象。LinearRing 是一个封闭的 LineString 即起始和终止点有相同的坐标值。
//Polygon由LinearRing围成。LinearRing的创建方法与LineString是一样的，
//惟一不同的LinearRing必须要闭合。
type LinearRing struct {
	geom1
}

// NewLinearRing 创建一个没有坐标的线环对象
func NewLinearRing(layout Layout) *LinearRing {
	return NewLinearRingFlat(layout, nil)
}

// NewLinearRingFlat 创建一个已传入坐标为控制点的LinearRing
func NewLinearRingFlat(layout Layout, flatCoords []float64) *LinearRing {
	lr := new(LinearRing)
	lr.layout = layout
	lr.stride = layout.Stride()
	lr.flatCoords = flatCoords
	return lr
}

/**
*------------------------------
*				LinearRing（线环）相关的方法
*---------------------------------
*/
// Area方法 返回线环的面积
func (lr *LinearRing) Area() float64 {
	return doubleArea1(lr.flatCoords, 0, len(lr.flatCoords), lr.stride) / 2
}

// Clone方法 深层拷贝一个LinearRing
func (lr *LinearRing) Clone() *LinearRing {
	return deriveCloneLinearRing(lr)
}

// Empty方法 返回False
func (lr *LinearRing) Empty() bool {
	return false
}

// Length方法 返回线环的周长
func (lr *LinearRing) Length() float64 {
	return length1(lr.flatCoords, 0, len(lr.flatCoords), lr.stride)
}

// MustSetCoords方法 设置点坐标，任何错误都将抛出
func (lr *LinearRing) MustSetCoords(coords []Coord) *LinearRing {
	Must(lr.SetCoords(coords))
	return lr
}

// SetCoords方法 设置坐标
func (lr *LinearRing) SetCoords(coords []Coord) (*LinearRing, error) {
	if err := lr.setCoords(coords); err != nil {
		return nil, err
	}
	return lr, nil
}

// SetSRID方法 设置Linestring的坐标系参考
func (lr *LinearRing) SetSRID(srid int) *LinearRing {
	lr.srid = srid
	return lr
}

// Swap方法 将本对象与传入的LinearRing互换
func (lr *LinearRing) Swap(lr2 *LinearRing) {
	*lr, *lr2 = *lr2, *lr
}
