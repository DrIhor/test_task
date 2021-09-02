package items

import (
	"context"

	itemModel "github.com/DrIhor/test_task/internal/models/items"
)

type ItemCRUDServ interface {
	AddNewItem(context.Context, itemModel.Item) (string, error)
	GetAllItems(context.Context) ([]byte, error)
	GetItem(context.Context, string) ([]byte, error)
	DeleteItem(context.Context, string) (bool, error)
	UpdateItem(context.Context, string) ([]byte, error)
}

func New(stor itemModel.ItemStorageServices) *ItemServices {
	return &ItemServices{
		storage: stor,
	}
}

func (itemSrv *ItemServices) AddNewItem(ctx context.Context, item itemModel.Item) (string, error) {
	return itemSrv.storage.AddNewItem(ctx, item)
}

func (itemSrv *ItemServices) GetAllItems(ctx context.Context) ([]byte, error) {
	return itemSrv.storage.GetAllItems(ctx)
}

func (itemSrv *ItemServices) GetItem(ctx context.Context, id string) ([]byte, error) {
	return itemSrv.storage.GetItem(ctx, id)
}

func (itemSrv *ItemServices) DeleteItem(ctx context.Context, id string) (bool, error) {
	return itemSrv.storage.DeleteItem(ctx, id)
}

func (itemSrv *ItemServices) UpdateItem(ctx context.Context, id string) ([]byte, error) {
	return itemSrv.storage.UpdateItem(ctx, id)
}
