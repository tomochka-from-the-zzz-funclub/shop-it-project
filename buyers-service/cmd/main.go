package main

import (
	"buyers-service/internal/config"
	"buyers-service/internal/logger"
	"buyers-service/internal/services/auth"
	buyersservice "buyers-service/internal/services/buyers"
	"buyers-service/internal/storages/postgres"
	"buyers-service/internal/transport"
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	cfg := config.New()
	if cfg == nil {
		fmt.Fprintln(os.Stderr, "error during parsing config")
		os.Exit(1)
	}

	log := logger.New(cfg.ServerConfig.ServiceName, cfg.ServerConfig.IsPretty, cfg.ServerConfig.LogLevel)

	db := postgres.NewPostgres(cfg.PostgresConfig, log)

	authService := auth.NewAuthService(&db, cfg.PostgresConfig.JWTSecret, 24*time.Hour)

	buyersService := buyersservice.NewBuyersService(&db, authService, log)

	transportMiddlewares := []transport.Option{
		transport.SetFiberCfg(fiber.Config{Immutable: true}),
		transport.Use(cors.New(cors.Config{
			AllowOrigins:  "*",
			AllowMethods:  "GET,POST,DELETE",
			AllowHeaders:  "Authorization,Content-Type,X-Request-Id",
			ExposeHeaders: "X-Request-Id",
		})),
		transport.WithRequestID("X-Request-Id"),
		transport.Use(func(c *fiber.Ctx) error {
			if c.Path() == "/api/v1/auth/register" || c.Path() == "/api/v1/auth/login" {
				return c.Next()
			}
			token := c.Get("Authorization")
			if token == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
			}
			if len(token) > 7 && token[:7] == "Bearer " {
				token = token[7:]
			}
			userID, err := authService.ValidateToken(token)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
			}
			c.Locals("userID", userID)
			return c.Next()
		}),
	}

	servicesOpts := []transport.Option{
		transport.BuyersService(transport.NewBuyersService(&buyersService)),
	}

	srv := transport.New(log.Logger, append(transportMiddlewares, servicesOpts...)...).WithMetrics().WithLog()

	srv.ServeMetrics(log.Logger, "/", cfg.ServerConfig.BindMetrics)

	srv.Fiber().Get("/readiness", func(ctx *fiber.Ctx) error {
		ctx.Status(fiber.StatusOK)
		return ctx.JSON("OK")
	})
	srv.Fiber().Get("/liveness", func(ctx *fiber.Ctx) error {
		ctx.Status(fiber.StatusOK)
		return ctx.JSON("OK")
	})

	go func() {
		if err := srv.Fiber().Listen(cfg.ServerConfig.SrvAddr); err != nil {
			log.LogErrorf("Failed to start server on %s: %v", cfg.ServerConfig.SrvAddr, err)
		}
	}()

	fmt.Printf("service %s started on %s\n", cfg.ServerConfig.ServiceName, cfg.ServerConfig.SrvAddr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ServerConfig.GracefulTimeout)
	defer cancel()

	if err := srv.Fiber().ShutdownWithContext(ctx); err != nil {
		log.LogErrorf("Failed to shutdown server: %v", err)
	}
	db.Close()
	log.LogInfo("service shut down")
}
