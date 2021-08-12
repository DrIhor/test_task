package memory

import (
	"encoding/json"
	"errors"

	itemsModel "github.com/DrIhor/test_task/internal/models/items"
)

type DB struct {
	items map[int]itemsModel.Item
}

func New() *DB {
	return &DB{
		items: make(map[int]itemsModel.Item),
	}
}

func (db *DB) AddNewItem(newItem itemsModel.Item) (int, error) {
	id := len(db.items) + 1
	if _, ok := db.items[id]; ok {
		return 0, errors.New("RowExist")
	}

	db.items[id] = newItem
	return id, nil
}

func (db *DB) GetAllItems() ([]byte, error) {
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

func (db *DB) GetItem(id int) ([]byte, error) {
	if _, ok := db.items[id]; !ok {
		return nil, errors.New("NotExist")
	}

	res, err := json.Marshal(db.items[id])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (db *DB) DeleteItem(id int) (bool, error) {
	if _, ok := db.items[id]; !ok {
		return false, nil
	}

	delete(db.items, id)
	return true, nil
}

func (db *DB) UpdateItem(id int) ([]byte, error) {
	val, ok := db.items[id]
	if !ok {
		return nil, errors.New("NotExist")
	}

	if val.ItemsNumber-1 <= 0 {
		delete(db.items, id)
		return nil, errors.New("WrongNumber")
	}

	// save new value
	val.ItemsNumber--
	db.items[id] = val

	// return result
	res, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	return res, nil
}
