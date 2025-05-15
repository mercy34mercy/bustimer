package main

import (
	"io/fs"
	"net/http"
	"path"

	"github.com/labstack/echo/v4"
	"github.com/shun-shun123/bus-timer/src/infrastructure"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

var e = echo.New()

func Routing() {
	// OpenTelemetryミドルウェアを追加
	e.Use(otelecho.Middleware("bustimer-service"))

	// 埋め込みファイルを使用した静的ファイル配信用のルートを追加
	embeddedFS := GetEmbeddedFile()

	// 埋め込まれた静的ファイル用のハンドラーを設定
	staticFs, err := fs.Sub(embeddedFS, "static")
	if err != nil {
		e.Logger.Fatal(err)
	}

	// 静的ファイル用のハンドラー
	fileServer := http.FileServer(http.FS(staticFs))

	// /embedded/ パスで埋め込まれたファイルにアクセスできるように設定
	e.GET("/static/*", echo.WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// URLパスから"/embedded/"を削除
		r.URL.Path = path.Join("/", r.URL.Path[len("/static/"):])
		fileServer.ServeHTTP(w, r)
	})))

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
