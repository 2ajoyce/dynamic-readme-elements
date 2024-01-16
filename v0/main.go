package main

import (
	"github.com/2ajoyce/dynamic-readme-elements/v0/internal"
	"github.com/2ajoyce/dynamic-readme-elements/v0/internal/svggen"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	// Route for a calendar
	router.GET("/calendar", svggen.HandleCalendar)

	// Route for a rectangular loading bar
	router.GET("/progress/bar", svggen.HandleProgressBar)

	// Route for a circular progress bar
	router.GET("/progress/circle", svggen.HandleProgressCircle)

	// Route for a gauge progress chart
	router.GET("/progress/gauge", svggen.HandleProgressGauge)

	// Route for a waffle progress chart
	router.GET("/progress/waffle", svggen.HandleProgressWaffle)

	// Route for health check
	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Route for version endpoint
	router.GET("/version", internal.HandleVersion)

	err := router.Run(":8080")
	if err != nil {
		return
	}
}
