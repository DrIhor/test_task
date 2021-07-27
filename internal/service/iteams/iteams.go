package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/DrIhor/test_task/internal/models/iteams"
	iteamModel "github.com/DrIhor/test_task/internal/models/iteams"
	mess "github.com/DrIhor/test_task/internal/service/messages"
)

// CRUD implementation
// Read
func ShowAllIteams(w http.ResponseWriter, r *http.Request, stor iteams.IteamServices) {

	res, err := stor.GetAllIteams()
	if err != nil {
		res = mess.CreateMsgResp(err.Error())
		w.Write(res)
		return
	}

	w.Write(res)
}

func ShowIteam(w http.ResponseWriter, r *http.Request, stor iteams.IteamServices) {
	params := mux.Vars(r)

	res, err := stor.GetIteam(params["name"])
	if err != nil {
		res = mess.CreateMsgResp(err.Error())
		w.Write(res)
		return
	}

	w.Write(res)
}

// Create
func AddNewIteam(w http.ResponseWriter, r *http.Request, stor iteams.IteamServices) {

	var obj iteamModel.Iteam
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		res := mess.CreateMsgResp(err.Error())
		w.Write(res)
		return
	}

	if (obj == iteamModel.Iteam{}) {
		res := mess.CreateMsgResp("Empty body. Change information")
		w.Write(res)
		return
	}

	err = stor.AddNewIteam(obj)
	if err != nil {
		res := mess.CreateMsgResp("Empty body. Change information")
		w.Write(res)
		return
	}

}

// Update
func BuyIteams(w http.ResponseWriter, r *http.Request, stor iteams.IteamServices) {
	params := mux.Vars(r)

	res, err := stor.UpdateIteam(params["name"])
	if err != nil {
		res := mess.CreateMsgResp("Empty body. Change information")
		w.Write(res)
		return
	}

	w.Write(res)
}

// Delete
func DeleteIteam(w http.ResponseWriter, r *http.Request, stor iteams.IteamServices) {
	params := mux.Vars(r)

	err := stor.DeleteIteam(params["name"])
	if err != nil {
		res := mess.CreateMsgResp("Empty body. Change information")
		w.Write(res)
		return
	}
}
