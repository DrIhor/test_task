package graph

//go:generate go run github.com/99designs/gqlgen

import (
	itemsModel "github.com/DrIhor/test_task/internal/models/items"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	services itemsModel.ItemStorageServices
}

func NewResolver(storage itemsModel.ItemStorageServices) *Resolver {
	return &Resolver{
		services: storage,
	}
}
