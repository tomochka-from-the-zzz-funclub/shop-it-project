package transport

import (
	"encoding/json"
	"fmt"

	config "goods/internal/cfg"
	"goods/internal/models"
	service "goods/internal/services"
	"net/http"

	"github.com/fasthttp/router"
	"github.com/google/uuid"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
)

type HandlersBuilder struct {
	srv  service.InterfaceService
	rout *router.Router
}

func HandleCreate(cfg config.Config, s service.InterfaceService) {

	hb := HandlersBuilder{
		srv:  s,
		rout: router.New(),
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8090", nil)
	}()

	// Создать новую карточку товара
	hb.rout.POST("/goodcards/create", hb.HandleCreateGoodCard())

	// Удалить карточку товара по UUID
	hb.rout.DELETE("/goodcards/{id}", hb.HandleDeleteGoodCard())

	// Обновить карточку товара по UUID
	hb.rout.PUT("/goodcards/{id}", hb.HandleUpdateGoodCard())

	// Создать товар с привязкой к карточке товара
	hb.rout.POST("/goods/create", hb.HandleCreateGood())

	// Удалить товар по UUID
	hb.rout.DELETE("/goods/{id}", hb.HandleDeleteGood())

	// Увеличить количество товара по UUID
	hb.rout.POST("/goods/{id}/add", hb.HandleAddCountGood())

	// Уменьшить количество товара по UUID
	hb.rout.POST("/goods/{id}/remove", hb.HandleDeleteCountGood())

	// Получить информацию о товаре и его карточке по UUID товара
	hb.rout.GET("/goods/{id}", hb.HandleReadCard())

	fmt.Println(fasthttp.ListenAndServe(":8080", hb.rout.Handler))
}

// Вспомогательная функция для отправки JSON ответа
func jsonResponse(ctx *fasthttp.RequestCtx, response interface{}) {
	respBody, err := json.Marshal(response)
	if err != nil {
		httpErrorResponse(ctx, fasthttp.StatusInternalServerError, "Failed to create response")
		return
	}
	ctx.SetContentType("application/json")
	ctx.SetBody(respBody)
}

// Вспомогательная функция для отправки ошибки
func httpErrorResponse(ctx *fasthttp.RequestCtx, statusCode int, message string) {
	ctx.SetStatusCode(statusCode)
	ctx.SetContentType("application/json")
	response := map[string]string{"error": message}
	jsonResponse(ctx, response)
}

func (hb *HandlersBuilder) HandleCreateGoodCard() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if !ctx.IsPost() {
			httpErrorResponse(ctx, fasthttp.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		var card models.GoodCard
		if err := json.Unmarshal(ctx.PostBody(), &card); err != nil {
			httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid request body")
			return
		}

		id, err := hb.srv.SrvCreateGoodCard(card)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusInternalServerError, "Failed to create good card")
			return
		}

		ctx.SetStatusCode(fasthttp.StatusCreated)
		jsonResponse(ctx, map[string]string{"id": id.String()})
	}, "HandleCreateGoodCard")
}

func (hb *HandlersBuilder) HandleDeleteGoodCard() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if !ctx.IsDelete() {
			httpErrorResponse(ctx, fasthttp.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		idStr := ctx.UserValue("id").(string)
		id, err := uuid.Parse(idStr)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid UUID")
			return
		}

		if err := hb.srv.SrvDeleteGoodCard(id); err != nil {
			httpErrorResponse(ctx, fasthttp.StatusInternalServerError, "Failed to delete good card")
			return
		}

		ctx.SetStatusCode(fasthttp.StatusNoContent) // 204 No Content
	}, "HandleDeleteGoodCard")
}

