package models

import (
	"fmt"
	"net/http"
)

type Item struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

type ItemList struct {
	Items []Item `json:"items"`
}

func (item *Item) Bind(request *http.Request) error {
	if item.Name == "" {
		return fmt.Errorf("name cannot be null")
	}
	return nil
}

func (*ItemList) Render(responseWriter http.ResponseWriter, request *http.Request) error {
	return nil
}

func (*Item) Render(responseWriter http.ResponseWriter, request *http.Request) error {
	return nil
}
