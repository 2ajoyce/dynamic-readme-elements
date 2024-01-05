# Dynamic README Elements

This project provides a Go API using the Gin framework for generating customizable SVG progress bars, with options for both linear and circular styles.

## Features

- **Linear Progress Bar**: Generates a horizontal bar indicating progress.
- **Circular Progress Bar**: Creates a circular, "donut" style progress indicator.
- **Customizable**: Adjust size, percentage, and other properties via query parameters.

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

**Linear Progress Bar**

- Endpoint: /bar
- Parameters: width, height, percentage
- Example: http://localhost:8080/bar?width=300&height=50&percentage=75

**Circular Progress Bar**

- Endpoint: /circle
- Parameters: size, percentage
- Example: http://localhost:8080/circle?size=120&percentage=75

## Customization

Modify query parameters for customization:

- width and height for the linear bar (in pixels).
- size for the diameter of the circular bar (in pixels).
- percentage for progress representation (0 to 100).

## Acknowledgments

- Inspired by [Frederico Jordan's progress-bar repository](https://github.com/fredericojordan/progress-bar)
- Thanks to the Go and Gin communities for their resources and support.
