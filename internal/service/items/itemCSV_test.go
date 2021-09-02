package items

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"log"
	"math/rand"
	"testing"

	er "github.com/DrIhor/test_task/internal/errors"
	"github.com/DrIhor/test_task/internal/models/items"
	mocks "github.com/DrIhor/test_task/mocks/intern/service/items"
	"github.com/google/uuid"
	"github.com/jszwec/csvutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func defTestItem() *items.Item {
	return &items.Item{
		Name:        "testItem",
		Price:       int32(rand.Int()),
		ItemsNumber: int32(rand.Int()),
		Description: "Some descr",
	}
}

func TestAddFromCSV(t *testing.T) {
	type request struct {
		ctx  context.Context
		data func([]items.Item) *csv.Reader

		defDataRequest items.Item
	}

	type wantResp struct {
		data func(items.Item) []byte
		err  error

		isErr bool
	}

	testCases := []struct {
		name            string
		serviceMock     *mocks.ItemSrv // service to mock
		serviceFuncResp func(*mocks.ItemSrv, string)
		request         request
		want            wantResp
	}{
		{
			name:        "read csv normal data",
			serviceMock: &mocks.ItemSrv{},
			serviceFuncResp: func(mc *mocks.ItemSrv, id string) {
				mc.On("AddNewItem", mock.Anything, mock.Anything).Return(func(_ context.Context, _ items.Item) string {
					return id
				}, nil)
			},
			request: request{ctx: context.Background(), data: func(item []items.Item) *csv.Reader {
				byteData, err := csvutil.Marshal(item)
				if err != nil {
					log.Fatal(err)
				}
				r := bytes.NewReader(byteData)

				return csv.NewReader(r)
			},
				defDataRequest: items.Item{
					ID:          uuid.New().String(),
					Name:        "test",
					Price:       int32(rand.Int()),
					ItemsNumber: int32(rand.Int()),
					Description: "Some descr",
				},
			},
			want: wantResp{
				data: func(item items.Item) []byte {
					res := []string{item.ID}

					data, _ := json.Marshal(res)
					return data
				},
				err:   nil,
				isErr: false,
			},
		},
		{
			name:        "return test error",
			serviceMock: &mocks.ItemSrv{},
			serviceFuncResp: func(mc *mocks.ItemSrv, id string) {
				mc.On("AddNewItem", mock.Anything, mock.Anything).Return(func(_ context.Context, _ items.Item) string {
					return ""
				}, er.DataNotExist)
			},
			request: request{ctx: context.Background(), data: func(item []items.Item) *csv.Reader {
				byteData, err := csvutil.Marshal(item)
				if err != nil {
					log.Fatal(err)
				}
				r := bytes.NewReader(byteData)

				return csv.NewReader(r)
			},
				defDataRequest: items.Item{
					ID:          uuid.New().String(),
					Name:        "test",
					Price:       int32(rand.Int()),
					ItemsNumber: int32(rand.Int()),
					Description: "Some descr",
				},
			},
			want: wantResp{
				data: func(_ items.Item) []byte {
					return nil
				},
				err:   er.DataNotExist,
				isErr: true,
			},
		},
	}
	for _, tc := range testCases {
		mockStor := tc.serviceMock
		itemServ := New(mockStor)

		itemData := tc.request.defDataRequest
		tc.serviceFuncResp(mockStor, itemData.ID)

		res, err := itemServ.AddFromCSV(tc.request.ctx, tc.request.data([]items.Item{itemData}))

		assert.Equal(t, tc.want.data(itemData), res, tc.name)
		assert.Equal(t, tc.want.err, err, tc.name)
	}
}

func TestReadData(t *testing.T) {
	type request struct {
		ctx context.Context

		defDataRequest []items.Item
	}

	type wantResp struct {
		data func([]items.Item) []byte
		err  error

		isErr bool
	}

	testCases := []struct {
		name            string
		serviceMock     *mocks.ItemSrv // service to mock
		serviceFuncResp func(*mocks.ItemSrv, []items.Item)
		request         request
		want            wantResp
	}{
		{
			name:        "get csv normal data",
			serviceMock: &mocks.ItemSrv{},
			serviceFuncResp: func(mc *mocks.ItemSrv, items []items.Item) {
				mc.On("GetAllItems", mock.Anything).Return(func(_ context.Context) []byte {
					res, err := json.Marshal(items)
					if err != nil {
						log.Fatalln(err)
					}

					return res
				}, nil)
			},
			request: request{ctx: context.Background(),
				defDataRequest: []items.Item{*defTestItem(), *defTestItem(), *defTestItem()},
			},
			want: wantResp{
				data: func(itemsData []items.Item) []byte {
					res, err := csvutil.Marshal(itemsData)
					if err != nil {
						log.Fatal(err)
					}

					return res
				},
				err:   nil,
				isErr: false,
			},
		},
		{
			name:        "return test error",
			serviceMock: &mocks.ItemSrv{},
			serviceFuncResp: func(mc *mocks.ItemSrv, items []items.Item) {
				mc.On("GetAllItems", mock.Anything).Return(func(_ context.Context) []byte {
					return nil
				}, er.DataNotExist)
			},
			request: request{
				ctx:            context.Background(),
				defDataRequest: []items.Item{*defTestItem(), *defTestItem(), *defTestItem(), *defTestItem()},
			},
			want: wantResp{
				data: func(_ []items.Item) []byte {
					return nil
				},
				err:   er.DataNotExist,
				isErr: true,
			},
		},
	}
	for _, tc := range testCases {
		mockStor := tc.serviceMock
		itemServ := New(mockStor)

		itemData := tc.request.defDataRequest
		tc.serviceFuncResp(mockStor, itemData)

		res, err := itemServ.GetAllItemsAsCSV(tc.request.ctx)

		assert.Equal(t, tc.want.data(itemData), res, tc.name)
		assert.Equal(t, tc.want.err, err, tc.name)
	}
}
