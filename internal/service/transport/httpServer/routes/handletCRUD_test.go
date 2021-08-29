package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DrIhor/test_task/internal/models/items"
	itemModel "github.com/DrIhor/test_task/internal/models/items"
	mocks "github.com/DrIhor/test_task/mocks/intern/service/items"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func defTestItem() *items.Item {
	return &items.Item{
		ID:          uuid.New().String(),
		Name:        "testItem",
		Price:       int32(rand.Int()),
		ItemsNumber: int32(rand.Int()),
		Description: "Some descr",
	}
}

func TestHTTPHandlerAllItems(t *testing.T) {
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
		code int
		body func() []byte
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
				mc.On("GetAllItems",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
				).Return(nil, nil)
			},
			request: request{
				endpoint: "/items",
				method:   "GET",
				body: func() []byte {
					res, _ := json.Marshal(nil)
					return res
				},
			},
			want: wantResp{
				code: 404,
				body: func() []byte {
					res, _ := json.Marshal(nil)
					return res
				},
			},
		},
		{
			name:    "get all items(is data)",
			service: &mocks.ItemSrv{},
			serviceFuncResp: func(mc *mocks.ItemSrv, items []byte) {
				mc.On("GetAllItems",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
				).Return(func(_ context.Context) []byte {
					return items
				}, nil)
			},
			request: request{
				endpoint: "/items",
				method:   "GET",
				body: func() []byte {
					info := items

					res, err := json.Marshal(info)
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

					res, err := json.Marshal(info)
					if err != nil {
						log.Fatal("Input data marshal: ", err)
					}

					return res
				},
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

		// results
		assert.Equal(t, tc.want.code, res.StatusCode, tc.name)
	}
}

func TestHTTPHandleGetItem(t *testing.T) {
	items := []items.Item{
		*defTestItem(),
		*defTestItem(),
		*defTestItem(),
	}

	type request struct {
		endpoint string
		method   string
		paramID  func() string
		body     func() []byte
	}

	type wantResp struct {
		code int
		body func() []byte
	}

	// all cases return second test item
	testCases := []struct {
		name            string
		service         *mocks.ItemSrv // service to mock
		serviceFuncResp func(*mocks.ItemSrv, []byte)
		request         request
		want            wantResp
	}{
		{
			name:    "get item(empty data to search)",
			service: &mocks.ItemSrv{},
			serviceFuncResp: func(mc *mocks.ItemSrv, items []byte) {
				mc.On("GetItem",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					mock.MatchedBy(func(id string) bool { return true }),
				).Return(nil, nil)
			},
			request: request{
				endpoint: "/item",
				method:   "GET",
				paramID: func() string {
					return uuid.New().String()
				},
				body: func() []byte {
					res, _ := json.Marshal(itemModel.Item{})
					return res
				},
			},
			want: wantResp{
				code: 404,
				body: func() []byte {
					res, _ := json.Marshal(itemModel.Item{})
					return res
				},
			},
		},
		{
			name:    "get item(wrong query)",
			service: &mocks.ItemSrv{},
			serviceFuncResp: func(mc *mocks.ItemSrv, items []byte) {
				mc.On("GetItem",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					mock.MatchedBy(func(id string) bool { return true }),
				).Return(nil, nil)
			},
			request: request{
				endpoint: "/item",
				method:   "GET",
				paramID: func() string {
					return "12345"
				},
				body: func() []byte {
					res, _ := json.Marshal(itemModel.Item{})
					return res
				},
			},
			want: wantResp{
				code: 400,
				body: func() []byte {
					res, _ := json.Marshal(itemModel.Item{})

					return res
				},
			},
		},
		{
			name:    "get item(not exist from exist rows)",
			service: &mocks.ItemSrv{},
			serviceFuncResp: func(mc *mocks.ItemSrv, items []byte) {
				mc.On("GetItem",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					mock.MatchedBy(func(id string) bool { return true }),
				).Return(nil, nil)
			},
			request: request{
				endpoint: "/item",
				method:   "GET",
				paramID: func() string {
					return uuid.New().String()
				},
				body: func() []byte {
					res, _ := json.Marshal(itemModel.Item{})
					return res
				},
			},
			want: wantResp{
				code: 404,
				body: func() []byte {
					res, _ := json.Marshal(itemModel.Item{})
					return res
				},
			},
		},
		{
			name:    "get item(is data)",
			service: &mocks.ItemSrv{},
			serviceFuncResp: func(mc *mocks.ItemSrv, item []byte) {
				mc.On("GetItem",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					mock.MatchedBy(func(id string) bool { return true }),
				).Return(func(_ context.Context, id string) []byte {
					return item
				}, nil)
			},
			request: request{
				endpoint: "/item",
				method:   "GET",
				paramID: func() string {
					return items[1].ID
				},
				body: func() []byte {
					info := items[1]

					res, err := json.Marshal(info)
					if err != nil {
						log.Fatal("Input data marshal: ", err)
					}

					return res
				},
			},
			want: wantResp{
				code: 200,
				body: func() []byte {
					info := items[1]

					res, err := json.Marshal(info)
					if err != nil {
						log.Fatal("Input data marshal: ", err)
					}

					return res
				},
			},
		},
	}

	for _, tc := range testCases {

		// init handler
		handler := New(context.Background(), mux.NewRouter(), tc.service)
		tc.serviceFuncResp(tc.service, tc.request.body()) // mock internal function
		handler.HandlerItems()

		// server
		hts := httptest.NewServer(handler.router)
		defer hts.Close()

		// client
		cl := hts.Client()
		req, _ := http.NewRequest(tc.request.method, hts.URL+tc.request.endpoint, bytes.NewReader(tc.request.body()))
		q := req.URL.Query()
		q.Add("id", tc.request.paramID())
		req.URL.RawQuery = q.Encode()

		// request proces
		res, err := cl.Do(req)
		if err != nil {
			t.Error("request error :", err)
		}

		var obj itemModel.Item
		json.NewDecoder(res.Body).Decode(&obj)

		result, err := json.Marshal(obj)
		if err != nil {
			log.Fatal(err)
		}

		// results
		assert.Equal(t, tc.want.code, res.StatusCode, tc.name)
		assert.Equal(t, tc.want.body(), result, tc.name)
	}
}

