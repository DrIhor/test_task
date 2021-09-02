package routes

import (
	"context"

	itemModel "github.com/DrIhor/test_task/internal/models/items"
	itemServ "github.com/DrIhor/test_task/internal/service/items"
	"github.com/gorilla/mux"
)

type HandlerItemsServ struct {
	router   *mux.Router
	services itemServ.ItemSrv
	ctx      context.Context
}

func New(ctx context.Context, router *mux.Router, stor itemModel.ItemStorageServices) *HandlerItemsServ {
	return &HandlerItemsServ{
		router:   router,
		services: itemServ.New(stor),
		ctx:      ctx,
	}
}

func (h *HandlerItemsServ) HandlerItems() {
	// item functions
	h.router.HandleFunc("/items", h.ShowAllItems).Methods("GET")
	h.router.HandleFunc("/item", h.ShowItem).Methods("GET")
	h.router.HandleFunc("/item", h.AddNewItem).Methods("POST")
	h.router.HandleFunc("/item", h.BuyItems).Methods("PUT")
	h.router.HandleFunc("/item", h.DeleteItem).Methods("DELETE")

	// work with files
	h.router.HandleFunc("/items/csv", h.ReturnAllItemsCSV).Methods("GET")
	h.router.HandleFunc("/items/csv", h.AddDataFromCSV).Methods("POST")
}
