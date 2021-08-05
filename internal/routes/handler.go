package routes

import (
	"encoding/json"
	"net/http"
	"os"

	itemModel "github.com/DrIhor/test_task/internal/models/items"
	"github.com/DrIhor/test_task/internal/service/connectors"
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

	var errData error
	var res []byte
	switch os.Getenv("STORAGE_TYPE") {
	case "":
		res, errData = h.services.GetAllItems()
	case "grpc":
		grpcConn := connectors.NewGRPC(os.Getenv("GRCP_ADDR"))
		res, errData = grpcConn.GetAllItems()
	}

	if errData != nil {
		res = msgServ.CreateMsgResp(errData.Error())
		w.Write(res)
		return
	}

	w.Write(res)
}

func (h *HandlerItemsServ) ShowItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var errData error
	var res []byte
	switch os.Getenv("STORAGE_TYPE") {
	case "":
		res, errData = h.services.UpdateItem(params["name"])
	case "grpc":
		grpcConn := connectors.NewGRPC(os.Getenv("GRCP_ADDR"))
		res, errData = grpcConn.GetItem(params["name"])
	}

	if errData != nil {
		res = msgServ.CreateMsgResp(errData.Error())
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

	var errData error
	switch os.Getenv("STORAGE_TYPE") {
	case "":
		errData = h.services.AddNewItem(obj)
	case "grpc":
		grpcConn := connectors.NewGRPC(os.Getenv("GRCP_ADDR"))
		_, errData = grpcConn.AddNewItem(obj)
	}

	if errData != nil {
		res := msgServ.CreateMsgResp(errData.Error())
		w.Write(res)
		return
	}
}

// Update
func (h *HandlerItemsServ) BuyItems(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var errData error
	var res []byte
	switch os.Getenv("STORAGE_TYPE") {
	case "":
		res, errData = h.services.UpdateItem(params["name"])
	case "grpc":
		grpcConn := connectors.NewGRPC(os.Getenv("GRCP_ADDR"))
		res, errData = grpcConn.UpdateItem(params["name"])
	}

	if errData != nil {
		res := msgServ.CreateMsgResp(errData.Error())
		w.Write(res)
		return
	}

	w.Write(res)
}

// Delete
func (h *HandlerItemsServ) DeleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var errData error
	switch os.Getenv("STORAGE_TYPE") {
	case "":
		errData = h.services.DeleteItem(params["name"])
	case "grpc":
		grpcConn := connectors.NewGRPC(os.Getenv("GRCP_ADDR"))
		_, errData = grpcConn.DeleteItem(params["name"])
	}

	if errData != nil {
		res := msgServ.CreateMsgResp(errData.Error())
		w.Write(res)
		return
	}
}
