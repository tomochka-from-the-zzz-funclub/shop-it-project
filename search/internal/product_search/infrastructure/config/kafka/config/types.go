package config

import (
	"time"
)

type KafkaCommonConfig struct {
	Username       string
	Password       string
	BootstrapHosts []string
}

type KafkaConsumerConfig struct {
	GroupID        string
	UserOrderTopic string
	StartOffset    int64
	CommitInterval time.Duration
	UseTLS         bool
	MinBytes       int
	MaxBytes       int
}
