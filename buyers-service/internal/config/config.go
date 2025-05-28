package config

import (
	"sync"
	"time"

	"github.com/caarlos0/env"
	"github.com/gofiber/fiber/v2/log"
)

var once sync.Once

type Options struct {
	PostgresConfig PostgresCFG
	ServerConfig   ServerCFG
}

type ServerCFG struct {
	ServiceName     string        `env:"SERVICE_NAME" envDefault:"checker"`
	LogLevel        string        `env:"LOG_LEVEL" envDefault:"debug"`
	IsPretty        bool          `env:"IS_PRETTY" envDefault:"true"`
	SrvAddr         string        `env:"SRV_ADDR" envDefault:":8080"`
	BindMetrics     string        `env:"BIND_METRICS" envDefault:":9090"`
	GracefulTimeout time.Duration `env:"GRACEFUL_TIMEOUT" envDefault:"15s"`
}

type PostgresCFG struct {
	Host      string `env:"POSTGRES_HOST" envDefault:"postgres"`
	Port      int    `env:"POSTGRES_PORT" envDefault:"5432"`
	User      string `env:"POSTGRES_USER" envDefault:"postgres"`
	Password  string `env:"POSTGRES_PASSWORD" envDefault:"secret"`
	Db        string `env:"POSTGRES_DB" envDefault:"buyer_info"`
	JWTSecret string `env:"POSTGRES_JWT_SECRET" envDefault:"my-super-secret-key-1234567890"`
}

func New() *Options {
	var err error
	var s ServerCFG
	var p PostgresCFG
	once.Do(func() {
		err = env.Parse(&s)
		if err != nil {
			return
		}
		err = env.Parse(&p)
		if err != nil {
			return
		}
	})
	if err != nil {
		log.Errorf("error during parsing config: %v", err)
		return nil
	}
	return &Options{
		PostgresConfig: p,
		ServerConfig:   s,
	}
}
