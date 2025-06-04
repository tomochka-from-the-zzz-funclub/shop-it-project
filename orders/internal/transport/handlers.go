package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"orders/internal/config"
	myErrors "orders/internal/errors"
	"orders/internal/models"
	"orders/internal/services"

	"github.com/fasthttp/router"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
)

type HandlersBuilder struct {
	srv       services.InterfaceService
	jwtSecret string
	rout      *router.Router
}

func HandleCreate(cfg config.Config, s services.InterfaceService) {
	hb := HandlersBuilder{
		srv:       s,
		jwtSecret: cfg.JwtSecret,
		rout:      router.New(),
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8090", nil)
	}()

	hb.rout.POST("/order/create", hb.HandleCreateOrder())
	hb.rout.GET("/order/get/:id", hb.HandleGetOrder())
	hb.rout.GET("/order/info", hb.HandleListOrders())
	hb.rout.POST("/cart/add", hb.HandleAddToCart())
	hb.rout.POST("/cart/remove", hb.HandleDeleteFromCart())
	hb.rout.POST("/order/create-from-cart", hb.HandleCreateOrderFromCart())

	fmt.Println(fasthttp.ListenAndServe(":8083", hb.rout.Handler))
}

// jsonResponse sends a JSON response
func jsonResponse(ctx *fasthttp.RequestCtx, response interface{}) {
	respBody, err := json.Marshal(response)
	if err != nil {
		httpErrorResponse(ctx, fasthttp.StatusInternalServerError, "Failed to create response")
		return
	}
	ctx.SetContentType("application/json")
	ctx.SetBody(respBody)
}

// httpErrorResponse sends an error response
func httpErrorResponse(ctx *fasthttp.RequestCtx, statusCode int, message string) {
	ctx.SetStatusCode(statusCode)
	ctx.SetContentType("application/json")
	response := map[string]string{"error": message}
	jsonResponse(ctx, response)
}

// extractCustomerID extracts customer_id from Authorization header
func (hb *HandlersBuilder) extractCustomerID(ctx *fasthttp.RequestCtx) (uuid.UUID, error) {
	authHeader := string(ctx.Request.Header.Peek("Authorization"))
	if authHeader == "" {
		return uuid.UUID{}, errors.New("missing authorization header")
	}
	if len(authHeader) > 7 && strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := authHeader[7:]
		return hb.Parse(tokenString)
	}
	return uuid.UUID{}, errors.New("invalid authorization header format")
}

func (hb *HandlersBuilder) HandleCreateOrder() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if !ctx.IsPost() {
			httpErrorResponse(ctx, fasthttp.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		customerID, err := hb.extractCustomerID(ctx)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusUnauthorized, err.Error())
			return
		}

		var order models.Order
		if err := json.Unmarshal(ctx.PostBody(), &order); err != nil {
			httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid request body")
			return
		}
		order.CustomerID = customerID

		id, err := hb.srv.SrvCreateOrder(order)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusInternalServerError, "Failed to create order")
			return
		}

		ctx.SetStatusCode(fasthttp.StatusCreated)
		jsonResponse(ctx, map[string]string{"id": id.String()})
	}, "HandleCreateOrder")
}

func (hb *HandlersBuilder) HandleGetOrder() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if !ctx.IsGet() {
			httpErrorResponse(ctx, fasthttp.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		customerID, err := hb.extractCustomerID(ctx)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusUnauthorized, err.Error())
			return
		}

		idStr := ctx.UserValue("id").(string)
		id, err := uuid.Parse(idStr)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid UUID format")
			return
		}

		order, goods, err := hb.srv.SrvGetOrderByID(customerID, id)
		if err != nil {
			if err == myErrors.ErrOrderNotFound {
				httpErrorResponse(ctx, fasthttp.StatusNotFound, "Order not found")
			} else {
				httpErrorResponse(ctx, fasthttp.StatusInternalServerError, "Internal server error")
			}
			return
		}

		if order.CustomerID != customerID {
			httpErrorResponse(ctx, fasthttp.StatusUnauthorized, "Unauthorized access to order")
			return
		}

		response := struct {
			Order models.Order  `json:"order"`
			Goods []models.Good `json:"goods"`
		}{
			Order: order,
			Goods: goods,
		}

		jsonResponse(ctx, response)
	}, "HandleGetOrder")
}

