package memory

import (
	"encoding/json"
	"errors"

	itemsModel "github.com/DrIhor/test_task/internal/models/items"
)

type DB struct {
	items map[string]itemsModel.Item
}

func New() *DB {
	return &DB{
		items: make(map[string]itemsModel.Item),
	}
}

func (db *DB) AddNewItem(newItem itemsModel.Item) error {
	if _, ok := db.items[newItem.Name]; ok {
		return errors.New("Current row data exist. Try modify!")
	}

	db.items[newItem.Name] = newItem
	return nil
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

func (db *DB) GetItem(name string) ([]byte, error) {
	if _, ok := db.items[name]; !ok {
		return nil, errors.New("Current row data exist. Try modify!")
	}

	res, err := json.Marshal(db.items[name])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (db *DB) DeleteItem(itemName string) error {
	if _, ok := db.items[itemName]; !ok {
		return errors.New("Current row data not exist. Try to add!")
	}

	delete(db.items, itemName)
	return nil
}

func (db *DB) UpdateItem(itemName string) ([]byte, error) {
	val, ok := db.items[itemName]
	if !ok {
		return nil, errors.New("Current row data not exist. Try to add!")
	}

	if val.ItemsNumber-1 <= 0 {
		delete(db.items, itemName)
		return nil, errors.New("Item ended!!!")
	}

	// save new value
	val.ItemsNumber--
	db.items[itemName] = val

	// return result
	res, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	return res, nil
}
