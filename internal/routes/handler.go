package routes

import (
	"encoding/json"
	"net/http"

	itemModel "github.com/DrIhor/test_task/internal/models/items"
	msgServ "github.com/DrIhor/test_task/internal/service/messages"

	"github.com/gorilla/mux"
)

// типу краще ініціалізовувати сервіси для БД чи передавати їх як параметр
func HandlerItems(router *mux.Router, stor itemModel.ItemStorageServices) {
	router.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		ShowAllItems(w, r, stor)
	}).Methods("GET")
	router.HandleFunc("/items/{name}", func(w http.ResponseWriter, r *http.Request) {
		ShowItem(w, r, stor)
	}).Methods("GET")
	router.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		AddNewItem(w, r, stor)
	}).Methods("POST")
	router.HandleFunc("/items/{name}", func(w http.ResponseWriter, r *http.Request) {
		BuyItems(w, r, stor)
	}).Methods("PUT")
	router.HandleFunc("/items/{name}", func(w http.ResponseWriter, r *http.Request) {
		DeleteItem(w, r, stor)
	}).Methods("DELETE")
}

// CRUD implementation for all endpoints
// Read
func ShowAllItems(w http.ResponseWriter, r *http.Request, stor itemModel.ItemStorageServices) {

	res, err := stor.GetAllItems()
	if err != nil {
		res = msgServ.CreateMsgResp(err.Error())
		w.Write(res)
		return
	}

	w.Write(res)
}

func ShowItem(w http.ResponseWriter, r *http.Request, stor itemModel.ItemStorageServices) {
	params := mux.Vars(r)

	res, err := stor.GetItem(params["name"])
	if err != nil {
		res = msgServ.CreateMsgResp(err.Error())
		w.Write(res)
		return
	}

	w.Write(res)
}

// Create
func AddNewItem(w http.ResponseWriter, r *http.Request, stor itemModel.ItemStorageServices) {

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

	err = stor.AddNewItem(obj)
	if err != nil {
		res := msgServ.CreateMsgResp("Empty body. Change information")
		w.Write(res)
		return
	}

}

// Update
func BuyItems(w http.ResponseWriter, r *http.Request, stor itemModel.ItemStorageServices) {
	params := mux.Vars(r)

	res, err := stor.UpdateItem(params["name"])
	if err != nil {
		res := msgServ.CreateMsgResp("Empty body. Change information")
		w.Write(res)
		return
	}

	w.Write(res)
}

// Delete
func DeleteItem(w http.ResponseWriter, r *http.Request, stor itemModel.ItemStorageServices) {
	params := mux.Vars(r)

	err := stor.DeleteItem(params["name"])
	if err != nil {
		res := msgServ.CreateMsgResp("Empty body. Change information")
		w.Write(res)
		return
	}
}
