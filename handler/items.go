package handler

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/peterkrauz/go-rest-api-playground/db"
	"github.com/peterkrauz/go-rest-api-playground/models"
	"net/http"
	"strconv"
)

var itemIdKey = "itemId"

func items(router chi.Router) {
	router.Get("/", getAllItems)
	router.Post("/", createItem)
	router.Route("/{itemId}/", detailRouteResolver)
}

func detailRouteResolver(router chi.Router) {
	router.Use(ItemContext)
	router.Get("/", getItem)
	router.Put("/", updateItem)
	router.Delete("/", deleteItem)
}

func ItemContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		itemId := chi.URLParam(request, "itemId")
		if itemId == "" {
			render.Render(responseWriter, request, ErrorRenderer(fmt.Errorf("item Id is required")))
			return
		}

		id, err := strconv.Atoi(itemId)
		if err != nil {
			render.Render(responseWriter, request, ErrorRenderer(fmt.Errorf("invalid item Id")))
		}
		ctx := context.WithValue(request.Context(), itemIdKey, id)
		next.ServeHTTP(responseWriter, request.WithContext(ctx))
	})
}

func createItem(responseWriter http.ResponseWriter, request *http.Request) {
	item := &models.Item{}

	if err := render.Bind(request, item); err != nil {
		render.Render(responseWriter, request, ErrorBadRequest)
		return
	}

	if err := dbInstance.CreateItem(item); err != nil {
		render.Render(responseWriter, request, ErrorRenderer(err))
		return
	}

	if err := render.Render(responseWriter, request, item); err != nil {
		render.Render(responseWriter, request, ServerErrorRenderer(err))
		return
	}
}

func getAllItems(responseWriter http.ResponseWriter, request *http.Request) {
	items, err := dbInstance.GetAllItems()

	if err != nil {
		render.Render(responseWriter, request, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(responseWriter, request, items); err != nil {
		render.Render(responseWriter, request, ErrorRenderer(err))
	}
}

func getItem(responseWriter http.ResponseWriter, request *http.Request) {
	itemId := request.Context().Value(itemIdKey).(int)
	item, err := dbInstance.GetItemById(itemId)

	if err != nil {
		if err == db.ErrorNoItemFound {
			render.Render(responseWriter, request, ErrorNotFound)
		} else {
			render.Render(responseWriter, request, ErrorRenderer(err))
		}
		return
	}

	if err := render.Render(responseWriter, request, &item); err != nil {
		render.Render(responseWriter, request, ServerErrorRenderer(err))
		return
	}
}

func deleteItem(responseWriter http.ResponseWriter, request *http.Request) {
	itemId := request.Context().Value(itemIdKey).(int)
	err := dbInstance.DeleteItem(itemId)
	if err != nil {
		if err == db.ErrorNoItemFound {
			render.Render(responseWriter, request, ErrorNotFound)
		} else {
			render.Render(responseWriter, request, ServerErrorRenderer(err))
		}
		return
	}
}

func updateItem(responseWriter http.ResponseWriter, request *http.Request) {
	itemId := request.Context().Value(itemIdKey).(int)
	itemData := models.Item{}

	if err := render.Bind(request, &itemData); err != nil {
		render.Render(responseWriter, request, ErrorBadRequest)
		return
	}

	item, err := dbInstance.UpdateItem(itemId, itemData)
	if err != nil {
		if err == db.ErrorNoItemFound {
			render.Render(responseWriter, request, ErrorNotFound)
		} else {
			render.Render(responseWriter, request, ServerErrorRenderer(err))
		}
	}

	if err := render.Render(responseWriter, request, &item); err != nil {
		render.Render(responseWriter, request, ServerErrorRenderer(err))
		return
	}
}
