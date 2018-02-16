package geom

// GeometryCollection是一个具有相同SRID的任意几何类型的集合
type GeometryCollection struct {
	geoms []T
	srid  int
}

// NewGeometryCollection函数 创建一个明确规定的GeometryCollection
func NewGeometryCollection() *GeometryCollection {
	return &GeometryCollection{}
}

/**
*------------------------------
*				GeometryCollection（几何图像集合）相关的方法
*---------------------------------
*/
// Geom方法 返回指定索引的几何图像
func (gc *GeometryCollection) Geom(i int) T {
	return gc.geoms[i]
}

// Geoms方法 返回所有的几何图像
func (gc *GeometryCollection) Geoms() []T {
	return gc.geoms
}

// Layout方法 返回最小的视图，在所有几何图形中
func (gc *GeometryCollection) Layout() Layout {
	maxLayout := NoLayout
	for _, g := range gc.geoms {
		switch l := g.Layout(); l {
		case XYZ:
			if maxLayout == XYM {
				maxLayout = XYZM
			} else if l > maxLayout {
				maxLayout = l
			}
		case XYM:
			if maxLayout == XYZ {
				maxLayout = XYZM
			} else if l > maxLayout {
				maxLayout = l
			}
		default:
			if l > maxLayout {
				maxLayout = l
			}
		}
	}
	return maxLayout
}

// NumGeoms方法 返回gc中几何图形的数目
func (gc *GeometryCollection) NumGeoms() int {
	return len(gc.geoms)
}

// Stride方法 返回gc中图形的视图维数
func (gc *GeometryCollection) Stride() int {
	return gc.Layout().Stride()
}

// Bounds returns the bounds of all the geometries in gc.
func (gc *GeometryCollection) Bounds() *Bounds {
	// FIXME this needs work for mixing layouts, e.g. XYZ and XYM
	b := NewBounds(gc.Layout())
	for _, g := range gc.geoms {
		b = b.Extend(g)
	}
	return b
}

// Empty方法 检测集合是否为空，为空时返回true
func (gc *GeometryCollection) Empty() bool {
	return len(gc.geoms) == 0
}

// FlatCoords方法 坐标报错
func (*GeometryCollection) FlatCoords() []float64 {
	panic("FlatCoords() called on a GeometryCollection")
}

// Ends panics.
func (*GeometryCollection) Ends() []int {
	panic("Ends() called on a GeometryCollection")
}

// Endss panics.
func (*GeometryCollection) Endss() [][]int {
	panic("Endss() called on a GeometryCollection")
}

// SRID方法 获取gc的坐标系参考
func (gc *GeometryCollection) SRID() int {
	return gc.srid
}

// MustPush方法 向集合中添加几何图形，如果有错误均将抛出
func (gc *GeometryCollection) MustPush(gs ...T) *GeometryCollection {
	if err := gc.Push(gs...); err != nil {
		panic(err)
	}
	return gc
}

// Push方法 向集合中添加几何元素
func (gc *GeometryCollection) Push(gs ...T) error {
	gc.geoms = append(gc.geoms, gs...)
	return nil
}

// SetSRID方法 设置集合的坐标系参考，这个是集合的基础属性
func (gc *GeometryCollection) SetSRID(srid int) *GeometryCollection {
	gc.srid = srid
	return gc
}
