package svggen

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCalculateGridSize(t *testing.T) {
	testCases := []struct {
		name            string
		width           int
		numberOfSquares int
		expectedRows    int
		expectedColumns int
	}{
		{
			name:            "Standard case",
			width:           300,
			numberOfSquares: 100,
			expectedRows:    10,
			expectedColumns: 10,
		},
		{
			name:            "Width less than ideal square size",
			width:           5,
			numberOfSquares: 25,
			expectedRows:    1,
			expectedColumns: 25,
		},
		{
			name:            "Width leading to fractional squares",
			width:           300,
			numberOfSquares: 105,
			expectedRows:    10,
			expectedColumns: 11,
		},
		{
			name:            "Minimum width",
			width:           1,
			numberOfSquares: 10,
			expectedRows:    1,
			expectedColumns: 10,
		},
		{
			name:            "Large number of squares",
			width:           300,
			numberOfSquares: 1000,
			expectedRows:    21,
			expectedColumns: 48,
		},
		{
			name:            "Width exactly divisible by square size",
			width:           300,
			numberOfSquares: 90, // 30x30 squares
			expectedRows:    9,
			expectedColumns: 10,
		},
		{
			name:            "Odd number of squares",
			width:           300,
			numberOfSquares: 99,
			expectedRows:    9,
			expectedColumns: 11,
		},
		{
			name:            "Zero width",
			width:           0,
			numberOfSquares: 100,
			expectedRows:    7,
			expectedColumns: 15,
		},
		{
			name:            "Negative width",
			width:           -300,
			numberOfSquares: 100,
			expectedRows:    7,
			expectedColumns: 15,
		}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rows, columns := CalculateGridSize(tc.width, tc.numberOfSquares, 3)

			if rows != tc.expectedRows || columns != tc.expectedColumns {
				t.Errorf("Test %s failed: expected (%d, %d), got (%d, %d)", tc.name, tc.expectedRows, tc.expectedColumns, rows, columns)
			}
		})
	}
}

func TestGenerateSquares(t *testing.T) {
	testCases := []struct {
		name            string
		squaresPerRow   int
		width           int
		numberOfSquares int
		filledSquares   int
		expectedSquares []struct {
			X, Y  int
			Color string
		}
	}{
		{
			name:            "Small mix case",
			squaresPerRow:   5,
			width:           100,
			numberOfSquares: 10,
			filledSquares:   6,
			expectedSquares: []struct {
				X, Y  int
				Color string
			}{
				{3, 3, Colors.Green},
				{23, 3, Colors.Green},
				{43, 3, Colors.Green},
				{63, 3, Colors.Green},
				{83, 3, Colors.Green},
				{3, 23, Colors.Green},
				{23, 23, Colors.Grey},
				{43, 23, Colors.Grey},
				{63, 23, Colors.Grey},
				{83, 23, Colors.Grey},
			},
		},
		{
			name:            "All squares filled",
			squaresPerRow:   5,
			width:           100,
			numberOfSquares: 10,
			filledSquares:   10,
			expectedSquares: []struct {
				X, Y  int
				Color string
			}{
				{3, 3, Colors.Green},
				{23, 3, Colors.Green},
				{43, 3, Colors.Green},
				{63, 3, Colors.Green},
				{83, 3, Colors.Green},
				{3, 23, Colors.Green},
				{23, 23, Colors.Green},
				{43, 23, Colors.Green},
				{63, 23, Colors.Green},
				{83, 23, Colors.Green},
			},
		},
		{
			name:            "No squares filled",
			squaresPerRow:   5,
			width:           100,
			numberOfSquares: 10,
			filledSquares:   0,
			expectedSquares: []struct {
				X, Y  int
				Color string
			}{
				{3, 3, Colors.Grey},
				{23, 3, Colors.Grey},
				{43, 3, Colors.Grey},
				{63, 3, Colors.Grey},
				{83, 3, Colors.Grey},
				{3, 23, Colors.Grey},
				{23, 23, Colors.Grey},
				{43, 23, Colors.Grey},
				{63, 23, Colors.Grey},
				{83, 23, Colors.Grey},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			squares := GenerateSquares(tc.width, tc.squaresPerRow, tc.numberOfSquares, tc.filledSquares, 3)

			// Validate the number of squares generated
			if len(squares) != tc.numberOfSquares {
				t.Errorf("Test %s failed: expected number of squares %d, got %d", tc.name, tc.numberOfSquares, len(squares))
			}

			// Validate the position and color of each square
			for i, expectedSquare := range tc.expectedSquares {
				if squares[i].X != expectedSquare.X || squares[i].Y != expectedSquare.Y || squares[i].Color != expectedSquare.Color {
					t.Errorf("Test %s failed at square %d: expected %v, got %v", tc.name, i, expectedSquare, squares[i])
				}
			}
		})
	}
}

