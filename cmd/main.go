package main

import (
	"log"

	"github.com/FacelessWayfarer/test-task-medods/internal/app"
)

// @title           Test task API
// @version         1.0
// @description     This is a test task
// @host      localhost:8080
// @BasePath  /
func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
