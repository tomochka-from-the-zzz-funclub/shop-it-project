package main

import (
	"fmt"
	config "goods/internal/cfg"
	"goods/internal/services"
	"goods/internal/transport"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Printf("%v", cfg)
	s := services.NewSrv(cfg)
	transport.HandleCreate(cfg, s)
}
