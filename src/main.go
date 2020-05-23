package main

import (
	"github.com/shun-shun123/bus-timer/src/infrastructure"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}
	go infrastructure.TimeTableCacheStart()
	Routing()
	e.Debug = true
	e.Logger.Fatal(e.Start(":" + port))
}