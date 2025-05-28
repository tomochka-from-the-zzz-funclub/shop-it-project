package kafka

import (
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/infrastructure/config/kafka/config"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/infrastructure/config/kafka/connection"
	"log"
)

type Connections struct {
	ConsumerConnection *connection.Consumer
}

func NewKafkaConnections() *Connections {
	log.Println("[kafka] loading Kafka configuration and initializing connections")

	cfg := config.LoadKafkaConfig()

	consumerConn := connection.NewConsumerConnection(cfg)

	log.Println("[kafka] Consumer connection initialized")

	log.Println("[kafka] all Kafka connections are ready")

	return &Connections{
		ConsumerConnection: consumerConn,
	}
}