func TestAddItem(t *testing.T) {
	type request struct {
		endpoint string
		method   string
		body     func() []byte
	}

	type wantResp struct {
		code   int
		bodyID func(string) string
	}

	// all cases return second test item
	testCases := []struct {
		name            string
		service         *mocks.ItemSrv // service to mock
		serviceFuncResp func(*mocks.ItemSrv)
		request         request
		want            wantResp
	}{
		{
			name:    "add(normal data)",
			service: &mocks.ItemSrv{},
			serviceFuncResp: func(mc *mocks.ItemSrv) {
				mc.On("AddNewItem",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					mock.MatchedBy(func(_ itemModel.Item) bool { return true }),
				).Return(uuid.New().String(), nil)
			},
			request: request{
				endpoint: "/item",
				method:   "POST",
				body: func() []byte {
					itemData := *defTestItem()
					itemData.ID = ""

					res, err := json.Marshal(itemData)
					if err != nil {
						log.Fatal("Input data marshal: ", err)
					}

					return res
				},
			},
			want: wantResp{
				code: 201,
				bodyID: func(id string) string {
					return id
				},
			},
		},
		{
			name:    "add(empty item)",
			service: &mocks.ItemSrv{},
			serviceFuncResp: func(mc *mocks.ItemSrv) {
				mc.On("AddNewItem",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					mock.MatchedBy(func(_ itemModel.Item) bool { return true }),
				).Return(nil, nil)
			},

			request: request{
				endpoint: "/item",
				method:   "POST",
				body: func() []byte {
					res, _ := json.Marshal(itemModel.Item{})
					return res
				},
			},
			want: wantResp{
				code: 400,
				bodyID: func(_ string) string {
					return ""
				},
			},
		},
	}

	for _, tc := range testCases {

		// init handler
		handler := New(context.Background(), mux.NewRouter(), tc.service)

		tc.serviceFuncResp(tc.service) // mock internal function
		handler.HandlerItems()

		// server
		hts := httptest.NewServer(handler.router)
		defer hts.Close()

		// client
		cl := hts.Client()
		req, err := http.NewRequest(tc.request.method, hts.URL+tc.request.endpoint, bytes.NewReader(tc.request.body()))
		if err != nil {
			log.Fatal("create request ", err)
		}
		// request proces
		res, err := cl.Do(req)
		if err != nil {
			t.Error("request error :", err)
		}

		var obj itemModel.Item
		if err := json.NewDecoder(res.Body).Decode(&obj); err != nil {
			log.Fatal("res decode ", err)
		}

		// results
		assert.Equal(t, tc.want.code, res.StatusCode, tc.name)
		assert.Equal(t, tc.want.bodyID(obj.ID), obj.ID, tc.name)
	}
}

