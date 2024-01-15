package svggen

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRadians(t *testing.T) {
	tests := []struct {
		degrees         float64
		expectedRadians float64
	}{
		{0.0, 0.0},
		{45.0, math.Pi / 4},
		{90.0, math.Pi / 2},
		{180.0, math.Pi},
		{360.0, 2 * math.Pi},
		{-45.0, -math.Pi / 4},
	}

	for _, test := range tests {
		result := radians(test.degrees)
		if result != test.expectedRadians {
			t.Errorf("radians(%f) = %f; expected %f", test.degrees, result, test.expectedRadians)
		}
	}
}

func TestHandleProgressGauge(t *testing.T) {
	tests := []struct {
		widthParam      string
		percentageParam string
		expectedSVG     string
	}{
		{"100", "50", "<svg height=\"50px\" width=\"100px\" viewBox=\"0 0 100 50\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M 0.000000,50.000000 A 50.000000,50.000000 0 0 1 9.549150,20.610737 L 50.000000,50.000000 L 0.000000,50.000000 Z\" fill=\"red\"/><path d=\"M 9.549150,20.610737 A 50.000000,50.000000 0 0 1 34.549150,2.447174 L 50.000000,50.000000 L 9.549150,20.610737 Z\" fill=\"orange\"/><path d=\"M 34.549150,2.447174 A 50.000000,50.000000 0 0 1 65.450850,2.447174 L 50.000000,50.000000 L 34.549150,2.447174 Z\" fill=\"yellow\"/><path d=\"M 65.450850,2.447174 A 50.000000,50.000000 0 0 1 90.450850,20.610737 L 50.000000,50.000000 L 65.450850,2.447174 Z\" fill=\"#99F255\"/><path d=\"M 90.450850,20.610737 A 50.000000,50.000000 0 0 1 100.000000,50.000000 L 50.000000,50.000000 L 90.450850,20.610737 Z\" fill=\"#44CC11\"/><path d=\"M50,50 L25,50 A25,25 0 1,1 75,50 Z\" fill=\"white\"/><polygon points=\"54.852857095672256,45.61110924292003 45.14714290432774,45.61110924292004 49.999999999999986,5\" fill=\"black\" /><path d=\"M50,50 L40,50 A10,10 0 1,1 60,50 Z\" fill=\"black\"/><path d=\"M50,50 L45,50 A5,5 0 1,1 55.00000000000001,50 Z\" fill=\"#7A7A7A\"/></svg>"},
		{"", "75", "<svg height=\"50px\" width=\"100px\" viewBox=\"0 0 100 50\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M 0.000000,50.000000 A 50.000000,50.000000 0 0 1 9.549150,20.610737 L 50.000000,50.000000 L 0.000000,50.000000 Z\" fill=\"red\"/><path d=\"M 9.549150,20.610737 A 50.000000,50.000000 0 0 1 34.549150,2.447174 L 50.000000,50.000000 L 9.549150,20.610737 Z\" fill=\"orange\"/><path d=\"M 34.549150,2.447174 A 50.000000,50.000000 0 0 1 65.450850,2.447174 L 50.000000,50.000000 L 34.549150,2.447174 Z\" fill=\"yellow\"/><path d=\"M 65.450850,2.447174 A 50.000000,50.000000 0 0 1 90.450850,20.610737 L 50.000000,50.000000 L 65.450850,2.447174 Z\" fill=\"#99F255\"/><path d=\"M 90.450850,20.610737 A 50.000000,50.000000 0 0 1 100.000000,50.000000 L 50.000000,50.000000 L 90.450850,20.610737 Z\" fill=\"#44CC11\"/><path d=\"M50,50 L25,50 A25,25 0 1,1 75,50 Z\" fill=\"white\"/><polygon points=\"56.53490257669732,50.328073744260905 49.6719262557391,43.46509742330269 81.81980515339464,18.180194846605357\" fill=\"black\" /><path d=\"M50,50 L40,50 A10,10 0 1,1 60,50 Z\" fill=\"black\"/><path d=\"M50,50 L45,50 A5,5 0 1,1 55.00000000000001,50 Z\" fill=\"#7A7A7A\"/></svg>"},                        // Defaults to width=100
		{"200", "", "<svg height=\"100px\" width=\"200px\" viewBox=\"0 0 200 100\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M 0.000000,100.000000 A 100.000000,100.000000 0 0 1 19.098301,41.221475 L 100.000000,100.000000 L 0.000000,100.000000 Z\" fill=\"red\"/><path d=\"M 19.098301,41.221475 A 100.000000,100.000000 0 0 1 69.098301,4.894348 L 100.000000,100.000000 L 19.098301,41.221475 Z\" fill=\"orange\"/><path d=\"M 69.098301,4.894348 A 100.000000,100.000000 0 0 1 130.901699,4.894348 L 100.000000,100.000000 L 69.098301,4.894348 Z\" fill=\"yellow\"/><path d=\"M 130.901699,4.894348 A 100.000000,100.000000 0 0 1 180.901699,41.221475 L 100.000000,100.000000 L 130.901699,4.894348 Z\" fill=\"#99F255\"/><path d=\"M 180.901699,41.221475 A 100.000000,100.000000 0 0 1 200.000000,100.000000 L 100.000000,100.000000 L 180.901699,41.221475 Z\" fill=\"#44CC11\"/><path d=\"M100,100 L50,100 A50,50 0 1,1 150,100 Z\" fill=\"white\"/><polygon points=\"91.22221848584006,90.29428580865546 91.22221848584006,109.70571419134453 10,100\" fill=\"black\" /><path d=\"M100,100 L80,100 A20,20 0 1,1 120,100 Z\" fill=\"black\"/><path d=\"M100,100 L90,100 A10,10 0 1,1 110.00000000000001,100 Z\" fill=\"#7A7A7A\"/></svg>"}, // Defaults to percentage=0
		{"-10", "30", "<svg height=\"50px\" width=\"100px\" viewBox=\"0 0 100 50\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M 0.000000,50.000000 A 50.000000,50.000000 0 0 1 9.549150,20.610737 L 50.000000,50.000000 L 0.000000,50.000000 Z\" fill=\"red\"/><path d=\"M 9.549150,20.610737 A 50.000000,50.000000 0 0 1 34.549150,2.447174 L 50.000000,50.000000 L 9.549150,20.610737 Z\" fill=\"orange\"/><path d=\"M 34.549150,2.447174 A 50.000000,50.000000 0 0 1 65.450850,2.447174 L 50.000000,50.000000 L 34.549150,2.447174 Z\" fill=\"yellow\"/><path d=\"M 65.450850,2.447174 A 50.000000,50.000000 0 0 1 90.450850,20.610737 L 50.000000,50.000000 L 65.450850,2.447174 Z\" fill=\"#99F255\"/><path d=\"M 90.450850,20.610737 A 50.000000,50.000000 0 0 1 100.000000,50.000000 L 50.000000,50.000000 L 90.450850,20.610737 Z\" fill=\"#44CC11\"/><path d=\"M50,50 L25,50 A25,25 0 1,1 75,50 Z\" fill=\"white\"/><polygon points=\"51.346318600737554,43.59687495874813 43.49423087739373,49.30175062338621 23.5496636468387,13.594235253127373\" fill=\"black\" /><path d=\"M50,50 L40,50 A10,10 0 1,1 60,50 Z\" fill=\"black\"/><path d=\"M50,50 L45,50 A5,5 0 1,1 55.00000000000001,50 Z\" fill=\"#7A7A7A\"/></svg>"},                     // Invalid width, defaults to 100
		{"100", "101", "<svg height=\"50px\" width=\"100px\" viewBox=\"0 0 100 50\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M 0.000000,50.000000 A 50.000000,50.000000 0 0 1 9.549150,20.610737 L 50.000000,50.000000 L 0.000000,50.000000 Z\" fill=\"red\"/><path d=\"M 9.549150,20.610737 A 50.000000,50.000000 0 0 1 34.549150,2.447174 L 50.000000,50.000000 L 9.549150,20.610737 Z\" fill=\"orange\"/><path d=\"M 34.549150,2.447174 A 50.000000,50.000000 0 0 1 65.450850,2.447174 L 50.000000,50.000000 L 34.549150,2.447174 Z\" fill=\"yellow\"/><path d=\"M 65.450850,2.447174 A 50.000000,50.000000 0 0 1 90.450850,20.610737 L 50.000000,50.000000 L 65.450850,2.447174 Z\" fill=\"#99F255\"/><path d=\"M 90.450850,20.610737 A 50.000000,50.000000 0 0 1 100.000000,50.000000 L 50.000000,50.000000 L 90.450850,20.610737 Z\" fill=\"#44CC11\"/><path d=\"M50,50 L25,50 A25,25 0 1,1 75,50 Z\" fill=\"white\"/><polygon points=\"54.38889075707997,54.852857095672256 54.38889075707996,45.14714290432774 95,49.999999999999986\" fill=\"black\" /><path d=\"M50,50 L40,50 A10,10 0 1,1 60,50 Z\" fill=\"black\"/><path d=\"M50,50 L45,50 A5,5 0 1,1 55.00000000000001,50 Z\" fill=\"#7A7A7A\"/></svg>"},                                  // Clamp percentage to 100
	}

	preprocessSVG := func(svg string) string {
		var lines []string
		for _, line := range strings.Split(svg, "\n") {
			trimmedLine := strings.TrimSpace(line)
			if trimmedLine != "" {
				lines = append(lines, trimmedLine)
			}
		}
		return strings.Join(lines, "")
	}

	for _, test := range tests {
		router := gin.Default()
		router.GET("/test", HandleProgressGauge)
		// Create a test request with the specified parameters
		req, _ := http.NewRequest("GET", "/test?width="+test.widthParam+"&percentage="+test.percentageParam, nil)
		resp := httptest.NewRecorder()

		// Perform the test request
		router.ServeHTTP(resp, req)

		// Check if the response content-type is correct
		if contentType := resp.Header().Get("Content-Type"); contentType != "image/svg+xml" {
			t.Errorf("Expected Content-Type 'image/svg+xml', got '%s'", contentType)
		}

		actualSVG := preprocessSVG(resp.Body.String())
		expectedSVG := preprocessSVG(test.expectedSVG)
		if actualSVG != expectedSVG {
			fmt.Printf("Test case: WidthParam=%s, PercentageParam=%s\n", test.widthParam, test.percentageParam)
			fmt.Printf("Response Body: %s\n", actualSVG)
			t.Errorf("SVG output does not match expected output for test case with widthParam=%s, percentageParam=%s", test.widthParam, test.percentageParam)
		}

	}

	// Additional checks can be added to verify specific aspects of the SVG output
}

