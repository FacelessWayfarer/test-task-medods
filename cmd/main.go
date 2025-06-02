package main

import (
	"log"

	"github.com/FacelessWayfarer/test-task-medods/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
