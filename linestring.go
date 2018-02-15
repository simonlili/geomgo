package geom

// LineString 代表了单线类型
type LineString struct {
	geom1
}

// NewLineString 分配了一个没有控制点并且符合Layout的 Linestring
func NewLineString(l Layout) *LineString {
	return NewLineStringFlat(l, nil)
}

// NewLineStringFlat 分配了一个有控制点（flatCoords）并且符合Layout的 Linestring
func NewLineStringFlat(layout Layout, flatCoords []float64) *LineString {
	ls := new(LineString)
	ls.layout = layout
	ls.stride = layout.Stride()
	ls.flatCoords = flatCoords
	return ls
}

/**
*------------------------------
*				Linestring（点）相关的方法
*---------------------------------
*/

// Area方法 返回 Linestring的面积
func (ls *LineString) Area() float64 {
	return 0
}

// Clone方法 返回一个Linestring的拷贝，这个不是ls的别名
func (ls *LineString) Clone() *LineString {
	return deriveCloneLineString(ls)
}

// Empty方法 返回false
func (ls *LineString) Empty() bool {
	return false
}

// Interpolate returns the index and delta of val in dimension dim.
func (ls *LineString) Interpolate(val float64, dim int) (int, float64) {
	n := len(ls.flatCoords)
	if n == 0 {
		panic("geom: empty linestring")
	}
	if val <= ls.flatCoords[dim] {
		return 0, 0
	}
	if ls.flatCoords[n-ls.stride+dim] <= val {
		return (n - 1) / ls.stride, 0
	}
	low := 0
	high := n / ls.stride
	for low < high {
		mid := (low + high) / 2
		if val < ls.flatCoords[mid*ls.stride+dim] {
			high = mid
		} else {
			low = mid + 1
		}
	}
	low--
	val0 := ls.flatCoords[low*ls.stride+dim]
	if val == val0 {
		return low, 0
	}
	val1 := ls.flatCoords[(low+1)*ls.stride+dim]
	return low, (val - val0) / (val1 - val0)
}

// Length 返回 Linestring 的长度
func (ls *LineString) Length() float64 {
	return length1(ls.flatCoords, 0, len(ls.flatCoords), ls.stride)
}

// MustSetCoords方法 设置Linestring的控制点，但是有任何 ERROR 都会抛出
func (ls *LineString) MustSetCoords(coords []Coord) *LineString {
	Must(ls.SetCoords(coords))
	return ls
}

// SetCoords方法 为Linestring 设置控制点
func (ls *LineString) SetCoords(coords []Coord) (*LineString, error) {
	if err := ls.setCoords(coords); err != nil {
		return nil, err
	}
	return ls, nil
}

// SetSRID方法 设置Linestring的 坐标系参考
func (ls *LineString) SetSRID(srid int) *LineString {
	ls.srid = srid
	return ls
}

// SubLineString方法返回一个开始和结束点相同的 Linestring. 返回的Linestring 替代了 ls
func (ls *LineString) SubLineString(start, stop int) *LineString {
	return NewLineStringFlat(ls.layout, ls.flatCoords[start*ls.stride:stop*ls.stride])
}

// Swap方法 与参数传入的 Linestring 互换
func (ls *LineString) Swap(ls2 *LineString) {
	*ls, *ls2 = *ls2, *ls
}
