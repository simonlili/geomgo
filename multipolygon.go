package geom

// MultiPolygon对象是一个多边形的集合
type MultiPolygon struct {
	geom3
}

// NewMultiPolygon函数 创建一个没有多边形的 MultiPolygon
func NewMultiPolygon(layout Layout) *MultiPolygon {
	return NewMultiPolygonFlat(layout, nil, nil)
}

// NewMultiPolygonFlat函数 根据传入参数构建一个非空的MultiPolygon
func NewMultiPolygonFlat(layout Layout, flatCoords []float64, endss [][]int) *MultiPolygon {
	mp := new(MultiPolygon)
	mp.layout = layout
	mp.stride = layout.Stride()
	mp.flatCoords = flatCoords
	mp.endss = endss
	return mp
}

/**
*------------------------------
*				MultiPolygon（多边形集合）相关的方法
*---------------------------------
*/

// Area方法 返回所有多边形面积之和
func (mp *MultiPolygon) Area() float64 {
	return doubleArea3(mp.flatCoords, 0, mp.endss, mp.stride) / 2
}

// Clone方法 创建一个深层拷贝.
func (mp *MultiPolygon) Clone() *MultiPolygon {
	return deriveCloneMultiPolygon(mp)
}

// Empty方法 检测集合是否为空，为空返回true
func (mp *MultiPolygon) Empty() bool {
	return mp.NumPolygons() == 0
}

// Length方法 返回所有多边形的周长之和
func (mp *MultiPolygon) Length() float64 {
	return length3(mp.flatCoords, 0, mp.endss, mp.stride)
}

// MustSetCoords方法 设置坐标，遇到任何错误都将抛出
func (mp *MultiPolygon) MustSetCoords(coords [][][]Coord) *MultiPolygon {
	Must(mp.SetCoords(coords))
	return mp
}

// NumPolygons方法 获取集合中多边行的数目
func (mp *MultiPolygon) NumPolygons() int {
	return len(mp.endss)
}

// Polygon方法 返回指定索引的多边形
func (mp *MultiPolygon) Polygon(i int) *Polygon {
	offset := 0
	if i > 0 {
		ends := mp.endss[i-1]
		offset = ends[len(ends)-1]
	}
	ends := make([]int, len(mp.endss[i]))
	if offset == 0 {
		copy(ends, mp.endss[i])
	} else {
		for j, end := range mp.endss[i] {
			ends[j] = end - offset
		}
	}
	return NewPolygonFlat(mp.layout, mp.flatCoords[offset:mp.endss[i][len(mp.endss[i])-1]], ends)
}

// Push方法 向集合中添加一个多边形.
func (mp *MultiPolygon) Push(p *Polygon) error {
	if p.layout != mp.layout {
		return ErrLayoutMismatch{Got: p.layout, Want: mp.layout}
	}
	offset := len(mp.flatCoords)
	ends := make([]int, len(p.ends))
	if offset == 0 {
		copy(ends, p.ends)
	} else {
		for i, end := range p.ends {
			ends[i] = end + offset
		}
	}
	mp.flatCoords = append(mp.flatCoords, p.flatCoords...)
	mp.endss = append(mp.endss, ends)
	return nil
}

// SetCoords方法 设置对象的坐标
func (mp *MultiPolygon) SetCoords(coords [][][]Coord) (*MultiPolygon, error) {
	if err := mp.setCoords(coords); err != nil {
		return nil, err
	}
	return mp, nil
}

// SetSRID方法 设置对象的坐标系参考
func (mp *MultiPolygon) SetSRID(srid int) *MultiPolygon {
	mp.srid = srid
	return mp
}

// Swap方法 将本对象与传入的多边形集合对象互相交换
func (mp *MultiPolygon) Swap(mp2 *MultiPolygon) {
	*mp, *mp2 = *mp2, *mp
}
