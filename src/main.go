package main

import "os"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}
	Routing()
	e.Debug = true
	e.Logger.Fatal(e.Start(":" + port))
}