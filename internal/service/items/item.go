package items

import (
	"context"

	itemModel "github.com/DrIhor/test_task/internal/models/items"
)

type ItemCRUDServ interface {
	AddNewItem(itemModel.Item) (string, error)
	GetAllItems() ([]byte, error)
	GetItem(string) ([]byte, error)
	DeleteItem(string) (bool, error)
	UpdateItem(string) ([]byte, error)
}

func New(ctx context.Context, stor itemModel.ItemStorageServices) *ItemServices {
	return &ItemServices{
		storage: stor,
		ctx:     ctx,
	}
}

func (itemSrv *ItemServices) AddNewItem(item itemModel.Item) (string, error) {
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
