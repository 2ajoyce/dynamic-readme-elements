package svggen

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"math"
	"strconv"
)

const waffleChartTemplateStr = `
<svg width="{{.Width}}px" height="{{.Height}}px" xmlns="http://www.w3.org/2000/svg">
	<rect x="0" y="0" width="{{.Width}}px" height="{{.Height}}px" fill="none" />
    {{range .Squares}}
        <rect class="gridSquare" x="{{.X}}px" y="{{.Y}}px" width="{{$.SquareSize}}px" height="{{$.SquareSize}}px" fill="{{.Color}}" />
    {{end}}
</svg>
`

var waffleChartTemplate = template.Must(template.New("waffleChart").Parse(waffleChartTemplateStr))

func CalculateGridSize(width, numberOfSquares, gap int) (int, int) {
	// Validate the inputs
	if width <= 0 {
		width = 100
	}
	if numberOfSquares <= 0 {
		numberOfSquares = 100
	}

	// Estimate the number of squares per row
	estimatedSquaresPerRow := int(math.Sqrt(float64(numberOfSquares)))

	// Calculate the total gap space for the estimated number of squares per row
	totalGapSpace := (estimatedSquaresPerRow - 1) * gap

	// Adjust the width to account for the gaps
	adjustedWidth := width - totalGapSpace

	// Calculate the ideal size of each square based on the adjusted width
	idealSquareSize := int(math.Max(10, float64(adjustedWidth)/float64(estimatedSquaresPerRow)))

	// Recalculate the number of squares per row based on the ideal square size
	squaresPerRow := adjustedWidth / idealSquareSize

	// Ensure there is at least one square per row
	if squaresPerRow == 0 {
		squaresPerRow = 1
	}

	// Calculate the number of squares per column
	squaresPerColumn := numberOfSquares / squaresPerRow

	// Adjust for any remaining squares
	if numberOfSquares%squaresPerRow != 0 {
		squaresPerColumn++
	}

	return squaresPerRow, squaresPerColumn
}

func GenerateSquares(width, squaresPerRow, numberOfSquares, filledSquares, gap int) []struct {
	X, Y  int
	Color string
} {
	squares := make([]struct {
		X, Y  int
		Color string
	}, numberOfSquares)

	// Calculate total available width for squares, excluding gaps
	totalWidthForSquares := width - (gap * (squaresPerRow - 1))

	// Calculate width of each square
	squareWidth := totalWidthForSquares / squaresPerRow

	// Calculate the extra width to be distributed among squares
	extraWidthPerSquare := (totalWidthForSquares % squaresPerRow) / squaresPerRow

	for i := 0; i < numberOfSquares; i++ {
		// Calculate the X and Y position, including gaps and extra width
		x := gap + (i%squaresPerRow)*(squareWidth+gap+extraWidthPerSquare)
		y := gap + (i/squaresPerRow)*(squareWidth+gap)

		color := Colors.ProgressInactive // Default color for unfilled squares
		if i < filledSquares {
			color = Colors.ProgressActive // Color for filled squares
		}

		squares[i] = struct {
			X, Y  int
			Color string
		}{X: x, Y: y, Color: color}
	}
	return squares
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

func HandleProgressWaffle(c *gin.Context) {
	width := parseOrDefault(c.DefaultQuery("width", "100"), 10) // Minimum width is 10
	numberOfSquares := parseOrDefault(c.DefaultQuery("numberOfSquares", "100"), 100)
	gap := 3 // Gap between squares

	squaresPerRow, squaresPerColumn := CalculateGridSize(width, numberOfSquares, gap)

	// Calculate the adjusted width to account for the gaps
	totalGapSpace := (squaresPerRow - 1) * gap
	adjustedWidth := width - totalGapSpace

	// Calculate the size of each square based on the adjusted width
	squareSize := adjustedWidth / squaresPerRow

	height := (squareSize+gap)*squaresPerColumn - gap // Adjust the height to include gaps between rows

	percentage := clamp(parseOrDefault(c.DefaultQuery("percentage", "0"), 0), 0, 100)
	filledSquares := numberOfSquares * percentage / 100

	squares := GenerateSquares(width, squaresPerRow, numberOfSquares, filledSquares, gap)

	data := struct {
		Width, Height, SquareSize int
		Squares                   []struct {
			X, Y  int
			Color string
		}
	}{
		Width:      ((squareSize + gap) * squaresPerRow) + gap,
		Height:     height + 2*gap,
		SquareSize: squareSize,
		Squares:    squares,
	}

	c.Writer.Header().Set("Content-Type", "image/svg+xml")
	if err := waffleChartTemplate.Execute(c.Writer, data); err != nil {
		log.Printf("Error executing waffle chart template: %v\n", err)
	}
}

func parseOrDefault(value string, defaultVal int) int {
	if parsed, err := strconv.Atoi(value); err == nil {
		return parsed
	}
	return defaultVal
}
