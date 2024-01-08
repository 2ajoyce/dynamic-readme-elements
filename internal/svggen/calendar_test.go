package svggen

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
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

func TestMult(t *testing.T) {
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
			result := mult(tc.a, tc.b)
			if result != tc.expected {
				t.Errorf("mult(%d, %d) = %d; want %d", tc.a, tc.b, result, tc.expected)
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

func TestHandleCalendar(t *testing.T) {
	router := gin.Default()
	router.GET("/calendar", HandleCalendar)

	testCases := []struct {
		name           string
		queryString    string
		expectedStatus int
		expectInBody   []string
	}{
		{
			name:           "Valid parameters - Spot Checks",
			queryString:    "/calendar?year=2023&month=1&progressDays=1,15",
			expectedStatus: http.StatusOK,
			expectInBody: []string{
				`<svg width="370px" height="310px" xmlns="http://www.w3.org/2000/svg" font-family="Arial">`,
				`<rect x="5" y="5" width="360px" height="300px" fill="white" rx="15" />`,                    // Background
				`<text x="180" y="35" font-size="20" text-anchor="middle" fill="black">January 2023</text>`, // Header
				// Check for a few specific days, including one that is marked as progress and one that is not
				`<rect x="15" y="45" width="40" height="40" fill="#4c1" stroke="#ddd" />`,       // Day 1 with progress
				`<text x="35" y="70" font-size="14" text-anchor="middle" fill="white">1</text>`, // Text for Day 1
				`<rect x="65" y="45" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,    // Day 2 without progress
				`<text x="85" y="70" font-size="14" text-anchor="middle" fill="black">2</text>`, // Text for Day 2
			},
		},
		{
			name:           "Valid parameters - Check Full Response",
			queryString:    "/calendar?year=2023&month=1&progressDays=1,15",
			expectedStatus: http.StatusOK,
			expectInBody: []string{
				`<svg width="370px" height="310px" xmlns="http://www.w3.org/2000/svg" font-family="Arial">`,
				`<rect x="5" y="5" width="360px" height="300px" fill="white" rx="15" />`,
				`<text x="180" y="35" font-size="20" text-anchor="middle" fill="black">January 2023</text><rect x="15" y="45" width="40" height="40" fill="#4c1" stroke="#ddd" />`,
				`<text x="35" y="70" font-size="14" text-anchor="middle" fill="white">1</text><rect x="65" y="45" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="85" y="70" font-size="14" text-anchor="middle" fill="black">2</text><rect x="115" y="45" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="135" y="70" font-size="14" text-anchor="middle" fill="black">3</text><rect x="165" y="45" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="185" y="70" font-size="14" text-anchor="middle" fill="black">4</text><rect x="215" y="45" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="235" y="70" font-size="14" text-anchor="middle" fill="black">5</text><rect x="265" y="45" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="285" y="70" font-size="14" text-anchor="middle" fill="black">6</text><rect x="315" y="45" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="335" y="70" font-size="14" text-anchor="middle" fill="black">7</text><rect x="15" y="95" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="35" y="120" font-size="14" text-anchor="middle" fill="black">8</text><rect x="65" y="95" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="85" y="120" font-size="14" text-anchor="middle" fill="black">9</text><rect x="115" y="95" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="135" y="120" font-size="14" text-anchor="middle" fill="black">10</text><rect x="165" y="95" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="185" y="120" font-size="14" text-anchor="middle" fill="black">11</text><rect x="215" y="95" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="235" y="120" font-size="14" text-anchor="middle" fill="black">12</text><rect x="265" y="95" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="285" y="120" font-size="14" text-anchor="middle" fill="black">13</text><rect x="315" y="95" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="335" y="120" font-size="14" text-anchor="middle" fill="black">14</text><rect x="15" y="145" width="40" height="40" fill="#4c1" stroke="#ddd" />`,
				`<text x="35" y="170" font-size="14" text-anchor="middle" fill="white">15</text><rect x="65" y="145" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="85" y="170" font-size="14" text-anchor="middle" fill="black">16</text><rect x="115" y="145" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="135" y="170" font-size="14" text-anchor="middle" fill="black">17</text><rect x="165" y="145" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="185" y="170" font-size="14" text-anchor="middle" fill="black">18</text><rect x="215" y="145" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="235" y="170" font-size="14" text-anchor="middle" fill="black">19</text><rect x="265" y="145" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="285" y="170" font-size="14" text-anchor="middle" fill="black">20</text><rect x="315" y="145" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="335" y="170" font-size="14" text-anchor="middle" fill="black">21</text><rect x="15" y="195" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="35" y="220" font-size="14" text-anchor="middle" fill="black">22</text><rect x="65" y="195" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="85" y="220" font-size="14" text-anchor="middle" fill="black">23</text><rect x="115" y="195" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="135" y="220" font-size="14" text-anchor="middle" fill="black">24</text><rect x="165" y="195" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="185" y="220" font-size="14" text-anchor="middle" fill="black">25</text><rect x="215" y="195" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="235" y="220" font-size="14" text-anchor="middle" fill="black">26</text><rect x="265" y="195" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="285" y="220" font-size="14" text-anchor="middle" fill="black">27</text><rect x="315" y="195" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="335" y="220" font-size="14" text-anchor="middle" fill="black">28</text><rect x="15" y="245" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="35" y="270" font-size="14" text-anchor="middle" fill="black">29</text><rect x="65" y="245" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="85" y="270" font-size="14" text-anchor="middle" fill="black">30</text><rect x="115" y="245" width="40" height="40" fill="#f0f0f0" stroke="#ddd" />`,
				`<text x="135" y="270" font-size="14" text-anchor="middle" fill="black">31</text>`,
				`</svg>`,
			},
		},
		{
			name:           "Invalid year",
			queryString:    "/calendar?year=abc&month=1",
			expectedStatus: http.StatusBadRequest,
			expectInBody:   []string{"Invalid year format"},
		},
		{
			name:           "Invalid month",
			queryString:    "/calendar?year=2023&month=13",
			expectedStatus: http.StatusBadRequest,
			expectInBody:   []string{"Month must be between 1 and 12"},
		},
		{
			name:           "Default to current year and month",
			queryString:    "/calendar",
			expectedStatus: http.StatusOK,
			expectInBody:   []string{time.Now().Format("January"), strconv.Itoa(time.Now().Year())},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tc.queryString, nil)
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			body := w.Body.String()

			for _, str := range tc.expectInBody {
				if !strings.Contains(body, str) {
					t.Errorf("Expected to find %s in response body", str)
				}
			}
		})
	}
}