func TestCreatePiePath(t *testing.T) {
	tests := []struct {
		center     float64
		radius     float64
		startAngle float64
		endAngle   float64
		expected   string
	}{
		{100, 100, 180, 216, "M 0.000000,100.000000 A 100.000000,100.000000 0 0 1 19.098301,41.221475 L 100.000000,100.000000 L 0.000000,100.000000 Z"},
		{100, 100, 216, 252, "M 19.098301,41.221475 A 100.000000,100.000000 0 0 1 69.098301,4.894348 L 100.000000,100.000000 L 19.098301,41.221475 Z"},
		{100, 100, 0, 360, "M 200.000000,100.000000 A 100.000000,100.000000 0 0 1 200.000000,100.000000 L 100.000000,100.000000 L 200.000000,100.000000 Z"},
		{100, 100, 0, 1, "M 200.000000,100.000000 A 100.000000,100.000000 0 0 1 199.984770,101.745241 L 100.000000,100.000000 L 200.000000,100.000000 Z"},
		{100, 100, 0, 270, "M 200.000000,100.000000 A 100.000000,100.000000 0 0 1 100.000000,0.000000 L 100.000000,100.000000 L 200.000000,100.000000 Z"},
		{100, 0, 0, 90, "M 100.000000,100.000000 A 0.000000,0.000000 0 0 1 100.000000,100.000000 L 100.000000,100.000000 L 100.000000,100.000000 Z"},
		{100, 100, -90, -45, "M 100.000000,0.000000 A 100.000000,100.000000 0 0 1 170.710678,29.289322 L 100.000000,100.000000 L 100.000000,0.000000 Z"},
	}

	for _, test := range tests {
		actual := createPiePath(test.center, test.radius, test.startAngle, test.endAngle)
		if actual != test.expected {
			t.Errorf("createPiePath(%f, %f, %f, %f) = %s; want %s", test.center, test.radius, test.startAngle, test.endAngle, actual, test.expected)
		}
	}
}

