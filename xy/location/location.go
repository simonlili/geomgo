package location

import (
	"fmt"
)

// Type 枚举了不同拓扑位置，可能是在{@link Geometry}.
// 常数是用行和列指数DE-9IM {@link IntersectionMatrix}es.
type Type int

const (
	// Interior 表明在几何体内部.
	// Also, DE-9IM row index of the interior of the first geometry and column index of
	// the interior of the second geometry.
	Interior Type = iota
	// Boundary 表示在几何体的边界.
	// Also, DE-9IM row index of the boundary of the first geometry and column index of
	// the boundary of the second geometry.
	Boundary
	// Exterior 表示在几何体外部.
	// Also, DE-9IM row index of the exterior of the first geometry and column index of
	// the exterior of the second geometry.
	Exterior
	// None 表示未知的拓扑关系.
	None
)

func (t Type) String() string {

	switch t {
	case Exterior:
		return "Exterior"
	case Boundary:
		return "Boundary"
	case Interior:
		return "Interior"
	case None:
		return "None"
	}

	panic(fmt.Sprintf("Unknown location value: %v", int(t)))
}

// Symbol方法 转换拓扑关系的表示方式, for example, Exterior => 'e'
// locationValue
// Returns either 'e', 'b', 'i' or '-'
func (t Type) Symbol() rune {
	switch t {
	case Exterior:
		return 'e'
	case Boundary:
		return 'b'
	case Interior:
		return 'i'
	case None:
		return '-'
	}
	panic(fmt.Sprintf("Unknown location value: %v", int(t)))
}
