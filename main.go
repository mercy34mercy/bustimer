package main

import (
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
		port = ":8000"
	}
	e := echo.New()
	e.Debug = true
	e.Use(middleware.Logger())
	e.GET("/bus/timetable", timetable.ScrapeTimeTable)
	// http.HandleFunc("/bus/time", getTimeTable)
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("/bus/timeにGETリクエストを送ってください"))
	// })
	// http.HandleFunc("/bus/timetable", scrapeTimeTable)
	// http.ListenAndServe(":"+port, nil)
	e.Logger.Fatal(e.Start(port))
}
