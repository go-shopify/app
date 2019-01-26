package app

import "net/http"

// ErrorHandler represents an error handler.
type ErrorHandler interface {
	ServeHTTPError(w http.ResponseWriter, req *http.Request, err error)
}
