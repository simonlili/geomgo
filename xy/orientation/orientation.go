package orientation

import "fmt"

// Type 枚举点与向量之间的角度关系.
type Type int

const (
	// Clockwise 表明向量或点相对于参照向量是顺时针方向的。
	Clockwise Type = iota - 1
	// Collinear 表示向量或点与参照向量沿着同一方向。
	Collinear
	// CounterClockwise 表明向量或点相对于参照向量是逆时针方向的
	CounterClockwise
)

var orientationLabels = [3]string{"Clockwise", "Collinear", "CounterClockwise"}

func (o Type) String() string {
	if o > 1 || o < -1 {
		return fmt.Sprintf("Unsafe to calculate: %v", int(o))
	}
	return orientationLabels[int(o+1)]
}
