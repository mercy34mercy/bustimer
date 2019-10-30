package main

import (
	"os"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	approach "github.com/shun-shun123/bus-timer/approach"
	timetable "github.com/shun-shun123/bus-timer/timetable"
)

var url = "https://ohmitetudo-bus.jorudan.biz/busstatedtl"

const (
	ritsumeikan   = "立命館大学〔近江鉄道・湖国バス〕"
	minamikusatsu = "立命館大学〔近江鉄道・湖国バス〕:2"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":1323"
	}
	e := echo.New()
	e.Debug = false
	e.Use(middleware.Logger())
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h1>バスです。API</h1>")
	})
	e.GET("/bus/timetable", timetable.ScrapeTimeTable)
	e.GET("/bus/time", approach.ScrapeApproachInfo)
	e.Logger.Fatal(e.Start(":1323"))
}
