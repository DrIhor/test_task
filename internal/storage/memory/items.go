package memory

import (
	"context"
	"encoding/json"

	er "github.com/DrIhor/test_task/internal/errors"

	itemsModel "github.com/DrIhor/test_task/internal/models/items"
	"github.com/google/uuid"
)

type DB struct {
	items map[uuid.UUID]itemsModel.Item
}

func New() *DB {
	return &DB{
		items: make(map[uuid.UUID]itemsModel.Item),
	}
}

// search not exist id to add to db
// recommend add len of db items as first argument
func (db *DB) getNewID(id uuid.UUID) uuid.UUID {

	if _, ok := db.items[id]; ok {
		return db.getNewID(uuid.New())
	}

	return id
}

func (db *DB) AddNewItem(ctx context.Context, newItem itemsModel.Item) (string, error) {
	id := db.getNewID(uuid.New())

	db.items[id] = newItem
	return id.String(), nil
}

func (db *DB) GetAllItems(ctx context.Context) ([]byte, error) {
	var itemsSlice []itemsModel.Item
	for _, value := range db.items {
		itemsSlice = append(itemsSlice, value)
	}

	res, err := json.Marshal(itemsSlice)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (db *DB) GetItem(ctx context.Context, id string) ([]byte, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	if _, ok := db.items[uid]; !ok {
		return nil, er.DataNotExist
	}

	res, err := json.Marshal(db.items[uid])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (db *DB) DeleteItem(ctx context.Context, id string) (bool, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return false, err
	}

	if _, ok := db.items[uid]; !ok {
		return false, nil
	}

	delete(db.items, uid)
	return true, nil
}

func (db *DB) UpdateItem(ctx context.Context, id string) ([]byte, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	val, ok := db.items[uid]
	if !ok {
		return nil, er.DataNotExist
	}

	if val.ItemsNumber-1 < 0 {
		delete(db.items, uid)
		return nil, nil
	}

	// save new value
	val.ItemsNumber--
	db.items[uid] = val

	// return result
	res, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	return res, nil
}
