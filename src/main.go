package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/shun-shun123/bus-timer/src/infrastructure"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}

	// Cloud Traceの初期化
	log.Println("Starting application with Cloud Trace...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("Calling InitTracer...")
	tp, err := infrastructure.InitTracer()
	if err != nil {
		log.Printf("Failed to initialize tracer: %v", err)
		// トレーサーの初期化に失敗してもアプリケーションは続行
	} else {
		log.Println("Cloud Trace initialized successfully, deferring shutdown")
		// グレースフルシャットダウン時にトレーサーも終了
		defer infrastructure.ShutdownTracer(ctx, tp)
	}

	log.Println("Starting TimeTableCache...")
	go infrastructure.TimeTableCacheStart()
	log.Println("Setting up routes...")
	Routing()
	e.Debug = true
	log.Println("Debug mode enabled")

	// シグナル処理を追加してグレースフルシャットダウン
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
		<-sigCh

		log.Println("Shutting down server...")
		if err := e.Shutdown(ctx); err != nil {
			log.Printf("Error during server shutdown: %v", err)
		}
		cancel()
	}()

	log.Printf("Starting server on port %s...", port)
	if err := e.Start(":" + port); err != nil {
		log.Printf("Server stopped: %v", err)
	}
}
