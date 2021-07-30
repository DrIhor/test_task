package routes

import (
	"github.com/DrIhor/test_task/internal/models/items"
	itemModel "github.com/DrIhor/test_task/internal/models/items"
)

type ItemServices struct {
	storage itemModel.ItemStorageServices
}

func New(stor itemModel.ItemStorageServices) *ItemServices {
	return &ItemServices{
		storage: stor,
	}
}

func (itemSrv *ItemServices) AddNewItem(item items.Item) error {
	return itemSrv.storage.AddNewItem(item)
}

func (itemSrv *ItemServices) GetItem(itemName string) ([]byte, error) {
	return itemSrv.storage.GetItem(itemName)
}

func (itemSrv *ItemServices) DeleteItem(itemName string) error {
	return itemSrv.storage.DeleteItem(itemName)
}

func (itemSrv *ItemServices) UpdateItem(itemName string) ([]byte, error) {
	return itemSrv.storage.UpdateItem(itemName)
}
