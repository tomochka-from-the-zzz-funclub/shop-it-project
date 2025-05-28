package transport

import (
	"encoding/json"
	"fmt"

	config "market/internal/cfg"
	myLog "market/internal/logger"
	"market/internal/models"
	service "market/internal/services"
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

	hb.rout.POST("/sellers/create", hb.HandleSellersCreate()) //общее количество просмотров товаров

	hb.rout.PUT("/sellers/update/{uuid}", hb.HandleUpdateSeller()) // — обновление информации о продавце

	hb.rout.DELETE("/sellers/delete", hb.HandleDeleteSeller()) // удаление продавца

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

func (hb *HandlersBuilder) HandleSellersCreate() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		myLog.Log.Debugf("Start func HandleSellersCreate")
		if ctx.IsPost() {
			// Создаем экземпляр структуры Seller
			var seller models.Seller

			// Считываем и декодируем JSON из тела запроса
			if err := json.Unmarshal(ctx.PostBody(), &seller); err != nil {
				httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid request body")
				fmt.Println(err.Error())
				return
			}

			id, err := hb.srv.Create(seller)
			if err != nil {
				httpErrorResponse(ctx, fasthttp.StatusInternalServerError, "Failed to create seller")
				return
			}

			// Если все прошло успешно, выводим ID продавца в тело ответа
			ctx.SetStatusCode(fasthttp.StatusCreated) // Устанавливаем статус код 201 Created
			ctx.SetContentType("application/json")
			response := map[string]interface{}{
				"id": id.String(), // Если id - это UUID, преобразуем его в строку
			}
			jsonResponse(ctx, response)
		}
	}, "HandleSellersCreate")
}

func (hb *HandlersBuilder) HandleUpdateSeller() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		myLog.Log.Debugf("Start func HandleUpdateSeller")

		// Проверяем, что метод запроса - PUT или PATCH (в зависимости от вашей структуры)
		if ctx.IsPut() {
			var seller models.Seller

			// Считываем и декодируем JSON из тела запроса
			if err := json.Unmarshal(ctx.PostBody(), &seller); err != nil {
				httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid request body")
				return
			}

			// Получаем ID из параметров URL
			id, ok := ctx.UserValue("uuid").(string) // Предполагается, что ID передаётся как параметр
			if !ok {
				myLog.Log.Errorf("Error unknown type id")
			}
			fmt.Println(id)
			sellerID, err := uuid.Parse(fmt.Sprintf("%v", id)) // Преобразование ID в UUID
			fmt.Println(sellerID)
			if err != nil {
				httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid seller ID")
				fmt.Println(err.Error())
				return
			}

			// Вызываем метод обновления
			err = hb.srv.Update(sellerID, seller)
			if err != nil {
				httpErrorResponse(ctx, fasthttp.StatusInternalServerError, "Failed to update seller")
				return
			}

			// Устанавливаем статус код 204 No Content для успешного обновления
			ctx.SetStatusCode(fasthttp.StatusNoContent)
			return
		}

		// Если метод не подходит, отвечаем ошибкой
		httpErrorResponse(ctx, fasthttp.StatusMethodNotAllowed, "Method not allowed")
	}, "HandleUpdateSeller")
}

func (hb *HandlersBuilder) HandleDeleteSeller() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		myLog.Log.Debugf("Start func HandleDeleteSeller")

		// Проверяем, что метод запроса - DELETE
		if ctx.IsDelete() {
			// Получаем ID продавца из параметров URL
			id := ctx.UserValue("uuid").(string) // ID передается как параметр
			fmt.Println(id)
			sellerID, err := uuid.Parse(fmt.Sprintf("%v", id)) // Преобразование ID в UUID
			fmt.Println(sellerID)
			if err != nil {
				httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid seller ID")
				return
			}

			// Вызываем метод удаления
			err = hb.srv.Delete(sellerID)
			if err != nil {
				// Если продавец не найден, возвращаем 404 Not Found
				if err.Error() == "not found" {
					httpErrorResponse(ctx, fasthttp.StatusNotFound, fmt.Sprintf("Seller with ID %s not found", sellerID))
					return
				}
				// В случае других ошибок возвращаем 500 Internal Server Error
				httpErrorResponse(ctx, fasthttp.StatusInternalServerError, "Failed to delete seller")
				return
			}

			// Устанавливаем статус код 204 No Content для успешного удаления
			ctx.SetStatusCode(fasthttp.StatusNoContent)
			return
		}

		// Если метод не подходит, отвечаем ошибкой
		httpErrorResponse(ctx, fasthttp.StatusMethodNotAllowed, "Method not allowed")
	}, "HandleDeleteSeller")
}
