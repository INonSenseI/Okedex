package main

import (
	"os"
)

func main() {
	exitCode := app.Run()
	os.Exit(exitCode)
}
