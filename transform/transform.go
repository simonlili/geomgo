package transform

import "github.com/chengxiaoer/geomGo"

// UniqueCoords函数 一个新的坐标数组 (具有与输入相同的视图类型)，在CoordData中包含了不同且唯一的坐标
// 坐标的顺序与输入的顺序相同
func UniqueCoords(layout geom.Layout, compare Compare, coordData []float64) []float64 {
	set := NewTreeSet(layout, compare)
	stride := layout.Stride()
	uniqueCoords := make([]float64, 0, len(coordData))
	numCoordsAdded := 0
	for i := 0; i < len(coordData); i += stride {
		coord := coordData[i : i+stride]
		added := set.Insert(geom.Coord(coord))

		if added {
			uniqueCoords = append(uniqueCoords, coord...)
			numCoordsAdded++
		}
	}
	return uniqueCoords[:numCoordsAdded*stride]
}
