package main

import (
	"log"

	"github.com/ilyas/flower/services/catalog/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
