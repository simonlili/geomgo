package xy

import (
	"github.com/chengxiaoer/geomGo"
	"math"
	"github.com/chengxiaoer/geomGo/encoding/wkb"
)

//点线关系的函数

//获取线的控制点中距离某点最近的点的索引
// FIXME 第一个参数可以是geom.LineString
func PointIndexOnLine(ls wkb.LineString, coord geom.Coord) int {
	//获取线的所有控制点
	var min = math.MaxFloat32

	var index int

	for key, controllerPoint := range ls.Coords() {
		//计算点距离
		distance := Distance(controllerPoint, coord)
		//获取最小值的点所在索引
		if distance <= min {
			min = distance
			index = key
		}
	}
	return index
}
