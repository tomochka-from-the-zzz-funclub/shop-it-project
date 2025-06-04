package errors

import "net/http"

var (
	ErrInternalServer = new("bad request", "", http.StatusInternalServerError)
	ErrPgCreateSeller = new("bad request", "error during create seller", http.StatusInternalServerError)
	ErrPgDeleteSeller = new("bad request", "error during deleting seller", http.StatusInternalServerError)
	ErrSellerNotFound = new("bad request", "buyer with such id not found", http.StatusNotFound)
	ErrDuplicateEmail = new("bad request", "email already exists", http.StatusBadRequest)
)
