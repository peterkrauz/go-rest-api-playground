package handler

import (
	"github.com/go-chi/render"
	"net/http"
)

type ErrorResponse struct {
	Error      error  `json:"-"`
	StatusCode int    `json:"-"`
	StatusText string `json:"status_text"`
	Message    string `json:"message"`
}

var (
	ErrorMethodNotAllowed = &ErrorResponse{StatusCode: 405, Message: "Method not allowed"}
	ErrorNotFound         = &ErrorResponse{StatusCode: 404, Message: "Resource not found"}
	ErrorBadRequest       = &ErrorResponse{StatusCode: 400, Message: "Bad request"}
)

func (response *ErrorResponse) Render(_ http.ResponseWriter, request *http.Request) error {
	render.Status(request, response.StatusCode)
	return nil
}

func ErrorRenderer(err error) *ErrorResponse {
	return &ErrorResponse{
		Error:      err,
		StatusCode: 400,
		StatusText: "Bad request",
		Message:    err.Error(),
	}
}

func ServerErrorRenderer(err error) *ErrorResponse {
	return &ErrorResponse{
		Error:      err,
		StatusCode: 500,
		StatusText: "Internal server error",
		Message:    err.Error(),
	}
}
