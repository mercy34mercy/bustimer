package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// setupMiddleware Echoのミドルウェアの初期設定を行う
func setupMiddleware(e *echo.Echo) {
	e.Use(middleware.Logger())
}
