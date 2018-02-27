package internal

import (
	"github.com/chengxiaoer/geomGo"
)

// CoordStack 是一个存放坐标的，简单的栈(in []float64 form)可存可取
// 这些坐标按照正常堆栈的顺序进行返回
//必须使用 NewCoordStack函数来创建
type CoordStack struct {
	// Data 是栈的数据.  遵循先进后出、后进先出的规则
	Data   []float64
	stride int
}

// NewCoordStack 根据视图类型创建一个栈
func NewCoordStack(layout geom.Layout) *CoordStack {
	return &CoordStack{stride: layout.Stride()}
}

// Push方法 存入一个坐标在栈的指定位置上.
func (stack *CoordStack) Push(data []float64, idx int) []float64 {
	c := data[idx : idx+stack.stride]
	stack.Data = append(stack.Data, c...)
	return c
}

// Pop方法 弹出栈顶存放的坐标
func (stack *CoordStack) Pop() ([]float64, int) {
	numOrds := len(stack.Data)
	start := numOrds - stack.stride
	coord := stack.Data[start:numOrds]
	stack.Data = stack.Data[:start]
	return coord, stack.Size()
}

// Peek方法 返回最新存入的坐标，不改变栈的结构
func (stack *CoordStack) Peek() []float64 {
	numOrds := len(stack.Data)
	start := numOrds - stack.stride
	coord := stack.Data[start:numOrds]
	return coord
}

// Size方法 返回栈中坐标的个数
func (stack *CoordStack) Size() int {
	return len(stack.Data) / stack.stride
}
