package handler

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/peterkrauz/go-rest-api-playground/db"
	"net/http"
)

var dbInstance db.Database

func NewHandler(db db.Database) http.Handler {
	router := chi.NewRouter()
	dbInstance = db

	router.MethodNotAllowed(methodNotAllowedHandler)
	router.NotFound(notFoundHandler)
	router.Route("/items", items)
	return router
}

func methodNotAllowedHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-type", "application/json")
	responseWriter.WriteHeader(405)
	render.Render(responseWriter, request, ErrorMethodNotAllowed)
}

func notFoundHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-type", "application/json")
	responseWriter.WriteHeader(400)
	render.Render(responseWriter, request, ErrorNotFound)
}
