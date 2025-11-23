package errcode

import "net/http"

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}

// Common Errors
var (
	ErrInternalServer = &Error{Code: http.StatusInternalServerError, Message: "Internal Server Error"}
	ErrInvalidRequest = &Error{Code: http.StatusBadRequest, Message: "Invalid Request"}
	ErrUnauthorized   = &Error{Code: http.StatusUnauthorized, Message: "Unauthorized"}
	ErrNotFound       = &Error{Code: http.StatusNotFound, Message: "Not Found"}
)

func New(code int, msg string) *Error {
	return &Error{Code: code, Message: msg}
}