func TestHTTPHandleBuyItem(t *testing.T) {
	items := []items.Item{
		itemModel.Item{
			ID:          uuid.New().String(),
			Name:        "test1",
			Price:       10,
			ItemsNumber: 10,
			Description: "Some desc",
		},
		itemModel.Item{
			ID:          uuid.New().String(),
			Name:        "test2",
			Price:       20,
			ItemsNumber: 3,
			Description: "Some desc",
		},
		itemModel.Item{
			ID:          uuid.New().String(),
			Name:        "test2",
			Price:       30,
			ItemsNumber: 1,
			Description: "Some desc",
		},
	}

	type request struct {
		endpoint  string
		method    string
		paramItem func() itemModel.Item
		body      func() []byte

		requestRepeats int32
		requestData    []itemModel.Item
	}

	type wantResp struct {
		code int
		body func(itemModel.Item) []byte
	}

	// all cases return second test item
	testCases := []struct {
		name            string
		service         *mocks.ItemSrv // service to mock
		serviceFuncResp func(*mocks.ItemSrv, itemModel.Item)
		request         request
		want            wantResp
	}{
		{
			name:    "normal update",
			service: &mocks.ItemSrv{},
			serviceFuncResp: func(mc *mocks.ItemSrv, item itemModel.Item) {
				mc.On("UpdateItem",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					mock.MatchedBy(func(id string) bool { return true }),
				).Return(func(_ context.Context, id string) []byte {
					for i := range items {
						if items[i].ID == id {
							items[i].ItemsNumber--
							if items[i].ItemsNumber <= 0 {
								break
							}

							res, _ := json.Marshal(items[i])
							return res
						}
					}

					res, _ := json.Marshal(nil)
					return res
				}, nil)
			},
			request: request{
				endpoint: "/item",
				method:   "PUT",
				paramItem: func() itemModel.Item {
					return items[0]
				},
				requestRepeats: 1,
				requestData:    items,
				body: func() []byte {
					return []byte{}
				},
			},
			want: wantResp{
				code: 200,
				body: func(item itemModel.Item) []byte {
					item.ItemsNumber--
					res, _ := json.Marshal(item)
					return res
				},
			},
		},
		{
			name:    "delete value",
			service: &mocks.ItemSrv{},
			serviceFuncResp: func(mc *mocks.ItemSrv, item itemModel.Item) {
				mc.On("UpdateItem",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					mock.MatchedBy(func(id string) bool { return true }),
				).Return(func(_ context.Context, id string) []byte {
					if item.ID == id {
						item.ItemsNumber--
						if item.ItemsNumber <= 0 {
							return nil
						}

						res, _ := json.Marshal(item)
						return res
					}

					return nil
				}, nil)
			},
			request: request{
				endpoint: "/item",
				method:   "PUT",
				paramItem: func() itemModel.Item {
					return items[2]
				},
				requestRepeats: 1,
				requestData:    items,
				body: func() []byte {
					return []byte{}
				},
			},
			want: wantResp{
				code: 404,
				body: func(item itemModel.Item) []byte {
					res, _ := json.Marshal(itemModel.Item{})
					return res
				},
			},
		},
		{
			name:    "wrong delete value",
			service: &mocks.ItemSrv{},
			serviceFuncResp: func(mc *mocks.ItemSrv, item itemModel.Item) {
				mc.On("UpdateItem",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					mock.MatchedBy(func(id string) bool { return true }),
				).Return(func(_ context.Context, id string) []byte {
					if item.ID == id {
						item.ItemsNumber--
						if item.ItemsNumber <= 0 {
							return nil
						}

						res, _ := json.Marshal(item)
						return res
					}

					return nil
				}, nil)
			},
			request: request{
				endpoint: "/item",
				method:   "PUT",
				paramItem: func() itemModel.Item {
					return itemModel.Item{}
				},
				requestRepeats: 1,
				requestData:    items,
				body: func() []byte {
					return []byte{}
				},
			},
			want: wantResp{
				code: 400,
				body: func(item itemModel.Item) []byte {
					res, _ := json.Marshal(itemModel.Item{})
					return res
				},
			},
		},
	}

	for _, tc := range testCases {

		// init handler
		handler := New(context.Background(), mux.NewRouter(), tc.service)

		data := tc.request.paramItem()
		tc.serviceFuncResp(tc.service, data) // mock internal function

		handler.HandlerItems()

		// server
		hts := httptest.NewServer(handler.router)
		defer hts.Close()

		// client
		cl := hts.Client()
		req, _ := http.NewRequest(tc.request.method, hts.URL+tc.request.endpoint, bytes.NewReader(tc.request.body()))
		item := tc.request.paramItem()
		q := req.URL.Query()
		q.Add("id", item.ID)
		req.URL.RawQuery = q.Encode()

		// request proces
		res, err := cl.Do(req)
		if err != nil {
			t.Error("request error :", err)
		}

		var obj itemModel.Item
		json.NewDecoder(res.Body).Decode(&obj)

		result, err := json.Marshal(obj)
		if err != nil {
			log.Fatal(err)
		}

		// results
		assert.Equal(t, tc.want.code, res.StatusCode, tc.name)
		assert.Equal(t, tc.want.body(item), result, tc.name)
	}
}

