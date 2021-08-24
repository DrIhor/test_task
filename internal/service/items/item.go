package routes

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"strings"

	"github.com/DrIhor/test_task/internal/models/items"
	itemModel "github.com/DrIhor/test_task/internal/models/items"
	"github.com/jszwec/csvutil"
)

type ItemServices struct {
	storage itemModel.ItemStorageServices
	ctx     context.Context
}

func New(ctx context.Context, stor itemModel.ItemStorageServices) *ItemServices {
	return &ItemServices{
		storage: stor,
		ctx:     ctx,
	}
}

func (itemSrv *ItemServices) AddNewItem(item items.Item) (string, error) {
	return itemSrv.storage.AddNewItem(itemSrv.ctx, item)
}

func (itemSrv *ItemServices) GetAllItems() ([]byte, error) {
	return itemSrv.storage.GetAllItems(itemSrv.ctx)
}

func (itemSrv *ItemServices) GetItem(id string) ([]byte, error) {
	return itemSrv.storage.GetItem(itemSrv.ctx, id)
}

func (itemSrv *ItemServices) DeleteItem(id string) (bool, error) {
	return itemSrv.storage.DeleteItem(itemSrv.ctx, id)
}

func (itemSrv *ItemServices) UpdateItem(id string) ([]byte, error) {
	return itemSrv.storage.UpdateItem(itemSrv.ctx, id)
}

func (itemSrv *ItemServices) AddFromCSV(rd *csv.Reader) ([]byte, error) {
	var (
		itemHeader string          // struct fields
		firstRow   bool     = true // if file header
		newIDs     []string        // result data
	)

	for {
		row, err := rd.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
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
			return nil, err
		}

		for _, item := range items {
			if item != (itemModel.Item{}) {
				id, err := itemSrv.AddNewItem(item)
				if err != nil {
					return nil, err
				}
				newIDs = append(newIDs, id)
			}
		}
	}

	res, err := json.Marshal(newIDs)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (itemSrv *ItemServices) GetAllItemsAsCSV() ([]byte, error) {
	byteData, err := itemSrv.GetAllItems()
	if err != nil {
		return nil, err
	}
	var itemsSlice []itemModel.Item
	if err := json.Unmarshal(byteData, &itemsSlice); err != nil {
		return nil, err
	}

	b, err := csvutil.Marshal(itemsSlice)
	if err != nil {
		return nil, err
	}

	return b, nil
}
