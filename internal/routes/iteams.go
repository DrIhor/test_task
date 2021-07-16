package routes

import (
	"encoding/json"
	"net/http"

	mem "github.com/DrIhor/test_task/internal/storage/memory"
	"github.com/gorilla/mux"
)

// CRUD implementation

// Read
func showAllIteams(w http.ResponseWriter, r *http.Request) {

	res, err := mem.ShowAllIteams()
	if err != nil {
		res = createMsgResp(err.Error())
		w.Write(res)
		return
	}

	w.Write(res)
}

func showIteam(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	res, err := mem.ShowIteam(params["name"])
	if err != nil {
		res = createMsgResp(err.Error())
		w.Write(res)
		return
	}

	w.Write(res)
}

// Create
func addNewIteam(w http.ResponseWriter, r *http.Request) {
	var obj mem.Iteam
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		res := createMsgResp(err.Error())
		w.Write(res)
		return
	}

	if (obj == mem.Iteam{}) {
		res := createMsgResp("Empty body. Change information")
		w.Write(res)
		return
	}

	err = mem.AddNewIteam(obj)
	if err != nil {
		res := createMsgResp("Empty body. Change information")
		w.Write(res)
		return
	}

}

// Update
func buyIteams(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	res, err := mem.UpdateIteam(params["name"])
	if err != nil {
		res := createMsgResp("Empty body. Change information")
		w.Write(res)
		return
	}

	w.Write(res)
}

// Delete
func deleteIteam(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := mem.DeleteIteam(params["name"])
	if err != nil {
		res := createMsgResp("Empty body. Change information")
		w.Write(res)
		return
	}
}
