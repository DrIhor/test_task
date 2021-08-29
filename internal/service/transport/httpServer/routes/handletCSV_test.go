package routes

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DrIhor/test_task/internal/models/items"
	mocks "github.com/DrIhor/test_task/mocks/intern/service/items"
	"github.com/gorilla/mux"
	"github.com/jszwec/csvutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// func TestFileUpload(t *testing.T) {
// 	// test info
// 	items := []items.Item{
// 		*defTestItem(),
// 		*defTestItem(),
// 		*defTestItem(),
// 	}

// 	type request struct {
// 		endpoint string
// 		method   string
// 		body     func() *bytes.Buffer
// 	}

// 	type wantResp struct {
// 		code       int
// 		contetType string
// 		body       func() []byte
// 	}

// 	testCases := []struct {
// 		name            string
// 		service         *mocks.ItemSrv // service to mock
// 		serviceFuncResp func(*mocks.ItemSrv, *bytes.Buffer)
// 		request         request
// 		want            wantResp
// 	}{
// 		{
// 			name:    "get all items(empty)",
// 			service: &mocks.ItemSrv{},
// 			serviceFuncResp: func(mc *mocks.ItemSrv, items *bytes.Buffer) {
// 				mc.On("AddFromCSV",
// 					mock.Anything,
// 				).Return(nil, nil)
// 			},
// 			request: request{
// 				endpoint: "/items/csv",
// 				method:   "POST",
// 				body: func() *bytes.Buffer {
// 					body := bytes.Buffer{}
// 					return &body
// 				},
// 			},
// 			want: wantResp{
// 				code: 400,
// 				body: func() []byte {
// 					res, _ := csvutil.Marshal(nil)
// 					return res
// 				},
// 				contetType: "",
// 			},
// 		},
// 		{
// 			name:    "get all items(is data)",
// 			service: &mocks.ItemSrv{},
// 			serviceFuncResp: func(mc *mocks.ItemSrv, items *bytes.Buffer) {
// 				mc.On("AddFromCSV",
// 					mock.Anything,
// 				).Return(func(_ context.Context) *bytes.Buffer {
// 					return items
// 				}, nil)
// 			},
// 			request: request{
// 				endpoint: "/items/csv",
// 				method:   "POST",
// 				body: func() *bytes.Buffer {
// 					body := bytes.Buffer{}

// 					info := items
// 					res, err := csvutil.Marshal(info)
// 					if err != nil {
// 						log.Fatal("Input data marshal: ", err)
// 					}

// 					writer := multipart.NewWriter(&body)
// 					w, _ := writer.CreateFormField("file")
// 					io.Copy(w, bytes.NewReader(res))

// 					return &body
// 				},
// 			},
// 			want: wantResp{
// 				code: 200,
// 				body: func() []byte {
// 					info := items

// 					res, err := csvutil.Marshal(info)
// 					if err != nil {
// 						log.Fatal("Input data marshal: ", err)
// 					}

// 					return res
// 				},
// 				contetType: "",
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {

// 		// init handler
// 		handler := HandlerItemsServ{
// 			router:   mux.NewRouter(),
// 			ctx:      context.Background(),
// 			services: tc.service,
// 		}
// 		body := tc.request.body()
// 		tc.serviceFuncResp(tc.service, body) // mock internal function
// 		handler.HandlerItems()

// 		// server
// 		hts := httptest.NewServer(handler.router)
// 		defer hts.Close()

// 		// client
// 		cl := hts.Client()
// 		req, _ := http.NewRequest(tc.request.method, hts.URL+tc.request.endpoint, body)

// 		// request proces
// 		res, err := cl.Do(req)
// 		if err != nil {
// 			t.Error("request error :", err)
// 		}

// 		contentType := res.Header.Get("Content-Type")

// 		// results
// 		assert.Equal(t, tc.want.contetType, contentType, tc.name)
// 		assert.Equal(t, tc.want.code, res.StatusCode, tc.name)
// 	}
// }

func TestFileDownload(t *testing.T) {
	// test info
	items := []items.Item{
		*defTestItem(),
		*defTestItem(),
		*defTestItem(),
	}

	type request struct {
		endpoint string
		method   string
		body     func() []byte
	}

	type wantResp struct {
		code       int
		contetType string
		body       func() []byte
	}

	testCases := []struct {
		name            string
		service         *mocks.ItemSrv // service to mock
		serviceFuncResp func(*mocks.ItemSrv, []byte)
		request         request
		want            wantResp
	}{
		{
			name:    "get all items(empty)",
			service: &mocks.ItemSrv{},
			serviceFuncResp: func(mc *mocks.ItemSrv, items []byte) {
				mc.On("GetAllItemsAsCSV",
					mock.Anything,
				).Return(nil, nil)
			},
			request: request{
				endpoint: "/items/csv",
				method:   "GET",
				body: func() []byte {
					res, _ := csvutil.Marshal(nil)
					return res
				},
			},
			want: wantResp{
				code: 404,
				body: func() []byte {
					res, _ := csvutil.Marshal(nil)
					return res
				},
				contetType: "",
			},
		},
		{
			name:    "get all items(is data)",
			service: &mocks.ItemSrv{},
			serviceFuncResp: func(mc *mocks.ItemSrv, items []byte) {
				mc.On("GetAllItemsAsCSV",
					mock.Anything,
				).Return(func(_ context.Context) []byte {
					return items
				}, nil)
			},
			request: request{
				endpoint: "/items/csv",
				method:   "GET",
				body: func() []byte {
					info := items

					res, err := csvutil.Marshal(info)
					if err != nil {
						log.Fatal("Input data marshal: ", err)
					}

					return res
				},
			},
			want: wantResp{
				code: 200,
				body: func() []byte {
					info := items

					res, err := csvutil.Marshal(info)
					if err != nil {
						log.Fatal("Input data marshal: ", err)
					}

					return res
				},
				contetType: "text/csv",
			},
		},
	}

	for _, tc := range testCases {

		// init handler
		handler := HandlerItemsServ{
			router:   mux.NewRouter(),
			ctx:      context.Background(),
			services: tc.service,
		}
		tc.serviceFuncResp(tc.service, tc.request.body()) // mock internal function
		handler.HandlerItems()

		// server
		hts := httptest.NewServer(handler.router)
		defer hts.Close()

		// client
		cl := hts.Client()
		req, _ := http.NewRequest(tc.request.method, hts.URL+tc.request.endpoint, bytes.NewReader(tc.request.body()))

		// request proces
		res, err := cl.Do(req)
		if err != nil {
			t.Error("request error :", err)
		}

		contentType := res.Header.Get("Content-Type")

		// results
		assert.Equal(t, tc.want.contetType, contentType, tc.name)
		assert.Equal(t, tc.want.code, res.StatusCode, tc.name)
	}
}
