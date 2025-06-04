package myErrors

import (
	"strconv"

	"github.com/valyala/fasthttp"
)

type Error struct {
	httpCode int
	cause    string
}

func NewError(code int, cause string) Error {
	return Error{
		httpCode: code,
		cause:    cause,
	}
}

func (e Error) GetHttpCode() int {
	return e.httpCode
}

func (e Error) GetCause() string {
	return e.cause
}
func (e Error) Error() string {
	return "Status code: " + strconv.Itoa(e.httpCode) + " cause: " + e.cause
}

// json
var ErrParseJSON = NewError(fasthttp.StatusBadRequest, "error decoding json")

// db
// Ошибки для операций с заказами
var ErrOrderAlreadyExists = NewError(fasthttp.StatusConflict, "error: order already exists with the same customer ID and status")
var ErrCreateOrderInternal = NewError(fasthttp.StatusInternalServerError, "error: internal error occurred while creating order")
var ErrUpdateOrderInternal = NewError(fasthttp.StatusInternalServerError, "error: internal error occurred while updating order")
var ErrDeleteOrderInternal = NewError(fasthttp.StatusInternalServerError, "error: internal error occurred while deleting order")
var ErrGetListOrdersInternal = NewError(fasthttp.StatusInternalServerError, "error: internal error occurred while retrieving list of orders")
var ErrGetOrderInternal = NewError(fasthttp.StatusInternalServerError, "error: internal error occurred while retrieving order")
var ErrOrderNotFound = NewError(fasthttp.StatusNotFound, "error: order not found")

// Ошибки для операций с позициями заказа
var ErrOrderItemAlreadyExists = NewError(fasthttp.StatusConflict, "error: order item already exists with the same order ID and good UUID")
var ErrCreateOrderItemInternal = NewError(fasthttp.StatusInternalServerError, "error: internal error occurred while creating order item")
var ErrUpdateOrderItemInternal = NewError(fasthttp.StatusInternalServerError, "error: internal error occurred while updating order item")
var ErrDeleteOrderItemInternal = NewError(fasthttp.StatusInternalServerError, "error: internal error occurred while deleting order item")
var ErrUpdateOrderItemQuantityInternal = NewError(fasthttp.StatusInternalServerError, "error: internal error occurred while updating order item quantity")
var ErrGetOrderItemsInternal = NewError(fasthttp.StatusInternalServerError, "error: internal error occurred while retrieving order items")
var ErrOrderItemNotFound = NewError(fasthttp.StatusNotFound, "error: order item not found")
var ErrInvalidOrderStatus = NewError(fasthttp.StatusBadRequest, "error: invalid order status for the operation")

// Ошибки для операций с корзиной
var ErrBagNotFound = NewError(fasthttp.StatusNotFound, "error: cart not found")
var ErrBagEmpty = NewError(fasthttp.StatusBadRequest, "error: cart is empty")
var ErrUnauthorized = NewError(fasthttp.StatusUnauthorized, "error: unauthorized access")

// Ошибки для операций с товарами
var ErrGoodNotFound = NewError(fasthttp.StatusNotFound, "error: product not found")
var ErrCreateGoodInternal = NewError(fasthttp.StatusInternalServerError, "error: internal error occurred while creating product")
var ErrGetGoodInternal = NewError(fasthttp.StatusInternalServerError, "error: internal error occurred while retrieving product")
