package memory

import (
	"encoding/json"
	"errors"

	iteamsModel "github.com/DrIhor/test_task/internal/models/iteams"
)

type IteamsStorage struct {
	list map[string]iteamsModel.Iteam
}

func NewIteamList() *IteamsStorage {
	return &IteamsStorage{
		list: make(map[string]iteamsModel.Iteam),
	}
}

func (iteams *IteamsStorage) AddNewIteam(newIteam iteamsModel.Iteam) error {
	if _, ok := iteams.list[newIteam.Name]; ok {
		return errors.New("Current row data exist. Try modify!")
	}

	iteams.list[newIteam.Name] = newIteam
	return nil
}

func (iteams *IteamsStorage) GetAllIteams() ([]byte, error) {
	var iteamsSlice []iteamsModel.Iteam
	for _, value := range iteams.list {
		iteamsSlice = append(iteamsSlice, value)
	}

	res, err := json.Marshal(iteamsSlice)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (iteams *IteamsStorage) GetIteam(name string) ([]byte, error) {
	if _, ok := iteams.list[name]; !ok {
		return nil, errors.New("Current row data exist. Try modify!")
	}

	res, err := json.Marshal(iteams.list[name])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (iteams *IteamsStorage) DeleteIteam(iteamName string) error {
	if _, ok := iteams.list[iteamName]; !ok {
		return errors.New("Current row data not exist. Try to add!")
	}

	delete(iteams.list, iteamName)
	return nil
}

func (iteams *IteamsStorage) UpdateIteam(iteamName string) ([]byte, error) {
	val, ok := iteams.list[iteamName]
	if !ok {
		return nil, errors.New("Current row data not exist. Try to add!")
	}

	if val.IteamsNumber-1 <= 0 {
		delete(iteams.list, iteamName)
		return nil, errors.New("Iteam ended!!!")
	}

	// save new value
	val.IteamsNumber--
	iteams.list[iteamName] = val

	// return result
	res, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	return res, nil
}
