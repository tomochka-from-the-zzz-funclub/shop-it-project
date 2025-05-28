package kafka_product_service

import (
	"context"
	repository "gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/infrastructure/repository/elasticsearch"
	"log"

	myKafka "gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/infrastructure/config/kafka"
	kafkaRepo "gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/infrastructure/repository/kafka/product"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/use_cases"
)

type Service struct {
	ProductService *use_cases.ProductService
	Consumer       *kafkaRepo.Consumer
}

func NewKafkaService(kafkaConn *myKafka.Connections, repo *repository.ProductRepository) *Service {
	log.Println("[kafka_product_service] Initializing Kafka service")
	svc := &Service{
		ProductService: use_cases.NewProductService(repo),
		Consumer:       kafkaRepo.NewConsumer(kafkaConn.ConsumerConnection),
	}

	log.Println("[kafka_product_service] Kafka service initialized")
	return svc
}

func (s *Service) Start(ctx context.Context, consumerWorkers, producerWorkers int) {
	log.Printf("[kafka_product_service] Starting %d consumer workers, %d producer workers", consumerWorkers, producerWorkers)
	go s.StartConsumerWorkers(ctx, consumerWorkers)
}
