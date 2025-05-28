package product

import (
	entityP "gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/domain/product/entity"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/transport/rest/product/product_dto"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/use_cases"
)

type Handler struct {
	ProductService *use_cases.ProductService
}

func NewProductHandler(productService *use_cases.ProductService) *Handler {
	log.Println("[ProductHandler] Initialized")
	return &Handler{ProductService: productService}
}

// SearchProducts выполняет поиск товаров
// @Summary Поиск товаров
// @Description Поиск товаров с фильтрацией, фасетами, сортировкой и подсветкой
// @Tags product
// @Accept json
// @Produce json
// @Param body body product_dto.SearchRequest true "Параметры поиска"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /product/search [post]
func (h *Handler) SearchProducts(c *gin.Context) {
	log.Println("[ProductHandler] SearchProducts called")

	var req product_dto.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ProductHandler][ERROR] Invalid request payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[ProductHandler] Payload: %+v", req)

	result, err := h.ProductService.SearchProducts(&req)
	if err != nil {
		log.Printf("[ProductHandler][ERROR] SearchProducts failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[ProductHandler] SearchProducts succeeded: found %d items", result.Total)
	c.JSON(http.StatusOK, result)
}

// GetClusterHealth возвращает статус кластера Elasticsearch
// @Summary Статус кластера Elasticsearch
// @Description Информация о состоянии (green, yellow, red) и метрики кластера
// @Tags product
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /product/health [get]
func (h *Handler) GetClusterHealth(c *gin.Context) {
	log.Println("[ProductHandler] GetClusterHealth called")

	health, err := h.ProductService.GetClusterHealth()
	if err != nil {
		log.Printf("[ProductHandler][ERROR] GetClusterHealth failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[ProductHandler] GetClusterHealth succeeded: status=%s", health.Status)
	c.JSON(http.StatusOK, health)
}

// CreateProduct индексирует один продукт
// @Summary Создать продукт
// @Description Индексация одного продукта
// @Tags product
// @Accept json
// @Produce json
// @Param body body product_dto.ProductRequest true "Данные продукта"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /product/create [post]
func (h *Handler) CreateProduct(c *gin.Context) {
	log.Println("[ProductHandler] CreateProduct called")

	var req product_dto.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ProductHandler][ERROR] Invalid payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[ProductHandler] Payload: %+v", req)

	p := &entityP.Product{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		Brand:       req.Brand,
	}

	if err := h.ProductService.IndexProduct(p); err != nil {
		log.Printf("[ProductHandler][ERROR] CreateProduct failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[ProductHandler] CreateProduct succeeded: ID=%s", p.ID)
	c.JSON(http.StatusOK, gin.H{"status": "created"})
}

// BulkCreateProducts индексирует сразу несколько продуктов
// @Summary Массовое создание продуктов
// @Description Индексация списка продуктов
// @Tags product
// @Accept json
// @Produce json
// @Param body body []product_dto.ProductRequest true "Список продуктов"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /product/bulk [post]
func (h *Handler) BulkCreateProducts(c *gin.Context) {
	log.Println("[ProductHandler] BulkCreateProducts called")

	var reqs []product_dto.ProductRequest
	if err := c.ShouldBindJSON(&reqs); err != nil {
		log.Printf("[ProductHandler][ERROR] Invalid payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[ProductHandler] Payload count: %d", len(reqs))

	products := make([]*entityP.Product, 0, len(reqs))
	for _, req := range reqs {
		products = append(products, &entityP.Product{
			ID:          req.ID,
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			Stock:       req.Stock,
			Category:    req.Category,
			Brand:       req.Brand,
		})
	}

	if err := h.ProductService.BulkIndexProducts(products); err != nil {
		log.Printf("[ProductHandler][ERROR] BulkCreateProducts failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[ProductHandler] BulkCreateProducts succeeded: indexed %d products", len(products))
	c.JSON(http.StatusOK, gin.H{"status": "bulk created"})
}

// DeleteProduct удаляет продукт по ID
// @Summary Удалить продукт
// @Description Удалить продукт из индекса по его ID
// @Tags product
// @Produce json
// @Param id path string true "ID продукта"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /product/delete/{id} [delete]
func (h *Handler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	log.Printf("[ProductHandler] DeleteProduct called: id=%s", id)
	if id == "" {
		log.Println("[ProductHandler][ERROR] ID is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter required"})
		return
	}

	if err := h.ProductService.DeleteProduct(id); err != nil {
		log.Printf("[ProductHandler][ERROR] DeleteProduct failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[ProductHandler] DeleteProduct succeeded: id=%s", id)
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
