package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"search/internal/config"
	"search/internal/models"
	"search/internal/service"
	"strconv"

	"github.com/fasthttp/router"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
)

type HandlersBuilder struct {
	srv  service.InterfaceSearchServicee
	rout *router.Router
}

func HandleCreate(cfg config.Config, s service.InterfaceSearchServicee) {

	hb := HandlersBuilder{
		srv:  s,
		rout: router.New(),
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8090", nil)
	}()
	//admin work methods
	hb.rout.GET("/goods/search", hb.HandleSearchGoods())

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

func (hb *HandlersBuilder) HandleSearchGoods() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if !ctx.IsGet() {
			httpErrorResponse(ctx, fasthttp.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		query := string(ctx.URI().QueryArgs().Peek("query"))
		minPriceStr := string(ctx.URI().QueryArgs().Peek("min_price"))
		maxPriceStr := string(ctx.URI().QueryArgs().Peek("max_price"))
		sortBy := string(ctx.URI().QueryArgs().Peek("sort_by"))
		pageStr := string(ctx.URI().QueryArgs().Peek("page"))
		pageSizeStr := string(ctx.URI().QueryArgs().Peek("page_size"))

		var minPrice, maxPrice float64
		var page, pageSize int
		var err error

		if minPriceStr != "" {
			minPrice, err = strconv.ParseFloat(minPriceStr, 64)
			if err != nil {
				httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid min_price")
				return
			}
		}
		if maxPriceStr != "" {
			maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
			if err != nil {
				httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid max_price")
				return
			}
		}
		if pageStr != "" {
			page, _ = strconv.Atoi(pageStr)
		}
		if pageSizeStr != "" {
			pageSize, _ = strconv.Atoi(pageSizeStr)
		}
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}

		req := models.SearchRequest{
			Query:    query,
			MinPrice: minPrice,
			MaxPrice: maxPrice,
			SortBy:   sortBy,
			Page:     page,
			PageSize: pageSize,
		}

		resp, err := hb.srv.SrvSearchGoods(context.Background(), req)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusInternalServerError, "Search failed")
			return
		}
		jsonResponse(ctx, resp)
	}, "HandleSearchGoods")
}
