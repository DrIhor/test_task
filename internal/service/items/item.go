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

func (itemSrv *ItemServices) AddNewItem(item items.Item) (int, error) {
	return itemSrv.storage.AddNewItem(item)
}

func (itemSrv *ItemServices) GetAllItems() ([]byte, error) {
	return itemSrv.storage.GetAllItems()
}

func (itemSrv *ItemServices) GetItem(id int) ([]byte, error) {
	return itemSrv.storage.GetItem(id)
}

func (itemSrv *ItemServices) DeleteItem(id int) (bool, error) {
	return itemSrv.storage.DeleteItem(id)
}

func (itemSrv *ItemServices) UpdateItem(id int) ([]byte, error) {
	return itemSrv.storage.UpdateItem(id)
}
