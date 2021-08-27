package routes

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	itemModel "github.com/DrIhor/test_task/internal/models/items"
	"github.com/DrIhor/test_task/internal/service/connectors"
	itemServ "github.com/DrIhor/test_task/internal/service/items"
	msgServ "github.com/DrIhor/test_task/internal/service/messages"
	"github.com/google/uuid"

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
		services: itemServ.New(ctx, stor),
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
	h.router.HandleFunc("/items/csv", h.AddDataFromCSV).Methods("Post")
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
		grpcConn := connectors.NewGRPC(h.ctx, os.Getenv("GRCP_ADDR"))
		res, errData = grpcConn.GetAllItems(h.ctx)
	}

	if errData != nil {
		res = msgServ.CreateMsgResp(errData.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(res)
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
		res, errData = h.services.GetItem(id)
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
		id, errData = h.services.AddNewItem(obj)
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
		res, errData = h.services.UpdateItem(id)
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
		done, errData = h.services.DeleteItem(id)
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

// file work
func (h *HandlerItemsServ) AddDataFromCSV(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(256) // limit your max input length!
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024) // max is 5 MB

	// upload file
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// log some data
	name := strings.Split(fileHeader.Filename, ".")
	fmt.Printf("File name %s\n", name[0])
	fmt.Printf("File type %v\n", fileHeader.Header["Content-Type"][0])
	fmt.Printf("File size %v\n", fileHeader.Size)

	contentType := fileHeader.Header["Content-Type"][0]

	var res []byte
	switch contentType {
	case "text/csv":
		rd := csv.NewReader(file)
		res, err = h.services.AddFromCSV(rd)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(err)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
}

func (h *HandlerItemsServ) ReturnAllItemsCSV(w http.ResponseWriter, r *http.Request) {

	res, err := h.services.GetAllItemsAsCSV()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// add some headers
	w.Header().Set("Content-Disposition", "multipart/form-data; boundary=something;")
	w.Header().Set("Content-Type", "text/csv")

	io.Copy(w, bytes.NewReader(res))
	return
}
