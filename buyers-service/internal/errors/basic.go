package errors

import "net/http"

var (
	ErrInternalServer = new("bad request", "", http.StatusInternalServerError)
	ErrPgCreateBuyer  = new("bad request", "error during create buyer", http.StatusInternalServerError)
	ErrPgDeleteBuyer  = new("bad request", "error during deleting buyer", http.StatusInternalServerError)
	ErrBuyerNotFound  = new("bad request", "buyer with such id not found", http.StatusNotFound)
	ErrDuplicateEmail = new("bad request", "email already exists", http.StatusBadRequest)
)
