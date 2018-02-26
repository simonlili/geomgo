package internal

// IsSameSignAndNonZero函数 检查参数a 和 b 的正负性是否一直
func IsSameSignAndNonZero(a, b float64) bool {
	if a == 0 || b == 0 {
		return false
	}
	return (a < 0 && b < 0) || (a > 0 && b > 0)
}

// Min函数 返回传入的四个参数中的最小值
func Min(v1, v2, v3, v4 float64) float64 {
	min := v1
	if v2 < min {
		min = v2
	}
	if v3 < min {
		min = v3
	}
	if v4 < min {
		min = v4
	}
	return min
}
