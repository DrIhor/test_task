package routes

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"
)

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
		res, err = h.services.AddFromCSV(h.ctx, rd)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
}

func (h *HandlerItemsServ) ReturnAllItemsCSV(w http.ResponseWriter, r *http.Request) {

	res, err := h.services.GetAllItemsAsCSV(h.ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(res) < 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// add some headers
	w.Header().Set("Content-Disposition", "multipart/form-data; boundary=something;")
	w.Header().Set("Content-Type", "text/csv")

	io.Copy(w, bytes.NewReader(res))
	return
}
