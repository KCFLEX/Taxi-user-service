package main

import (
	"log"

	"github.com/KCFLEX/Taxi-user-service/internal/handlers"
)

func main() {
	app := handlers.New()

	err := app.Serve()

	if err != nil {
		log.Fatal(err)
	}
}
