package svggen

import (
	"reflect"
	"testing"
)

func TestSeq(t *testing.T) {
	testCases := []struct {
		name     string
		start    int
		end      int
		expected []int
	}{
		{
			name:     "Normal range",
			start:    1,
			end:      5,
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "Zero start",
			start:    0,
			end:      3,
			expected: []int{0, 1, 2, 3},
		},
		{
			name:     "Negative start",
			start:    -2,
			end:      2,
			expected: []int{-2, -1, 0, 1, 2},
		},
		{
			name:     "Single element",
			start:    4,
			end:      4,
			expected: []int{4},
		},
		{
			name:     "Empty range",
			start:    5,
			end:      4,
			expected: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := seq(tc.start, tc.end)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("seq(%d, %d) = %v; want %v", tc.start, tc.end, result, tc.expected)
			}
		})
	}
}

func TestMod(t *testing.T) {
	testCases := []struct {
		name        string
		a, b        int
		expected    int
		expectPanic bool
	}{
		{
			name:        "Normal case",
			a:           10,
			b:           3,
			expected:    1,
			expectPanic: false,
		},
		{
			name:        "Zero divisor",
			a:           5,
			b:           0,
			expected:    0, // Go panics on division by zero
			expectPanic: true,
		},
		{
			name:        "Negative dividend",
			a:           -7,
			b:           3,
			expected:    -1, // Follows the sign of the dividend
			expectPanic: false,
		},
		{
			name:        "Negative divisor",
			a:           7,
			b:           -3,
			expected:    1, // Result is positive as the dividend is positive
			expectPanic: false,
		},
		{
			name:        "Both negative",
			a:           -7,
			b:           -3,
			expected:    -1, // Follows the sign of the dividend
			expectPanic: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("mod(%d, %d) did not panic; wanted panic", tc.a, tc.b)
					}
				}()
			}

			result := mod(tc.a, tc.b)
			if !tc.expectPanic && result != tc.expected {
				t.Errorf("mod(%d, %d) = %d; want %d", tc.a, tc.b, result, tc.expected)
			}
		})
	}
}

func TestDiv(t *testing.T) {
	testCases := []struct {
		name        string
		a, b        int
		expected    int
		expectPanic bool
	}{
		{
			name:        "Normal division",
			a:           10,
			b:           2,
			expected:    5,
			expectPanic: false,
		},
		{
			name:        "Division by zero",
			a:           5,
			b:           0,
			expected:    0, // This value is irrelevant as we expect a panic
			expectPanic: true,
		},
		{
			name:        "Negative dividend",
			a:           -10,
			b:           2,
			expected:    -5,
			expectPanic: false,
		},
		{
			name:        "Negative divisor",
			a:           10,
			b:           -2,
			expected:    -5,
			expectPanic: false,
		},
		{
			name:        "Both negative",
			a:           -10,
			b:           -2,
			expected:    5,
			expectPanic: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("div(%d, %d) did not panic; wanted panic", tc.a, tc.b)
					}
				}()
			}

			result := div(tc.a, tc.b)
			if !tc.expectPanic && result != tc.expected {
				t.Errorf("div(%d, %d) = %d; want %d", tc.a, tc.b, result, tc.expected)
			}
		})
	}
}

func TestMultInt(t *testing.T) {
	testCases := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "Normal multiplication",
			a:        5,
			b:        4,
			expected: 20,
		},
		{
			name:     "Multiplication by zero",
			a:        5,
			b:        0,
			expected: 0,
		},
		{
			name:     "Zero multiplied by number",
			a:        0,
			b:        3,
			expected: 0,
		},
		{
			name:     "Negative number",
			a:        -5,
			b:        3,
			expected: -15,
		},
		{
			name:     "Both numbers negative",
			a:        -5,
			b:        -4,
			expected: 20,
		},
		{
			name:     "One negative, one positive",
			a:        -5,
			b:        4,
			expected: -20,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := multInt(tc.a, tc.b)
			if result != tc.expected {
				t.Errorf("mult(%d, %d) = %d; want %d", tc.a, tc.b, result, tc.expected)
			}
		})
	}
}

