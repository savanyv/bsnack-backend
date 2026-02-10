package main

import (
	"log"
	"os"

	"github.com/savanyv/bsnack-backend/config"
	"github.com/savanyv/bsnack-backend/internal/app"
)

func main() {
	cfg := config.LoadConfig()

	server := app.NewServer(cfg)
	if err := server.Start(); err != nil {
		log.Printf("Server error: %v", err)
		os.Exit(1)
	}
}
