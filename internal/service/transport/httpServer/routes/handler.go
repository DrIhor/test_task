package routes

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	itemModel "github.com/DrIhor/test_task/internal/models/items"
	"github.com/DrIhor/test_task/internal/service/connectors"
	itemServ "github.com/DrIhor/test_task/internal/service/items"
	msgServ "github.com/DrIhor/test_task/internal/service/messages"

	"github.com/gorilla/mux"
	"github.com/jszwec/csvutil"
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
		grpcConn := connectors.NewGRPC(os.Getenv("GRCP_ADDR"))
		res, errData = grpcConn.GetAllItems()
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

	id, err := strconv.Atoi(keys[0])
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
		grpcConn := connectors.NewGRPC(os.Getenv("GRCP_ADDR"))
		res, errData = grpcConn.GetItem(id)
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

	var id int
	var errData error
	switch os.Getenv("STORAGE_TYPE") {
	case "":
		id, errData = h.services.AddNewItem(obj)
	case "grpc":
		grpcConn := connectors.NewGRPC(os.Getenv("GRCP_ADDR"))
		id, errData = grpcConn.AddNewItem(obj)
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

	id, err := strconv.Atoi(keys[0])
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
		grpcConn := connectors.NewGRPC(os.Getenv("GRCP_ADDR"))
		res, errData = grpcConn.UpdateItem(id)
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
	id, err := strconv.Atoi(keys[0])
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
		grpcConn := connectors.NewGRPC(os.Getenv("GRCP_ADDR"))
		done, errData = grpcConn.DeleteItem(id)
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
	fmt.Printf("File type %v\n", fileHeader.Header)
	fmt.Printf("File size %v\n", fileHeader.Size)

	//
	// start read file data
	//
	rd := csv.NewReader(file)
	var (
		itemHeader string        // struct fields
		firstRow   bool   = true // if file header
		newIDs     []int         // result data
	)

	for {
		row, err := rd.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// init header
		if firstRow {
			itemHeader = strings.Join(row, ",") + "\n" // create single row for csvutil lib
			firstRow = false
			continue
		}

		csvIteam := []byte(itemHeader + strings.Join(row, ","))
		var items []itemModel.Item
		if err := csvutil.Unmarshal(csvIteam, &items); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, item := range items {
			if item != (itemModel.Item{}) {
				id, err := h.services.AddNewItem(item)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				newIDs = append(newIDs, id)
			}
		}
	}

	res, err := json.Marshal(newIDs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
}

func (h *HandlerItemsServ) ReturnAllItemsCSV(w http.ResponseWriter, r *http.Request) {
	byteData, err := h.services.GetAllItems()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var itemsSlice []itemModel.Item
	if err := json.Unmarshal(byteData, &itemsSlice); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := csvutil.Marshal(itemsSlice)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "multipart/form-data; boundary=something;")
	w.Header().Set("Content-Type", "text/csv")

	io.Copy(w, bytes.NewReader(b))
	return
}