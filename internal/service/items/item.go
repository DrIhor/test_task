package items

import (
	"context"

	itemModel "github.com/DrIhor/test_task/internal/models/items"
)

// all services from item
type ItemCRUDServ interface {
	AddNewItem(context.Context, itemModel.Item) (string, error)
	GetAllItems(context.Context) ([]byte, error)
	GetItem(context.Context, string) ([]byte, error)
	DeleteItem(context.Context, string) (bool, error)
	UpdateItem(context.Context, string) ([]byte, error)
}

// init item services
func New(stor itemModel.ItemStorageServices) *ItemServices {
	return &ItemServices{
		storage: stor,
	}
}

/**
 * CRUD items services logic
 */

// add new record
func (itemSrv *ItemServices) AddNewItem(ctx context.Context, item itemModel.Item) (string, error) {
	return itemSrv.storage.AddNewItem(ctx, item)
}

// return each item from DB
func (itemSrv *ItemServices) GetAllItems(ctx context.Context) ([]byte, error) {
	return itemSrv.storage.GetAllItems(ctx)
}

// return single item
func (itemSrv *ItemServices) GetItem(ctx context.Context, id string) ([]byte, error) {
	return itemSrv.storage.GetItem(ctx, id)
}

// remove item from DB
func (itemSrv *ItemServices) DeleteItem(ctx context.Context, id string) (bool, error) {
	return itemSrv.storage.DeleteItem(ctx, id)
}

// update number of items
func (itemSrv *ItemServices) UpdateItem(ctx context.Context, id string) ([]byte, error) {
	return itemSrv.storage.UpdateItem(ctx, id)
}
