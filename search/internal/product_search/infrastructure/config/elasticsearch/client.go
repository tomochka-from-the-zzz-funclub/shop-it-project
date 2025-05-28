package elasticsearch

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"net/http"
	"time"
)

type Client struct {
	Connection *elasticsearch.Client
}

func NewClient() (*Client, error) {
	cfg := ReadConfigFromEnv()

	esCfg := elasticsearch.Config{
		Addresses: []string{"http://" + cfg.Address},
		Username:  cfg.Username,
		Password:  cfg.Password,
		Transport: &http.Transport{
			MaxIdleConns:        cfg.MaxIdleConns,
			MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
			IdleConnTimeout:     cfg.IdleConnTimeout,
		},
	}

	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("[elasticsearch][FATAL] failed to create ES client: %w", err)
	}

	const maxAttempts = 5
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		log.Printf("[elasticsearch] ping attempt %d/%d", attempt, maxAttempts)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		res, err := es.Info(es.Info.WithContext(ctx))
		if err != nil {
			lastErr = err
			log.Printf("[elasticsearch][WARN] ping error: %v", err)
		} else {
			defer res.Body.Close()
			if res.IsError() {
				body := res.String()
				lastErr = fmt.Errorf("status %s, body %s", res.Status(), body)
				log.Printf("[elasticsearch][WARN] ping returned error: %s", body)
			} else {
				log.Printf("[elasticsearch] connected to %s (status: %s)", cfg.Address, res.Status())
				return &Client{Connection: es}, nil
			}
		}

		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("[elasticsearch][FATAL] could not ping cluster after %d attempts: %w", maxAttempts, lastErr)
}
