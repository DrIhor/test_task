package memory

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/DrIhor/test_task/internal/models/items"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

var testItem = items.Item{
	Name:        "testItem",
	Price:       10,
	ItemsNumber: 1, // one item delete and next must return null data
	Description: "test descr",
}

var db *DB // single connection to db

func init() {
	db = New()
}

func TestAddItem(t *testing.T) {
	id, err := db.AddNewItem(context.Background(), testItem)

	assert.Nil(t, err, "add item error")
	assert.True(t, isValidUUID(id), "add item uuid validation")

	testItem.ID = id // add for other tests
}

func TestGetAllItems(t *testing.T) {
	_, err := db.GetAllItems(context.Background())

	assert.Nil(t, err, "get all items err")
}

func TestUpdateItem(t *testing.T) {
	type data struct {
		item *items.Item
		ctx  context.Context
	}

	type expect struct {
		itemsNumber     func(items.Item) int32
		isErr           bool
		endData         bool
		checkReturnData bool
	}

	testCases := []struct {
		name    string
		request data
		want    expect
	}{
		{
			name: "update item data",
			request: data{
				item: &testItem,
				ctx:  context.Background(),
			},
			want: expect{
				itemsNumber: func(currentItem items.Item) int32 {
					return currentItem.ItemsNumber - 1
				},
				isErr:           false,
				endData:         false,
				checkReturnData: true,
			},
		},
		{
			name: "update when data end",
			request: data{
				item: &testItem,
				ctx:  context.Background(),
			},
			want: expect{
				itemsNumber: func(_ items.Item) int32 {
					return 0
				},
				isErr:           false,
				endData:         true,
				checkReturnData: false,
			},
		},
	}

	for _, tc := range testCases {
		data, err := db.UpdateItem(tc.request.ctx, tc.request.item.ID)

		if tc.want.isErr {
			assert.Error(t, err, tc.name)
		} else {
			assert.Nil(t, err, tc.name)
		}

		if tc.want.endData {
			assert.Nil(t, data, tc.name)
		} else {
			assert.NotNil(t, data, tc.name)
		}

		if tc.want.checkReturnData {
			var resItem items.Item
			errUnmarsh := json.Unmarshal(data, &resItem)
			if errUnmarsh != nil {
				log.Fatal(errUnmarsh)
			}

			assert.Equal(t, tc.want.itemsNumber(*tc.request.item), resItem.ItemsNumber, tc.name)
		} else {
			assert.Nil(t, data, tc.name)
		}

	}
}

func TestDeleteItem(t *testing.T) {
	_, err := db.DeleteItem(context.Background(), testItem.ID)

	assert.Nil(t, err, "delete item err")
}