func TestCalculateNeedlePosition(t *testing.T) {
	tests := []struct {
		center     float64
		percentage int
		expectedX1 float64
		expectedY1 float64
		expectedX2 float64
		expectedY2 float64
		expectedX3 float64
		expectedY3 float64
	}{
		{100, 0, 91.222218, 90.294286, 91.222218, 109.705714, 10.000000, 100.000000},      // 0%
		{100, 50, 109.705714, 91.222218, 90.294286, 91.222218, 100.000000, 10.000000},     // 50%
		{100, 100, 108.777782, 109.705714, 108.777782, 90.294286, 190.000000, 100.000000}, // 100%
	}

	for _, test := range tests {
		needle := calculateNeedlePosition(test.center, test.percentage)

		if floatEquals(needle.X1, test.expectedX1) && floatEquals(needle.Y1, test.expectedY1) &&
			floatEquals(needle.X2, test.expectedX2) && floatEquals(needle.Y2, test.expectedY2) &&
			floatEquals(needle.X3, test.expectedX3) && floatEquals(needle.Y3, test.expectedY3) {
			continue
		}
		t.Errorf("calculateNeedlePosition(%f, %d) produced incorrect triangle points\n"+
			"Expected: X1=%f, Y1=%f, X2=%f, Y2=%f, X3=%f, Y3=%f\n"+
			"Actual: X1=%f, Y1=%f, X2=%f, Y2=%f, X3=%f, Y3=%f",
			test.center, test.percentage,
			test.expectedX1, test.expectedY1, test.expectedX2, test.expectedY2, test.expectedX3, test.expectedY3,
			needle.X1, needle.Y1, needle.X2, needle.Y2, needle.X3, needle.Y3)
	}
}

func floatEquals(a, b float64) bool {
	const epsilon = 1e-5
	diff := math.Abs(a - b)
	if diff < epsilon {
		return true
	}
	fmt.Printf("Difference: %g, Epsilon: %g\n", diff, epsilon)
	return false
}
