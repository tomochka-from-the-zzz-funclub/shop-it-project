package main

import (
	"fmt"
	"search/internal/config"
	"search/internal/service"
	"search/internal/transport"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Printf("%v", cfg)
	s := service.NewSearchService(cfg)
	transport.HandleCreate(cfg, s)
}
