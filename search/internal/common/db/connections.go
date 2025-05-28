package db

import (
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/infrastructure/config/elasticsearch"
	"log"
)

type Connections struct {
	ESClient *elasticsearch.Client
}

func NewConnections() (*Connections, error) {
	log.Println("[db] Initializing database connections")

	// Elasticsearch
	esClient, err := elasticsearch.NewClient()
	if err != nil {
		log.Printf("[db][ERROR] failed to initialize Elasticsearch client: %v", err)
		return nil, err
	}
	log.Println("[db] Elasticsearch client created")

	return &Connections{
		ESClient: esClient,
	}, nil
}
