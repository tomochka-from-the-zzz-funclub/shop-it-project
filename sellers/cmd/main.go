package main

import (
	"fmt"
	config "market/internal/cfg"
	"market/internal/services"
	"market/internal/transport"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Printf("%v", cfg)
	s := services.NewSrv(cfg)
	transport.HandleCreate(cfg, s)
}