func TestHandleProgressWaffle(t *testing.T) {
	// Setup Gin with test mode
	gin.SetMode(gin.TestMode)

	// Create router and add route
	router := gin.Default()
	router.GET("/progress/waffle", HandleProgressWaffle)

	testCases := []struct {
		name             string
		query            string
		expectedStatus   int
		expectedFilled   int
		expectedUnfilled int
	}{
		{
			name:             "Valid request",
			query:            "/progress/waffle?width=300&numberOfSquares=100&percentage=50",
			expectedStatus:   http.StatusOK,
			expectedFilled:   50,
			expectedUnfilled: 50,
		},
		{
			name:             "Width only",
			query:            "/progress/waffle?width=200&percentage=25",
			expectedStatus:   http.StatusOK,
			expectedFilled:   25,
			expectedUnfilled: 75,
		},
		{
			name:             "Invalid percentage",
			query:            "/progress/waffle?width=200&numberOfSquares=100&percentage=150",
			expectedStatus:   http.StatusOK,
			expectedFilled:   100, // Clamped to 100%
			expectedUnfilled: 0,
		},
		{
			name:             "Zero squares",
			query:            "/progress/waffle?width=200&numberOfSquares=0&percentage=50",
			expectedStatus:   http.StatusOK,
			expectedFilled:   0,
			expectedUnfilled: 0,
		},
		{
			name:             "Negative width",
			query:            "/progress/waffle?width=-200&numberOfSquares=100&percentage=50",
			expectedStatus:   http.StatusOK,
			expectedFilled:   50,
			expectedUnfilled: 50,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions for status code
			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d for test %s", tc.expectedStatus, w.Code, tc.name)
			}

			// Validate the response body
			responseBody := w.Body.String()
			if err := validateWaffleSVG(responseBody, tc.expectedFilled, tc.expectedUnfilled); err != nil {
				t.Errorf("Validation failed for test %s: %v", tc.name, err)
			}
		})
	}
}

func validateWaffleSVG(svgContent string, filledCount, unfilledCount int) error {
	// Check basic SVG structure
	if !strings.Contains(svgContent, "<svg") || !strings.Contains(svgContent, "</svg>") {
		return errors.New("missing SVG tags")
	}

	// Use regular expressions to match squares regardless of attribute order
	filledRegex := regexp.MustCompile(fmt.Sprintf(`class="gridSquare"[^>]*fill="%s"`, Colors.Green))
	unfilledRegex := regexp.MustCompile(fmt.Sprintf(`class="gridSquare"[^>]*fill="%s"`, Colors.Grey))

	// Count occurrences of filled and unfilled squares
	actualFilled := len(filledRegex.FindAllStringIndex(svgContent, -1))
	actualUnfilled := len(unfilledRegex.FindAllStringIndex(svgContent, -1))

	if actualFilled != filledCount || actualUnfilled != unfilledCount {
		return fmt.Errorf("incorrect count of filled (%d) or unfilled (%d) squares, expected filled: %d, expected unfilled: %d", actualFilled, actualUnfilled, filledCount, unfilledCount)
	}

	return nil
}
