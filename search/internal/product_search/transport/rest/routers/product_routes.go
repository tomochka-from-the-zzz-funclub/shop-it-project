package routers

import (
	"context"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/infrastructure/config/elasticsearch"
	"log"

	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/infrastructure/repository/elasticsearch"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/transport/rest/product"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/use_cases"

	"github.com/gin-gonic/gin"
)

type ProductRoute struct {
	engine   *gin.Engine
	ctx      context.Context
	ESClient *elasticsearch.Client
}

func NewProductRoute(ctx context.Context, engine *gin.Engine, esClient *elasticsearch.Client) *ProductRoute {
	log.Println("[routers:product] initializing ProductRoute")
	return &ProductRoute{ctx: ctx, engine: engine, ESClient: esClient}
}

func (pr *ProductRoute) RegisterRoutes(repo *repository.ProductRepository) {

	service := use_cases.NewProductService(repo)
	h := product.NewProductHandler(service)

	api := pr.engine.Group("/product")
	{
		api.POST("/search", h.SearchProducts)
		api.GET("/health", h.GetClusterHealth)
		api.POST("/create", h.CreateProduct)
		api.POST("/bulk", h.BulkCreateProducts)
		api.DELETE("/delete/:id", h.DeleteProduct)

	}

	log.Println("[routers:product] routes registered: POST /product/search, GET /product/health")
}
