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
var ErrCreateSellerInternal = NewError(fasthttp.StatusInternalServerError, "error create seller")
var ErrCreateSellerFound = NewError(fasthttp.StatusFound, "error create seller: seller not new")

var ErrUpdateSellerInternal = NewError(fasthttp.StatusInternalServerError, "error update seller")
var ErrUpdateSellerNotFound = NewError(fasthttp.StatusInternalServerError, "error update seller: seller not found")
var ErrUpdateSellerNotFields = NewError(fasthttp.StatusInternalServerError, "error update seller: no fields to update")

var ErrDeleteSellerInternal = NewError(fasthttp.StatusInternalServerError, "error delete seller")
var ErrDeleteSellerNotFound = NewError(fasthttp.StatusNotFound, "error delete seller: seller not found")
