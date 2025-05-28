package main

import (
	"context"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/app"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/common/logger"
	"log"
)

func main() {
	ctx := context.Background()
	logger.Init()

	log.Println("Starting user-service")

	if err := app.Run(ctx); err != nil {
		log.Fatalf("Application terminated with error: %v", err)
	}
}
