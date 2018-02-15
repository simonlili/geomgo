package geom

//Polygon对象 多变形对象是一个LinearRing(线环)的集合。第一个LinearRing作为外边界，
//随后的LinearRing对象作为内边界
type Polygon struct {
	geom2
}

// NewPolygon函数 创建一个空的多边形
func NewPolygon(layout Layout) *Polygon {
	return NewPolygonFlat(layout, nil, nil)
}

// NewPolygonFlat函数 根据传入的坐标和视图类型创建多边形
func NewPolygonFlat(layout Layout, flatCoords []float64, ends []int) *Polygon {
	p := new(Polygon)
	p.layout = layout
	p.stride = layout.Stride()
	p.flatCoords = flatCoords
	p.ends = ends
	return p
}

/**
*------------------------------
*				Polygon（多边形）相关的方法
*---------------------------------
*/

// Area方法 返回多边形的面积
func (p *Polygon) Area() float64 {
	return doubleArea2(p.flatCoords, 0, p.ends, p.stride) / 2
}

// Clone方法 深层拷贝多边形
func (p *Polygon) Clone() *Polygon {
	return deriveClonePolygon(p)
}

// Empty方法 返回False
func (p *Polygon) Empty() bool {
	return false
}

// Length方法 返回周长
func (p *Polygon) Length() float64 {
	return length2(p.flatCoords, 0, p.ends, p.stride)
}

// LinearRing方法 获取指定索引的线环
func (p *Polygon) LinearRing(i int) *LinearRing {
	offset := 0
	if i > 0 {
		offset = p.ends[i-1]
	}
	return NewLinearRingFlat(p.layout, p.flatCoords[offset:p.ends[i]])
}

// MustSetCoords方法 设置坐标，任何错误都将抛出
func (p *Polygon) MustSetCoords(coords [][]Coord) *Polygon {
	Must(p.SetCoords(coords))
	return p
}

// NumLinearRings方法 返回LinearRing的数目
func (p *Polygon) NumLinearRings() int {
	return len(p.ends)
}

// Push方法 向多边形中添加LinearRing
func (p *Polygon) Push(lr *LinearRing) error {
	if lr.layout != p.layout {
		return ErrLayoutMismatch{Got: lr.layout, Want: p.layout}
	}
	p.flatCoords = append(p.flatCoords, lr.flatCoords...)
	p.ends = append(p.ends, len(p.flatCoords))
	return nil
}

// SetCoords方法 设置坐标
func (p *Polygon) SetCoords(coords [][]Coord) (*Polygon, error) {
	if err := p.setCoords(coords); err != nil {
		return nil, err
	}
	return p, nil
}

// SetSRID方法 设置多变形的坐标系参考
func (p *Polygon) SetSRID(srid int) *Polygon {
	p.srid = srid
	return p
}

// Swap方法 将本对象与传入的多边形对象互相交换
func (p *Polygon) Swap(p2 *Polygon) {
	*p, *p2 = *p2, *p
}
