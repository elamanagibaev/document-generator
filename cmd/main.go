package main

import (
	"log"
	"os"

	"document-generator/internal/app"
	"document-generator/internal/config"
)

func main() {
	cfg := config.MustLoad()

	application := app.New(cfg)

	if err := application.Run(); err != nil {
		log.Println("server error:", err)
		os.Exit(1)
	}
}
