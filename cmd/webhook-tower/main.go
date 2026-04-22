package main

import (
	"flag"
	"log"
	"webhook-tower/internal/config"
	"webhook-tower/internal/router"
)

func main() {
	configFile := flag.String("config", "config.yaml", "path to configuration file")
	flag.Parse()

	cfg, err := config.LoadConfigFromFile(*configFile)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	r := router.NewRouter(cfg)

	log.Printf("Starting Webhook Tower on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
