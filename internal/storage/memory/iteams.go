package memory

import (
	"encoding/json"
	"errors"
)

type Iteam struct {
	Name         string `json:"name"`
	Price        int    `json:"price"`
	IteamsNumber int    `json:"iteamsNumber"`
	Description  string `json:"desc"`
}

var IteamList = map[string]Iteam{}

func AddNewIteam(newIteam Iteam) error {
	if _, ok := IteamList[newIteam.Name]; ok {
		return errors.New("Current row data exist. Try modify!")
	}

	IteamList[newIteam.Name] = newIteam
	return nil
}

func ShowAllIteams() ([]byte, error) {
	var iteamsSlice []Iteam
	for _, value := range IteamList {
		iteamsSlice = append(iteamsSlice, value)
	}

	res, err := json.Marshal(iteamsSlice)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func ShowIteam(name string) ([]byte, error) {
	if _, ok := IteamList[name]; !ok {
		return nil, errors.New("Current row data exist. Try modify!")
	}

	res, err := json.Marshal(IteamList[name])
	if err != nil {
		return nil, err
	}

	return res, nil
}

func DeleteIteam(iteamName string) error {
	if _, ok := IteamList[iteamName]; !ok {
		return errors.New("Current row data not exist. Try to add!")
	}

	delete(IteamList, iteamName)
	return nil
}

func UpdateIteam(iteamName string) ([]byte, error) {
	val, ok := IteamList[iteamName]
	if !ok {
		return nil, errors.New("Current row data not exist. Try to add!")
	}

	if val.IteamsNumber-1 <= 0 {
		delete(IteamList, iteamName)
		return nil, errors.New("Iteam ended!!!")
	}

	// save new value
	val.IteamsNumber--
	IteamList[iteamName] = val

	// return result
	res, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	return res, nil
}
