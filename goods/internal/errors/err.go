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

// postrgess
// Ошибки, связанные с товарами
var ErrGoodNotFound = NewError(fasthttp.StatusNotFound, "error: good not found")
var ErrCreateGoodInternal = NewError(fasthttp.StatusInternalServerError, "error creating good")
var ErrUpdateGoodInternal = NewError(fasthttp.StatusInternalServerError, "error updating good")
var ErrDeleteGoodInternal = NewError(fasthttp.StatusInternalServerError, "error deleting good")
var ErrGoodAlreadyExists = NewError(fasthttp.StatusConflict, "error: good already exists")
var ErrInsufficientQuantity = NewError(fasthttp.StatusBadRequest, "error: insufficient quantity for operation")

// Ошибки, связанные с карточками товаров
var ErrGoodCardNotFound = NewError(fasthttp.StatusNotFound, "error: good card not found")
var ErrCreateGoodCardInternal = NewError(fasthttp.StatusInternalServerError, "error creating good card")
var ErrUpdateGoodCardInternal = NewError(fasthttp.StatusInternalServerError, "error updating good card")
var ErrDeleteGoodCardInternal = NewError(fasthttp.StatusInternalServerError, "error deleting good card")
var ErrGoodCardAlreadyExists = NewError(fasthttp.StatusConflict, "error: good card already exists")
var ErrUpdateGoodCardNotFields = NewError(fasthttp.StatusBadRequest, "error: no fields to update")

// Ошибки для функции AddCountGood
var (
	ErrAddCountGoodInternal = NewError(fasthttp.StatusInternalServerError, "error: internal server error while adding count to good")
	ErrAddCountGoodInvalid  = NewError(fasthttp.StatusBadRequest, "error: invalid count for adding to good")
)

// Ошибки для функции DeleteCountGood
var (
	ErrDeleteCountGoodInternal = NewError(fasthttp.StatusInternalServerError, "error: internal server error while deleting count from good")
	ErrNotEnoughQuantity       = NewError(fasthttp.StatusBadRequest, "error: not enough quantity to delete")
	ErrDeleteCountGoodInvalid  = NewError(fasthttp.StatusBadRequest, "error: invalid count for deleting from good")
)

var ErrReadCardInternal = NewError(fasthttp.StatusInternalServerError, "error: internal error read good")
