// Package app
// @title Product Search API
// @version 1.0
// @description API для поиска и управления каталогом продуктов
// @host localhost:8080
// @BasePath /
// @schemes http
package app

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gitlab.mai.ru/4-bogatyra/backend/search/docs"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/common/db"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/common/settings"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/app"
	"time"

	"log"
)

func Run(ctx context.Context) error {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(loggingMiddleware)

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	docs.SwaggerInfo.BasePath = ""
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	if err := settings.LoadEnv(".env"); err != nil {

		log.Printf("failed to load env: %v", err)

		return err
	}

	log.Println("Environment variables loaded")

	connections, err := db.NewConnections()
	if err != nil {
		log.Printf("[ERROR] failed to establish connections: %v", err)
		return err
	}
	log.Println("Database connections established")

	app.Run(ctx, engine, connections.ESClient)

	addr := ":8080"
	swaggerURL := fmt.Sprintf("http://localhost%s/swagger/index.html", addr)
	log.Printf("Swagger UI available at %s", swaggerURL)

	if err := engine.Run(addr); err != nil {

		log.Printf("Error starting the application: %v", err)

		return err
	}

	return nil
}

func loggingMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()

	latency := time.Since(start)
	status := c.Writer.Status()
	method := c.Request.Method
	path := c.Request.URL.Path
	ip := c.ClientIP()

	if len(c.Errors) > 0 {

		log.Printf("ERROR: status=%d method=%s path=%s ip=%s latency=%v errors=%s",
			status, method, path, ip, latency, c.Errors.String())
	} else {

		log.Printf("INFO: status=%d method=%s path=%s ip=%s latency=%v",
			status, method, path, ip, latency)
	}
}
