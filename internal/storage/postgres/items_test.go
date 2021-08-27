package postgres

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"testing"

	"github.com/DrIhor/test_task/internal/models/items"
	mocks "github.com/DrIhor/test_task/mocks/intern/models/items"
	"github.com/google/uuid"
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

func TestSaveToDB(t *testing.T) {

	testCases := []struct {
		name string

		// test mocks
		mock      mocks.ItemStorageServices
		caseCheck func(*mocks.ItemStorageServices)

		// test data
		testCtx  context.Context
		testItem items.Item

		// result
		expectedResult string
		hasErr         bool
	}{
		{
			name: "Add new good item",
			mock: mocks.ItemStorageServices{},
			caseCheck: func(mc *mocks.ItemStorageServices) {
				mc.On("AddNewItem",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
					mock.MatchedBy(func(it items.Item) bool { return it != items.Item{} }),
				).Return(
					func(ctx context.Context, s items.Item) string {
						return s.ID
					}, nil)
			},
			testCtx:  context.Background(),
			testItem: *defTestItem(),
			// expectedResult: "testItem.ID", //
			hasErr: false,
		},
	}

	for _, test := range testCases {
		tsMock := test.mock
		test.caseCheck(&tsMock)
		actual, err := tsMock.AddNewItem(test.testCtx, test.testItem)
		assert.Equal(t, test.testItem.ID, actual, test.name)

		if test.hasErr {
			assert.NotNil(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}

	}
}

func TestGetAllItems(t *testing.T) {

	// testItem :=
	testCases := []struct {
		name string

		// test mocks
		mock      mocks.ItemStorageServices
		caseCheck func(*mocks.ItemStorageServices, []items.Item)

		// test data
		testCtx   context.Context
		testItems []items.Item

		// result
		expectedResult func([]items.Item) []byte
		hasErr         bool
	}{
		{
			name: "Get empty data",
			mock: mocks.ItemStorageServices{},
			caseCheck: func(mc *mocks.ItemStorageServices, items []items.Item) {
				mc.On("GetAllItems",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
				).Return(nil, nil)
			},
			testCtx:   context.Background(),
			testItems: nil,
			expectedResult: func(items []items.Item) []byte {
				return nil
			},
			hasErr: false,
		},
		{
			name: "Get full data",
			mock: mocks.ItemStorageServices{},
			caseCheck: func(mc *mocks.ItemStorageServices, items []items.Item) {
				mc.On("GetAllItems",
					mock.MatchedBy(func(_ context.Context) bool { return true }),
				).Return(func(_ context.Context) []byte {
					res, err := json.Marshal(items)
					if err != nil {
						log.Fatalln("Marshal err: ", err)
					}

					return res
				}, nil)
			},
			testCtx:   context.Background(),
			testItems: []items.Item{*defTestItem(), *defTestItem()},
			expectedResult: func(items []items.Item) []byte {
				if len(items) == 0 {
					return nil
				}

				res, err := json.Marshal(items)
				if err != nil {
					log.Fatal(err)
				}

				return res
			},
			hasErr: false,
		},
	}

	for _, test := range testCases {
		tsMock := test.mock
		test.caseCheck(&tsMock, test.testItems)
		actual, err := tsMock.GetAllItems(test.testCtx)
		if ok := assert.Equal(t, test.expectedResult(test.testItems), actual, test.name); !ok {
			t.Error(test.name)
		}

		if test.hasErr {
			assert.NotNil(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}
	}
}
