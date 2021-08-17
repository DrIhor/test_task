package graph

import "github.com/DrIhor/test_task/internal/models/items"

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	services items.ItemStorageServices
}

func NewResolver(storage items.ItemStorageServices) *Resolver {
	return &Resolver{
		services: storage,
	}
}
