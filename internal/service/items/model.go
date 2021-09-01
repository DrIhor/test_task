package items

import (
	itemModel "github.com/DrIhor/test_task/internal/models/items"
)

type ItemServices struct {
	storage itemModel.ItemStorageServices
}

type ItemSrv interface {
	ItemCRUDServ
	ItemCSVServ
}
