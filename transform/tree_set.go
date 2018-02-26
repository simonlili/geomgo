package transform

import (
	"fmt"

	"github.com/chengxiaoer/go-geom"
)

// Compare 比较两个两个坐标的大小和是否相同
type Compare interface {
	IsEquals(x, y geom.Coord) bool
	IsLess(x, y geom.Coord) bool
}

type tree struct {
	left  *tree
	value geom.Coord
	right *tree
}

// TreeSet 使用Compare里面的方法，对坐标进行排序。根据Equals函数移除坐标中的重复项
type TreeSet struct {
	compare Compare
	tree    *tree
	size    int
	layout  geom.Layout
	stride  int
}

// NewTreeSet函数 创建一个新的 TreeSet 实例。
func NewTreeSet(layout geom.Layout, compare Compare) *TreeSet {
	return &TreeSet{
		layout:  layout,
		stride:  layout.Stride(),
		compare: compare,
	}
}

/**
*------------------------------
*				TreeSet  相关的方法
*---------------------------------
*/
// Insert方法 向TreeSet中添加新的坐标
// 添加的坐标必须具有相同的尸体布局维数
// Returns true 如果点成功添加
// Returns false 如果添加的点在 TreeSet 中已经存在
func (set *TreeSet) Insert(coord geom.Coord) bool {
	if set.stride == 0 {
		set.stride = set.layout.Stride()
	}
	if len(coord) < set.stride {
		panic(fmt.Sprintf("Coordinate inserted into tree does not have a sufficient number of points for the provided layout.  Length of Coord was %v but should have been %v", len(coord), set.stride))
	}
	tree, added := set.insertImpl(set.tree, coord)
	if added {
		set.tree = tree
		set.size++
	}

	return added
}

// ToFlatArray方法 返回一个浮点数组包含treeSet中所有的坐标
func (set *TreeSet) ToFlatArray() []float64 {
	stride := set.layout.Stride()
	array := make([]float64, set.size*stride)

	i := 0
	set.walk(set.tree, func(v []float64) {
		for j := 0; j < stride; j++ {
			array[i+j] = v[j]
		}
		i += stride
	})

	return array
}

func (set *TreeSet) walk(t *tree, visitor func([]float64)) {
	if t == nil {
		return
	}
	set.walk(t.left, visitor)
	visitor(t.value)
	set.walk(t.right, visitor)
}

func (set *TreeSet) insertImpl(t *tree, v []float64) (*tree, bool) {
	if t == nil {
		return &tree{nil, v, nil}, true
	}

	if set.compare.IsEquals(geom.Coord(v), t.value) {
		return t, false
	}

	var added bool
	if set.compare.IsLess(geom.Coord(v), t.value) {
		t.left, added = set.insertImpl(t.left, v)
	} else {
		t.right, added = set.insertImpl(t.right, v)
	}

	return t, added
}
