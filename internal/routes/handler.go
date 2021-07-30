package routes

import (
	"encoding/json"
	"net/http"

	iteamModel "github.com/DrIhor/test_task/internal/models/iteams"
	msgServ "github.com/DrIhor/test_task/internal/service/messages"

	"github.com/gorilla/mux"
)

// типу краще ініціалізовувати сервіси для БД чи передавати їх як параметр
func HandlerItems(router *mux.Router, stor iteamModel.IteamStorageServices) {
	router.HandleFunc("/iteams", func(w http.ResponseWriter, r *http.Request) {
		ShowAllIteams(w, r, stor)
	}).Methods("GET")
	router.HandleFunc("/iteams/{name}", func(w http.ResponseWriter, r *http.Request) {
		ShowIteam(w, r, stor)
	}).Methods("GET")
	router.HandleFunc("/iteams", func(w http.ResponseWriter, r *http.Request) {
		AddNewIteam(w, r, stor)
	}).Methods("POST")
	router.HandleFunc("/iteams/{name}", func(w http.ResponseWriter, r *http.Request) {
		BuyIteams(w, r, stor)
	}).Methods("PUT")
	router.HandleFunc("/iteams/{name}", func(w http.ResponseWriter, r *http.Request) {
		DeleteIteam(w, r, stor)
	}).Methods("DELETE")
}

// CRUD implementation
// Read
func ShowAllIteams(w http.ResponseWriter, r *http.Request, stor iteamModel.IteamStorageServices) {

	res, err := stor.GetAllIteams()
	if err != nil {
		res = msgServ.CreateMsgResp(err.Error())
		w.Write(res)
		return
	}

	w.Write(res)
}

func ShowIteam(w http.ResponseWriter, r *http.Request, stor iteamModel.IteamStorageServices) {
	params := mux.Vars(r)

	res, err := stor.GetIteam(params["name"])
	if err != nil {
		res = msgServ.CreateMsgResp(err.Error())
		w.Write(res)
		return
	}

	w.Write(res)
}

// Create
func AddNewIteam(w http.ResponseWriter, r *http.Request, stor iteamModel.IteamStorageServices) {

	var obj iteamModel.Iteam
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		res := msgServ.CreateMsgResp(err.Error())
		w.Write(res)
		return
	}

	if (obj == iteamModel.Iteam{}) {
		res := msgServ.CreateMsgResp("Empty body. Change information")
		w.Write(res)
		return
	}

	err = stor.AddNewIteam(obj)
	if err != nil {
		res := msgServ.CreateMsgResp("Empty body. Change information")
		w.Write(res)
		return
	}

}

// Update
func BuyIteams(w http.ResponseWriter, r *http.Request, stor iteamModel.IteamStorageServices) {
	params := mux.Vars(r)

	res, err := stor.UpdateIteam(params["name"])
	if err != nil {
		res := msgServ.CreateMsgResp("Empty body. Change information")
		w.Write(res)
		return
	}

	w.Write(res)
}

// Delete
func DeleteIteam(w http.ResponseWriter, r *http.Request, stor iteamModel.IteamStorageServices) {
	params := mux.Vars(r)

	err := stor.DeleteIteam(params["name"])
	if err != nil {
		res := msgServ.CreateMsgResp("Empty body. Change information")
		w.Write(res)
		return
	}
}
