package items

import (
	"context"
	"encoding/json"
)

type Item struct {
	ID          int    `bson:"_id" json:"id,omitempty" csv:"-,omitempty"`
	Name        string `bson:"name,omitempty" json:"name,omitempty" csv:"name"`
	Price       int32  `bson:"price,omitempty" json:"price,omitempty" csv:"price"`
	ItemsNumber int32  `bson:"itemsNumber,omitempty" json:"itemsNumber,omitempty" csv:"itemsNumber"`
	Description string `bson:"desc,omitempty" json:"desc,omitempty" csv:"desc"`
}

func (i *Item) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

// all main services for Item to work with DB
type ItemStorageServices interface {
	AddNewItem(context.Context, Item) (int, error)
	GetAllItems(context.Context) ([]byte, error)
	GetItem(context.Context, int) ([]byte, error)
	DeleteItem(context.Context, int) (bool, error)
	UpdateItem(context.Context, int) ([]byte, error)
}
