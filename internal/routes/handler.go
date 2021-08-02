package routes

import (
	"encoding/json"
	"net/http"

	itemModel "github.com/DrIhor/test_task/internal/models/items"
	itemServ "github.com/DrIhor/test_task/internal/service/items"

	msgServ "github.com/DrIhor/test_task/internal/service/messages"

	"github.com/gorilla/mux"
)

type HandlerItemsServ struct {
	router   *mux.Router
	services *itemServ.ItemServices
}

func New(router *mux.Router, stor itemModel.ItemStorageServices) *HandlerItemsServ {
	return &HandlerItemsServ{
		router:   router,
		services: itemServ.New(stor),
	}
}

// типу краще ініціалізовувати сервіси для БД чи передавати їх як параметр
func (h *HandlerItemsServ) HandlerItems() {
	h.router.HandleFunc("/items", h.ShowAllItems).Methods("GET")
	h.router.HandleFunc("/items/{name}", h.ShowItem).Methods("GET")
	h.router.HandleFunc("/items", h.AddNewItem).Methods("POST")
	h.router.HandleFunc("/items/{name}", h.BuyItems).Methods("PUT")
	h.router.HandleFunc("/items/{name}", h.DeleteItem).Methods("DELETE")
}

// CRUD implementation for all endpoints
// Read
func (h *HandlerItemsServ) ShowAllItems(w http.ResponseWriter, r *http.Request) {

	res, err := h.services.GetAllItems()
	if err != nil {
		res = msgServ.CreateMsgResp(err.Error())
		w.Write(res)
		return
	}

	w.Write(res)
}

func (h *HandlerItemsServ) ShowItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	res, err := h.services.GetItem(params["name"])
	if err != nil {
		res = msgServ.CreateMsgResp(err.Error())
		w.Write(res)
		return
	}

	w.Write(res)
}

// Create
func (h *HandlerItemsServ) AddNewItem(w http.ResponseWriter, r *http.Request) {

	var obj itemModel.Item
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		res := msgServ.CreateMsgResp(err.Error())
		w.Write(res)
		return
	}

	if (obj == itemModel.Item{}) {
		res := msgServ.CreateMsgResp("Empty body. Change information")
		w.Write(res)
		return
	}

	err = h.services.AddNewItem(obj)
	if err != nil {
		res := msgServ.CreateMsgResp(err.Error())
		w.Write(res)
		return
	}

}

// Update
func (h *HandlerItemsServ) BuyItems(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	res, err := h.services.UpdateItem(params["name"])
	if err != nil {
		res := msgServ.CreateMsgResp(err.Error())
		w.Write(res)
		return
	}

	w.Write(res)
}

// Delete
func (h *HandlerItemsServ) DeleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	err := h.services.DeleteItem(params["name"])
	if err != nil {
		res := msgServ.CreateMsgResp(err.Error())
		w.Write(res)
		return
	}
}
