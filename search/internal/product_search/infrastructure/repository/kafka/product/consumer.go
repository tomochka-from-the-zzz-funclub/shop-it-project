package product

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/domain/product/entity"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/infrastructure/config/kafka/connection"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/infrastructure/config/kafka/contracts"
	"log"
)

type Consumer struct {
	Reader *kafka.Reader
}

func NewConsumer(consumerConnection *connection.Consumer) *Consumer {
	cfg := consumerConnection.Config

	log.Printf("[kafka:product] initializing consumer: topic=%s, groupID=%s", cfg.UserOrderTopic, cfg.GroupID)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        consumerConnection.Brokers,
		Topic:          cfg.UserOrderTopic,
		GroupID:        cfg.GroupID,
		Dialer:         consumerConnection.Connection,
		MinBytes:       cfg.MinBytes,
		MaxBytes:       cfg.MaxBytes,
		StartOffset:    cfg.StartOffset,
		CommitInterval: cfg.CommitInterval,
	})

	log.Println("[kafka:product] consumer reader created")

	return &Consumer{Reader: reader}
}

func (c *Consumer) ReadMessage(ctx context.Context) (*entity.Product, error) {
	log.Println("[kafka:product] reading message from Kafka")

	msg, err := c.Reader.ReadMessage(ctx)
	if err != nil {
		log.Printf("[kafka:product][ERROR] reading message failed: %v", err)
		return nil, err
	}
	log.Printf(
		"[kafka:product] message received: topic=%s partition=%d offset=%d",
		msg.Topic, msg.Partition, msg.Offset,
	)

	var dto contracts.ProductKafkaDTO
	if err = json.Unmarshal(msg.Value, &dto); err != nil {
		log.Printf("[kafka:product][ERROR] unmarshalling message failed: %v", err)
		return nil, err
	}
	log.Printf("[kafka:product] DTO parsed: %+v", dto)

	product := &entity.Product{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		Stock:       dto.Stock,
		Category:    dto.Category,
		Brand:       dto.Brand,
	}

	log.Printf("[kafka:product] Product entity created: %+v", product)
	return product, nil
}

func (c *Consumer) Close() error {

	log.Println("[kafka:product] closing consumer reader")

	if err := c.Reader.Close(); err != nil {

		return fmt.Errorf("can't close reader: %w", err)

	}

	log.Println("[kafka:product] consumer reader closed")

	return nil
}