func TestHTTPHandleDeleteItem(t *testing.T) {
	items := []items.Item{
		itemModel.Item{
			ID:          uuid.New().String(),
			Name:        "test1",
			Price:       10,
			ItemsNumber: 10,
			Description: "Some desc",
		},
		itemModel.Item{
			ID:          uuid.New().String(),
			Name:        "test2",
			Price:       20,
			ItemsNumber: 3,
			Description: "Some desc",
		},
		itemModel.Item{
			ID:          uuid.New().String(),
			Name:        "test2",
			Price:       30,
			ItemsNumber: 1,
			Description: "Some desc",
		},
	}

	type request struct {
		endpoint    string
		method      string
		paramItemID func() string
		body        func() []byte

		requestData []itemModel.Item
	}

	type wantResp struct {
		code int
	}

	// all cases return second test item
	testCases := []struct {
		name            string
		service         *mocks.ItemSrv // service to mock
		serviceFuncResp func(*mocks.ItemSrv, string)
		request         request
		want            wantResp
	}{
		{
			name:    "normal delete",
			service: &mocks.ItemSrv{},
			serviceFuncResp: func(mc *mocks.ItemSrv, id string) {
				mc.On("DeleteItem",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					mock.MatchedBy(func(id string) bool { return true }),
				).Return(func(_ context.Context, id string) bool {
					for i := range items {
						if items[i].ID == id {
							return true
						}
					}

					return false
				}, nil)
			},
			request: request{
				endpoint: "/item",
				method:   "DELETE",
				paramItemID: func() string {
					return items[0].ID
				},
				requestData: items,
				body: func() []byte {
					return []byte{}
				},
			},
			want: wantResp{
				code: 200,
			},
		},
		{
			name:    "delete not exist value",
			service: &mocks.ItemSrv{},
			serviceFuncResp: func(mc *mocks.ItemSrv, id string) {
				mc.On("DeleteItem",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					mock.MatchedBy(func(id string) bool { return true }),
				).Return(func(_ context.Context, id string) bool {
					for i := range items {
						if items[i].ID == id {
							return true
						}
					}

					return false
				}, nil)
			},
			request: request{
				endpoint: "/item",
				method:   "DELETE",
				paramItemID: func() string {
					return uuid.New().String()
				},
				requestData: items,
				body: func() []byte {
					return []byte{}
				},
			},
			want: wantResp{
				code: 404,
			},
		},
		{
			name:    "wrong delete value",
			service: &mocks.ItemSrv{},
			serviceFuncResp: func(mc *mocks.ItemSrv, id string) {
				mc.On("DeleteItem",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					mock.MatchedBy(func(id string) bool { return true }),
				).Return(func(_ context.Context, id string) bool {
					for i := range items {
						if items[i].ID == id {
							return true
						}
					}

					return false
				}, nil)
			},
			request: request{
				endpoint: "/item",
				method:   "DELETE",
				paramItemID: func() string {
					return ""
				},
				requestData: items,
				body: func() []byte {
					return []byte{}
				},
			},
			want: wantResp{
				code: 400,
			},
		},
	}

	for _, tc := range testCases {

		// init handler
		handler := New(context.Background(), mux.NewRouter(), tc.service)
		tc.serviceFuncResp(tc.service, tc.request.paramItemID()) // mock internal function
		handler.HandlerItems()

		// server
		hts := httptest.NewServer(handler.router)
		defer hts.Close()

		// client
		cl := hts.Client()
		req, _ := http.NewRequest(tc.request.method, hts.URL+tc.request.endpoint, bytes.NewReader(tc.request.body()))
		q := req.URL.Query()
		q.Add("id", tc.request.paramItemID())
		req.URL.RawQuery = q.Encode()

		// request proces
		res, err := cl.Do(req)
		if err != nil {
			t.Error("request error :", err)
		}

		// results
		assert.Equal(t, tc.want.code, res.StatusCode, tc.name)
	}
}
