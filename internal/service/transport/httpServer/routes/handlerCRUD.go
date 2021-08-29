package routes

import (
	"encoding/json"
	"net/http"
	"os"

	itemModel "github.com/DrIhor/test_task/internal/models/items"
	"github.com/DrIhor/test_task/internal/service/connectors"
	msgServ "github.com/DrIhor/test_task/internal/service/messages"
	"github.com/google/uuid"
)

// CRUD implementation for all endpoints
// Read
func (h *HandlerItemsServ) ShowAllItems(w http.ResponseWriter, r *http.Request) {

	var errData error
	var res []byte
	switch os.Getenv("STORAGE_TYPE") {
	case "":
		res, errData = h.services.GetAllItems(h.ctx)
	case "grpc":
		grpcConn := connectors.NewGRPC(h.ctx, os.Getenv("GRCP_ADDR"))
		res, errData = grpcConn.GetAllItems(h.ctx)
	}

	if errData != nil {
		res = msgServ.CreateMsgResp(errData.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	// check if struct is empty
	if len(res) < 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Write(res)
}

func (h *HandlerItemsServ) ShowItem(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := keys[0]

	_, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var errData error
	var res []byte
	switch os.Getenv("STORAGE_TYPE") {
	case "":
		res, errData = h.services.GetItem(h.ctx, id)
	case "grpc":
		grpcConn := connectors.NewGRPC(h.ctx, os.Getenv("GRCP_ADDR"))
		res, errData = grpcConn.GetItem(h.ctx, id)
	}

	if errData != nil {
		res = msgServ.CreateMsgResp(errData.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	// check if struct is empty
	if len(res) < 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// Create
func (h *HandlerItemsServ) AddNewItem(w http.ResponseWriter, r *http.Request) {

	var obj itemModel.Item
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil || (obj == itemModel.Item{}) {
		res := msgServ.CreateMsgResp("Empty body. Change information")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	var id string
	var errData error
	switch os.Getenv("STORAGE_TYPE") {
	case "":
		id, errData = h.services.AddNewItem(h.ctx, obj)
	case "grpc":
		grpcConn := connectors.NewGRPC(h.ctx, os.Getenv("GRCP_ADDR"))
		id, errData = grpcConn.AddNewItem(h.ctx, obj)
	}

	if errData != nil {
		res := msgServ.CreateMsgResp(errData.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	item := itemModel.Item{
		ID: id,
	}

	res, err := json.Marshal(item)
	if err != nil {
		res := msgServ.CreateMsgResp(errData.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

// Update
func (h *HandlerItemsServ) BuyItems(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := keys[0]
	_, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var errData error
	var res []byte
	switch os.Getenv("STORAGE_TYPE") {
	case "":
		res, errData = h.services.UpdateItem(h.ctx, id)
	case "grpc":
		grpcConn := connectors.NewGRPC(h.ctx, os.Getenv("GRCP_ADDR"))
		res, errData = grpcConn.UpdateItem(h.ctx, id)
	}

	// check if struct is empty
	if len(res) < 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if errData != nil {
		res := msgServ.CreateMsgResp(errData.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	w.Write(res)
}

// Delete
func (h *HandlerItemsServ) DeleteItem(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := keys[0]
	_, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var done bool
	var errData error
	switch os.Getenv("STORAGE_TYPE") {
	case "":
		done, errData = h.services.DeleteItem(h.ctx, id)
	case "grpc":
		grpcConn := connectors.NewGRPC(h.ctx, os.Getenv("GRCP_ADDR"))
		done, errData = grpcConn.DeleteItem(h.ctx, id)
	}

	if errData != nil {
		res := msgServ.CreateMsgResp(errData.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
		return
	}

	if !done {
		w.WriteHeader(http.StatusNotFound)
	}
}
