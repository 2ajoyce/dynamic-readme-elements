package main

import (
	"github.com/2ajoyce/dynamic-readme-elements/internal/svggen"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Route for a calendar
	router.GET("/calendar", svggen.HandleCalendar)

	// Route for a rectangular loading bar
	router.GET("/progress/bar", svggen.HandleProgressBar)

	// Route for a circular progress bar
	router.GET("/progress/circle", svggen.HandleProgressCircle)

	err := router.Run(":8080")
	if err != nil {
		return
	}
}
