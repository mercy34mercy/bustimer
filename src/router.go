package main

import (
	"github.com/labstack/echo/v4"
	"github.com/shun-shun123/bus-timer/src/infrastructure"
	"net/http"
)

var e = echo.New()

func Routing() {
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h1>Busdes! Clean Architecture API</h1>")
	})
	e.GET("/bus/time/v3", func(c echo.Context) error {
		cc := &CustomContext{c}
		return infrastructure.ApproachInfoRequest(cc)
	})
}