func (hb *HandlersBuilder) HandleListOrders() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if !ctx.IsGet() {
			httpErrorResponse(ctx, fasthttp.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		customerID, err := hb.extractCustomerID(ctx)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusUnauthorized, err.Error())
			return
		}

		status := string(ctx.QueryArgs().Peek("status"))
		limit := ctx.QueryArgs().GetUintOrZero("limit")
		offset := ctx.QueryArgs().GetUintOrZero("offset")

		orders, err := hb.srv.SrvGetListOrders(status, limit, offset, customerID)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusInternalServerError, "Failed to retrieve orders")
			return
		}

		// Filter orders by customer_id
		var filteredOrders []models.Order
		for _, order := range orders {
			if order.CustomerID == customerID {
				filteredOrders = append(filteredOrders, order)
			}
		}

		jsonResponse(ctx, filteredOrders)
	}, "HandleListOrders")
}

func (hb *HandlersBuilder) HandleAddToCart() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if !ctx.IsPost() {
			httpErrorResponse(ctx, fasthttp.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		customerID, err := hb.extractCustomerID(ctx)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusUnauthorized, err.Error())
			return
		}

		var request struct {
			Product string `json:"product"`
		}
		if err := json.Unmarshal(ctx.PostBody(), &request); err != nil {
			httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid request body")
			return
		}

		parsedProductUUID, err := uuid.Parse(request.Product)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid Product format")
			return
		}

		bagID, err := hb.srv.GetBagIDByUser(customerID)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusNotFound, "User's bag not found")
			return
		}

		if err := hb.srv.AddGoodToBag(bagID, parsedProductUUID, customerID); err != nil {
			httpErrorResponse(ctx, fasthttp.StatusInternalServerError, "Failed to add product to cart")
			return
		}

		ctx.SetStatusCode(fasthttp.StatusOK)
		jsonResponse(ctx, map[string]string{"message": "Product added to cart"})
	}, "HandleAddToCart")
}

func (hb *HandlersBuilder) HandleDeleteFromCart() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if !ctx.IsPost() {
			httpErrorResponse(ctx, fasthttp.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		customerID, err := hb.extractCustomerID(ctx)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusUnauthorized, err.Error())
			return
		}

		var request struct {
			Product string `json:"product"`
		}
		if err := json.Unmarshal(ctx.PostBody(), &request); err != nil {
			httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid request body")
			return
		}

		productUUID, err := uuid.Parse(request.Product)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusBadRequest, "Invalid Product format")
			return
		}

		bagID, err := hb.srv.GetBagIDByUser(customerID)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusNotFound, "User's bag not found")
			return
		}

		if err := hb.srv.RemoveGoodFromBag(bagID, productUUID, customerID); err != nil {
			httpErrorResponse(ctx, fasthttp.StatusInternalServerError, "Failed to remove product from cart")
			return
		}

		ctx.SetStatusCode(fasthttp.StatusOK)
		jsonResponse(ctx, map[string]string{"message": "Product removed from cart"})
	}, "HandleDeleteFromCart")
}

func (hb *HandlersBuilder) HandleCreateOrderFromCart() func(ctx *fasthttp.RequestCtx) {
	return metrics(func(ctx *fasthttp.RequestCtx) {
		if !ctx.IsPost() {
			httpErrorResponse(ctx, fasthttp.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		customerID, err := hb.extractCustomerID(ctx)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusUnauthorized, err.Error())
			return
		}

		bagID, err := hb.srv.GetBagIDByUser(customerID)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusNotFound, "User's bag not found")
			return
		}

		order, err := hb.srv.CheckoutBag(bagID, customerID)
		if err != nil {
			httpErrorResponse(ctx, fasthttp.StatusInternalServerError, "Failed to create order from cart")
			return
		}

		ctx.SetStatusCode(fasthttp.StatusCreated)
		jsonResponse(ctx, map[string]string{"order_id": order.UUID.String()})
	}, "HandleCreateOrderFromCart")
}
