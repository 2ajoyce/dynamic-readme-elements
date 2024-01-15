package svggen

func seq(start, end int) []int {
	sequence := make([]int, end-start+1)
	for i := range sequence {
		sequence[i] = start + i
	}
	return sequence
}

func mod(a, b int) int {
	return a % b
}

func div(a, b int) int {
	return a / b
}

func multInt(a, b int) int {
	return a * b
}

func multFloat64(a, b float64) float64 {
	return a * b
}

func add(a, b int) int {
	return a + b
}

func hasElem(slice []int, elem int) bool {
	for _, v := range slice {
		if v == elem {
			return true
		}
	}
	return false
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
