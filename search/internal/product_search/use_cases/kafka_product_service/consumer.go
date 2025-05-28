package kafka_product_service

import (
	"context"
	"log"
	"sync"
)

func (kps *Service) StartConsumerWorkers(ctx context.Context, numWorkers int) {

	log.Printf("[kafka_service] starting %d consumer workers", numWorkers)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {

		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			log.Printf("[kafka_service][consumer %d] worker started", workerID)

			kps.consumeLoop(ctx, workerID)

			log.Printf("[kafka_service][consumer %d] worker stopped", workerID)
		}(i)
	}

	wg.Wait()

	log.Println("[kafka_service] all consumer workers exited")
}

func (kps *Service) consumeLoop(ctx context.Context, workerID int) {
	defer func() {
		if r := recover(); r != nil {

			log.Printf("[kafka_service][consumer %d][PANIC] recovered: %v", workerID, r)
		}
	}()

	for {
		select {
		case <-ctx.Done():

			log.Printf("[kafka_service][consumer %d] context cancelled, exiting", workerID)

			return

		default:

			log.Printf("[kafka_service][consumer %d] polling for message", workerID)

			if err := kps.CreateProductFromKafka(ctx); err != nil {

				log.Printf("[kafka_service][consumer %d][ERROR] processing message: %v", workerID, err)

			} else {

				log.Printf("[kafka_service][consumer %d] message processed successfully", workerID)
			}
		}
	}
}

func (kps *Service) CreateProductFromKafka(ctx context.Context) error {
	log.Println("[kafka_product_service] CreateProductFromKafka invoked")

	product, err := kps.Consumer.ReadMessage(ctx)
	if err != nil {
		log.Printf("[kafka_product_service][consumer] Error reading message from Kafka: %v", err)
		return err
	}

	log.Printf("[kafka_product_service][consumer] Received Product: %+v", product)

	if err = kps.ProductService.Repo.Save(product); err != nil {
		log.Printf("[kafka_product_service][consumer][ERROR] Error saving product: %v", err)
		return err
	}

	log.Printf("[kafka_product_service][consumer] Product saved successfully: ID=%s", product.ID)
	return nil
}
