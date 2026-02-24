package main

import (
	"log"

	"github.com/ilyas/flower/services/auth/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
