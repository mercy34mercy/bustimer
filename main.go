package main

import (
	"os"

	"github.com/labstack/echo/v4"
)

var url = "https://ohmitetudo-bus.jorudan.biz/busstatedtl"

const (
	ritsumeikan   = "立命館大学〔近江鉄道・湖国バス〕"
	minamikusatsu = "立命館大学〔近江鉄道・湖国バス〕:2"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}
	e := echo.New()
	e.Debug = false
	setupMiddleware(e)
	routing(e)
	e.Logger.Fatal(e.Start(":" + port))
}
