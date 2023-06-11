package main

import (
	"log"

	"github.com/Pedrommb91/go-auth/config"
	"github.com/Pedrommb91/go-auth/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	app.Run(cfg)
}
