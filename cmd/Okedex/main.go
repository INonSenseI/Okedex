package main

import (
	app "Okedex/internal/app"
	"log"
	"os"
)

func main() {
	f, err := os.OpenFile("logs/log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	exitCode := app.Run()
	os.Exit(exitCode)
}
