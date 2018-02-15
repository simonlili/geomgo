package geom

//MultiLineString 是 LineStrings的集合.
type MultiLineString struct {
	geom2
}

// NewMultiLineString函数返回一个没有 Linestring 的MultiLinestring
func NewMultiLineString(layout Layout) *MultiLineString {
	return NewMultiLineStringFlat(layout, nil, nil)
}

// NewMultiLineStringFlat函数 返回一个新的、有Linestring的MultiLinestring.
func NewMultiLineStringFlat(layout Layout, flatCoords []float64, ends []int) *MultiLineString {
	mls := new(MultiLineString)
	mls.layout = layout
	mls.stride = layout.Stride()
	mls.flatCoords = flatCoords
	mls.ends = ends
	return mls
}

/**
*------------------------------
*				MultiPoint（多点）相关的方法
*---------------------------------
*/

// Area方法 返回0
func (mls *MultiLineString) Area() float64 {
	return 0
}

// Clone方法 对MultiPoint进行深层拷贝
func (mls *MultiLineString) Clone() *MultiLineString {
	return deriveCloneMultiLineString(mls)
}

// Empty方法 如果此集合为空则返回true
func (mls *MultiLineString) Empty() bool {
	return mls.NumLineStrings() == 0
}

// Length方法 返回集合中 Linestring 长度的和
func (mls *MultiLineString) Length() float64 {
	return length2(mls.flatCoords, 0, mls.ends, mls.stride)
}

// LineString方法 返回指定索引的 Lienstring
func (mls *MultiLineString) LineString(i int) *LineString {
	offset := 0
	if i > 0 {
		offset = mls.ends[i-1]
	}
	return NewLineStringFlat(mls.layout, mls.flatCoords[offset:mls.ends[i]])
}

// MustSetCoords方法 设置坐标，如果有错将会抛出
func (mls *MultiLineString) MustSetCoords(coords [][]Coord) *MultiLineString {
	Must(mls.SetCoords(coords))
	return mls
}

// NumLineStrings方法 返回集合中Lienstring的个数
func (mls *MultiLineString) NumLineStrings() int {
	return len(mls.ends)
}

// Push方法 向集合中添加一个Linestring
func (mls *MultiLineString) Push(ls *LineString) error {
	if ls.layout != mls.layout {
		return ErrLayoutMismatch{Got: ls.layout, Want: mls.layout}
	}
	mls.flatCoords = append(mls.flatCoords, ls.flatCoords...)
	mls.ends = append(mls.ends, len(mls.flatCoords))
	return nil
}

// SetCoords方法 设置坐标
func (mls *MultiLineString) SetCoords(coords [][]Coord) (*MultiLineString, error) {
	if err := mls.setCoords(coords); err != nil {
		return nil, err
	}
	return mls, nil
}

// SetSRID方法 设置MultiLinestring的坐标系参考
func (mls *MultiLineString) SetSRID(srid int) *MultiLineString {
	mls.srid = srid
	return mls
}

// Swap方法 将本集合与传入的集合互相交换
func (mls *MultiLineString) Swap(mls2 *MultiLineString) {
	*mls, *mls2 = *mls2, *mls
}
