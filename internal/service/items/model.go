package items

import (
	"context"

	itemModel "github.com/DrIhor/test_task/internal/models/items"
)

type ItemServices struct {
	storage itemModel.ItemStorageServices
	ctx     context.Context
}

type ItemSrv interface {
	ItemCRUDServ
	ItemCSVServ
}
