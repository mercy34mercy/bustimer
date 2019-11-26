package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	approach "github.com/shun-shun123/bus-timer/approach"
	timetable "github.com/shun-shun123/bus-timer/timetable"
)

// routing APIのルーティングを行う
func routing(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h1>busdes!API</h1>")
	})
	e.GET("/bus/timetable", timetable.ScrapeTimeTable)
	e.GET("/bus/time", approach.ScrapeApproachInfo)
}
