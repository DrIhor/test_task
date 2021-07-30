package memory

import (
	"encoding/json"
	"errors"

	iteamsModel "github.com/DrIhor/test_task/internal/models/iteams"
)

type DB struct {
	items map[string]iteamsModel.Iteam
}

func New() *DB {
	return &DB{
		items: make(map[string]iteamsModel.Iteam),
	}
}

func (db *DB) AddNewIteam(newIteam iteamsModel.Iteam) error {
	if _, ok := db.items[newIteam.Name]; ok {
		return errors.New("Current row data exist. Try modify!")
	}

	db.items[newIteam.Name] = newIteam
	return nil
}

func (db *DB) GetAllIteams() ([]byte, error) {
	var iteamsSlice []iteamsModel.Iteam
	for _, value := range db.items {
		iteamsSlice = append(iteamsSlice, value)
	}

	res, err := json.Marshal(iteamsSlice)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (db *DB) GetIteam(name string) ([]byte, error) {
	if _, ok := db.items[name]; !ok {
		return nil, errors.New("Current row data exist. Try modify!")
	}

	res, err := json.Marshal(db.items[name])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (db *DB) DeleteIteam(iteamName string) error {
	if _, ok := db.items[iteamName]; !ok {
		return errors.New("Current row data not exist. Try to add!")
	}

	delete(db.items, iteamName)
	return nil
}

func (db *DB) UpdateIteam(iteamName string) ([]byte, error) {
	val, ok := db.items[iteamName]
	if !ok {
		return nil, errors.New("Current row data not exist. Try to add!")
	}

	if val.IteamsNumber-1 <= 0 {
		delete(db.items, iteamName)
		return nil, errors.New("Iteam ended!!!")
	}

	// save new value
	val.IteamsNumber--
	db.items[iteamName] = val

	// return result
	res, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	return res, nil
}
