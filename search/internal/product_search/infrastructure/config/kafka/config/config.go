package config

import "log"

type KafkaConfig struct {
	Common   KafkaCommonConfig
	Consumer KafkaConsumerConfig
}

func LoadKafkaConfig() KafkaConfig {
	log.Println("[config:kafka] loading Kafka config from environment")

	cfg := KafkaConfig{
		Common: KafkaCommonConfig{
			Username:       getEnv("KAFKA_USERNAME"),
			Password:       getEnv("KAFKA_PASSWORD"),
			BootstrapHosts: parseHosts("KAFKA_BOOTSTRAP_SERVERS"),
		},
		Consumer: KafkaConsumerConfig{
			GroupID:        getEnv("KAFKA_CONSUMER_GROUP_ID"),
			UserOrderTopic: getEnv("KAFKA_CONSUMER_USER_ORDER_TOPIC"),
			StartOffset:    parseInt64("KAFKA_CONSUMER_START_OFFSET"),
			CommitInterval: parseDuration("KAFKA_CONSUMER_COMMIT_INTERVAL"),
			UseTLS:         parseBool("KAFKA_USE_TLS"),
			MinBytes:       parseInt("KAFKA_CONSUMER_MIN_BYTES"),
			MaxBytes:       parseInt("KAFKA_CONSUMER_MAX_BYTES"),
		},
	}

	log.Printf("[config:kafka] loaded config: common=%+v consumer=%+v",
		cfg.Common, cfg.Consumer)

	return cfg
}
