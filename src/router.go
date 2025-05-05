package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shun-shun123/bus-timer/src/infrastructure"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

var e = echo.New()

func Routing() {
	// OpenTelemetryミドルウェアを追加
	e.Use(otelecho.Middleware("bustimer-service"))

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h1>Busdes! Clean Architecture API</h1>")
	})
	e.GET("/bus/time/v3", func(c echo.Context) error {
		cc := &CustomContext{c}
		return infrastructure.ApproachInfoRequest(cc)
	})
	e.GET("/bus/timetable", func(c echo.Context) error {
		cc := &CustomContext{c}
		return infrastructure.TimeTableRequest(cc)
	})
	e.GET("/system/info", func(c echo.Context) error {
		cc := &CustomContext{c}
		return infrastructure.SystemInfoRequest(cc)
	})
	e.GET("/debug/system/info/error", func(c echo.Context) error {
		cc := &CustomContext{c}
		return infrastructure.DebugErrorSystemInfoRequest(cc)
	})
	e.GET("/debug/system/info/success", func(c echo.Context) error {
		cc := &CustomContext{c}
		return infrastructure.DebugSuccessSystemInfoRequest(cc)
	})
}