func (hb *HandlersBuilder) HandleUpdateGoodCard() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if !ctx.IsPut() {
			httpErrorResponse(ctx, fasthttp.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		idStr := ctx.UserValue("id").(string)
		id, err := uuid.Parse(idStr)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid UUID")
			return
		}

		var card models.GoodCard
		if err := json.Unmarshal(ctx.PostBody(), &card); err != nil {
			httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid request body")
			return
		}

		if err := hb.srv.SrvUpdateGoodCard(id, card); err != nil {
			httpErrorResponse(ctx, fasthttp.StatusInternalServerError, "Failed to update good card")
			return
		}

		ctx.SetStatusCode(fasthttp.StatusOK)
		jsonResponse(ctx, map[string]string{"status": "updated"})
	}, "HandleUpdateGoodCard")
}

func (hb *HandlersBuilder) HandleReadCard() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if !ctx.IsGet() {
			httpErrorResponse(ctx, fasthttp.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		idStr := ctx.UserValue("id").(string)
		id, err := uuid.Parse(idStr)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid UUID")
			return
		}

		card, err := hb.srv.SrvReadGood(id)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusInternalServerError, "Failed to read good card")
			return
		}

		ctx.SetStatusCode(fasthttp.StatusOK)
		jsonResponse(ctx, card)
	}, "HandleReadCard")
}

func (hb *HandlersBuilder) HandleAddCountGood() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if !ctx.IsPost() {
			httpErrorResponse(ctx, fasthttp.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		idStr := ctx.UserValue("id").(string)
		id, err := uuid.Parse(idStr)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid UUID")
			return
		}

		var req struct {
			Number int `json:"number"`
		}
		if err := json.Unmarshal(ctx.PostBody(), &req); err != nil || req.Number <= 0 {
			httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid number")
			return
		}

		newCount, err := hb.srv.SrvAddCountGood(id, req.Number)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusInternalServerError, "Failed to add count")
			return
		}

		ctx.SetStatusCode(fasthttp.StatusOK)
		jsonResponse(ctx, map[string]interface{}{"new_count": newCount})
	}, "HandleAddCountGood")
}

func (hb *HandlersBuilder) HandleDeleteCountGood() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if !ctx.IsPost() {
			httpErrorResponse(ctx, fasthttp.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		idStr := ctx.UserValue("id").(string)
		id, err := uuid.Parse(idStr)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid UUID")
			return
		}

		var req struct {
			Number int `json:"number"`
		}
		if err := json.Unmarshal(ctx.PostBody(), &req); err != nil || req.Number <= 0 {
			httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid number")
			return
		}

		newCount, err := hb.srv.SrvDeleteCountGood(id, req.Number)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusInternalServerError, "Failed to delete count")
			return
		}

		ctx.SetStatusCode(fasthttp.StatusOK)
		jsonResponse(ctx, map[string]interface{}{"new_count": newCount})
	}, "HandleDeleteCountGood")
}

func (hb *HandlersBuilder) HandleCreateGood() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if !ctx.IsPost() {
			httpErrorResponse(ctx, fasthttp.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		var card models.GoodCard
		if err := json.Unmarshal(ctx.PostBody(), &card); err != nil {
			httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid request body")
			return
		}

		// Создаём карточку товара, что является созданием товара
		id, err := hb.srv.SrvCreateGoodCard(card)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusInternalServerError, "Failed to create good")
			return
		}

		ctx.SetStatusCode(fasthttp.StatusCreated)
		jsonResponse(ctx, map[string]string{"id": id.String()})
	}, "HandleCreateGood")
}

func (hb *HandlersBuilder) HandleDeleteGood() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if !ctx.IsDelete() {
			httpErrorResponse(ctx, fasthttp.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		idStr := ctx.UserValue("id").(string)
		id, err := uuid.Parse(idStr)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid UUID")
			return
		}

		// Удаляем карточку товара по cardId
		if err := hb.srv.SrvDeleteGoodCard(id); err != nil {
			httpErrorResponse(ctx, fasthttp.StatusInternalServerError, "Failed to delete good")
			return
		}

		ctx.SetStatusCode(fasthttp.StatusNoContent)
	}, "HandleDeleteGood")
}
