package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/infrastructure/config/elasticsearch"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/infrastructure/repository/elasticsearch"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/transport/rest/routers"
	"log"
)

func RegisterRoutes(ctx context.Context, engine *gin.Engine, es *elasticsearch.Client, repo *repository.ProductRepository) {
	log.Println("[rest:product] Registering product routes")

	productRoute := routers.NewProductRoute(ctx, engine, es)

	productRoute.RegisterRoutes(repo)

	log.Println("[rest:product] Product routes registered successfully")
}