func TestMultFloat64(t *testing.T) {
	// Test cases
	testCases := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"MultiplyPositiveNumbers", 2.0, 3.0, 6.0},
		{"MultiplyPositiveDecimals", 0.5, 0.2, 0.1},
		{"MultiplyNegativeNumbers", -1.0, 4.0, -4.0},
		{"MultiplyByZero", 0.0, 100.0, 0.0},
		{"MultiplyNegativeDecimals", -2.5, -2.0, 5.0},
		{"MultiplyMixedSigns", 1.5, -1.5, -2.25},
	}

	// Iterate through test cases
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := multFloat64(testCase.a, testCase.b)
			if result != testCase.expected {
				t.Errorf("%s: multFloat64(%f, %f) = %f, expected %f", testCase.name, testCase.a, testCase.b, result, testCase.expected)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	testCases := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "Normal addition",
			a:        5,
			b:        3,
			expected: 8,
		},
		{
			name:     "Addition with zero",
			a:        5,
			b:        0,
			expected: 5,
		},
		{
			name:     "Zero addition",
			a:        0,
			b:        3,
			expected: 3,
		},
		{
			name:     "Negative number addition",
			a:        -5,
			b:        3,
			expected: -2,
		},
		{
			name:     "Both numbers negative",
			a:        -5,
			b:        -4,
			expected: -9,
		},
		{
			name:     "One negative, one positive",
			a:        -5,
			b:        5,
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := add(tc.a, tc.b)
			if result != tc.expected {
				t.Errorf("add(%d, %d) = %d; want %d", tc.a, tc.b, result, tc.expected)
			}
		})
	}
}

func TestHasElem(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []int
		elem     int
		expected bool
	}{
		{
			name:     "Element present",
			slice:    []int{1, 2, 3, 4, 5},
			elem:     3,
			expected: true,
		},
		{
			name:     "Element absent",
			slice:    []int{1, 2, 4, 5},
			elem:     3,
			expected: false,
		},
		{
			name:     "Empty slice",
			slice:    []int{},
			elem:     3,
			expected: false,
		},
		{
			name:     "Single element slice, present",
			slice:    []int{3},
			elem:     3,
			expected: true,
		},
		{
			name:     "Single element slice, absent",
			slice:    []int{1},
			elem:     3,
			expected: false,
		},
		{
			name:     "Negative element, present",
			slice:    []int{-1, -2, -3},
			elem:     -3,
			expected: true,
		},
		{
			name:     "Negative element, absent",
			slice:    []int{-1, -2, -4},
			elem:     -3,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := hasElem(tc.slice, tc.elem)
			if result != tc.expected {
				t.Errorf("hasElem(%v, %d) = %v; want %v", tc.slice, tc.elem, result, tc.expected)
			}
		})
	}
}

func TestClamp(t *testing.T) {
	testCases := []struct {
		name     string
		value    int
		min      int
		max      int
		expected int
	}{
		{
			name:     "Value within range",
			value:    50,
			min:      0,
			max:      100,
			expected: 50,
		},
		{
			name:     "Value below minimum",
			value:    -10,
			min:      0,
			max:      100,
			expected: 0,
		},
		{
			name:     "Value above maximum",
			value:    150,
			min:      0,
			max:      100,
			expected: 100,
		},
		{
			name:     "Value equal to minimum",
			value:    0,
			min:      0,
			max:      100,
			expected: 0,
		},
		{
			name:     "Value equal to maximum",
			value:    100,
			min:      0,
			max:      100,
			expected: 100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := clamp(tc.value, tc.min, tc.max)
			if result != tc.expected {
				t.Errorf("Test %s failed: expected %d, got %d", tc.name, tc.expected, result)
			}
		})
	}
}
