# Dynamic README Elements

This project provides a Go API using the Gin framework for generating customizable SVG progress bars, with options for both linear and circular styles.

## Features

- **Linear Progress Bar**: Generates a horizontal bar to visually represent progress. Customizable in size and fill percentage.
- **Circular Progress Bar**: Creates a circular or "donut" style progress indicator. Size and progress fill are adjustable.
- **Waffle Progress Chart**: Displays progress in a grid or 'waffle' format. Offers customization in grid size, square count, and filled percentage.
- **Calendar Progress Chart**: Shows a monthly calendar view with specific days marked to indicate progress. Customizable by year, month, and progress days.

## Getting Started

### Prerequisites

- Go (version 1.15 or later)
- Gin web framework

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/2ajoyce/dynamic-readme-elements.git
   ```
2. Navigate to the project directory:
   ```bash
   cd dynamic-readme-elements
   ```
3. Run the server:
   ```bash
   go run main.go
   ```

## Usage

Generate SVG progress bars by accessing the endpoints with specific query parameters:

### Linear Progress Bar

- **Endpoint**: `/progress/bar`
- **Parameters**: `width`, `height`, `percentage`
- **Example**: `http://localhost:8080/progress/bar?width=100&height=25&percentage=72`

![Linear Progress Bar](https://progress.2ajoyce.com/progress/bar?width=100&height=25&percentage=72)

### Circular Progress Bar

- **Endpoint**: `/progress/circle`
- **Parameters**: `size`, `percentage`
- **Example**: `http://localhost:8080/progress/circle?size=100&percentage=72`

![Circular Progress Bar](https://progress.2ajoyce.com/progress/circle?size=100&percentage=72)

### Waffle Progress Chart

- **Endpoint**: `/progress/waffle`
- **Parameters**: `width`, `numberOfSquares`,`percentage`
- **Example**: `http://localhost:8080/progress/waffle?width=100&numberOfSquares=100&percentage=99`

![Waffle Progress Chart](https://progress.2ajoyce.com/progress/waffle?width=100&numberOfSquares=100&percentage=99)

### Calendar Progress Chart

- **Endpoint**: `/calendar`
- **Parameters**: `year`, `month`, `progressDays` (optional; comma-separated list of days)
- **Default**: Defaults to the current year and month if not provided.
- **Example**: `http://localhost:8080/calendar?year=2023&month=1&progressDays=2,15,20`

![Calendar Progress Chart](https://progress.2ajoyce.com/calendar)

## Customization

Each progress indicator type offers specific customization options through query parameters:

- **Linear Progress Bar**: Adjust the `width`, `height`, and `percentage` to control the bar's dimensions and progress.
- **Circular Progress Bar**: Modify the `size` for the diameter and `percentage` for progress representation.
- **Waffle Progress Chart**: Change the `width` to control the overall size, `numberOfSquares` for grid density, and `percentage` for filled squares.
- **Calendar Progress Chart**: Set `year`, `month`, and optionally `progressDays` to display progress on specific days of a month.

## Acknowledgments

- Inspired by [Frederico Jordan's progress-bar repository](https://github.com/fredericojordan/progress-bar)
- Thanks to the Go and Gin communities for their resources and support.
