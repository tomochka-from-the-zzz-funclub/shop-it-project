package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/infrastructure/config/elasticsearch"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/infrastructure/config/kafka"
	repository "gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/infrastructure/repository/elasticsearch"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/transport/rest"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/use_cases/kafka_product_service"
	"log"
	"os"
	"strconv"
)

func Run(ctx context.Context, engine *gin.Engine, es *elasticsearch.Client) {
	log.Printf("[search_run] starting product search service")

	repo := repository.NewProductRepository(es, "product")

	log.Printf("[search_run] product repository created")

	log.Printf("[search_run] starting Kafka service")

	kafkaConn := kafka.NewKafkaConnections()

	log.Printf("[search_run] Kafka connections created")

	kafkaProductService := kafka_product_service.NewKafkaService(kafkaConn, repo)

	log.Printf("[search_run] Kafka service created")

	log.Printf("[search_run] starting REST service")

	rest.RegisterRoutes(ctx, engine, es, repo)

	log.Printf("[search_run] REST service started")

	cw, pw, err := loadKafkaWorkers()
	if err != nil {
		log.Printf("[search_run][ERROR] loading Kafka workers: %v", err)
	} else {
		kafkaProductService.Start(ctx, cw, pw)
		log.Printf("[search_run] Kafka service started (consumers=%d, producers=%d)", cw, pw)
	}
}

func loadKafkaWorkers() (int, int, error) {
	log.Println("[search_run] reading Kafka worker settings from env")

	cwStr := os.Getenv("KAFKA_CONSUMER_WORKERS")

	cw, err := strconv.Atoi(cwStr)

	if err != nil {
		return 0, 0, fmt.Errorf("invalid KAFKA_CONSUMER_WORKERS: %w", err)
	}

	pwStr := os.Getenv("KAFKA_PRODUCER_WORKERS")

	pw, err := strconv.Atoi(pwStr)

	if err != nil {
		return 0, 0, fmt.Errorf("invalid KAFKA_PRODUCER_WORKERS: %w", err)
	}

	log.Printf("[search_run] Kafka workers parsed: consumers=%d, producers=%d", cw, pw)

	return cw, pw, nil
}
