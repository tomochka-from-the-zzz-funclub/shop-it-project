package connection

import (
	"crypto/tls"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/infrastructure/config/kafka/config"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type Consumer struct {
	Connection *kafka.Dialer
	Brokers    []string
	Config     config.KafkaConsumerConfig
}

func NewConsumerConnection(kafkaCfg config.KafkaConfig) *Consumer {

	log.Println("[kafka:consumer] initializing consumer connection")

	log.Printf("[kafka:consumer] brokers=%v", kafkaCfg.Common.BootstrapHosts)

	log.Printf("[kafka:consumer] groupID=%s topic=%s startOffset=%d useTLS=%t",
		kafkaCfg.Consumer.GroupID,
		kafkaCfg.Consumer.UserOrderTopic,
		kafkaCfg.Consumer.StartOffset,
		kafkaCfg.Consumer.UseTLS,
	)

	mechanism := plain.Mechanism{
		Username: kafkaCfg.Common.Username,
		Password: kafkaCfg.Common.Password,
	}

	log.Printf("[kafka:consumer] SASL mechanism configured for product=%s", kafkaCfg.Common.Username)

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	if kafkaCfg.Consumer.UseTLS {

		dialer.TLS = &tls.Config{}

		log.Println("[kafka:consumer] TLS enabled")
	}

	log.Println("[kafka:consumer] consumer dialer created")

	return &Consumer{
		Connection: dialer,
		Brokers:    kafkaCfg.Common.BootstrapHosts,
		Config:     kafkaCfg.Consumer,
	}
}
