package main

import (
	"context"
	"log"
	"os"

	"github.com/prometheus-agi/prometheus/gateway/internal/router"
)

func main() {
	var repo router.Repository = router.NewSwarmStore()
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		pgRepo, err := router.NewPostgresRepository(context.Background(), databaseURL)
		if err != nil {
			log.Printf("failed to connect postgres, falling back to in-memory store: %v", err)
		} else {
			defer pgRepo.Close()
			repo = pgRepo
			log.Printf("postgres repository enabled")
		}
	}
	deps := router.NewDependencies(repo)
	r := router.SetupRouter(deps)
	if err := r.Run(":3001"); err != nil {
		log.Fatal(err)
	}
}
