package items

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"strings"

	itemModel "github.com/DrIhor/test_task/internal/models/items"
	"github.com/jszwec/csvutil"
)

type ItemCSVServ interface {
	AddFromCSV(rd *csv.Reader) ([]byte, error)
	GetAllItemsAsCSV() ([]byte, error)
}

func (itemSrv *ItemServices) AddFromCSV(rd *csv.Reader) ([]byte, error) {
	var (
		itemHeader string          // struct fields
		firstRow   bool     = true // if file header
		newIDs     []string        // result data
	)

	for {
		row, err := rd.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		// init header
		if firstRow {
			itemHeader = strings.Join(row, ",") + "\n" // create single row for csvutil lib
			firstRow = false
			continue
		}

		csvIteam := []byte(itemHeader + strings.Join(row, ","))
		var items []itemModel.Item
		if err := csvutil.Unmarshal(csvIteam, &items); err != nil {
			return nil, err
		}

		for _, item := range items {
			if item != (itemModel.Item{}) {
				id, err := itemSrv.AddNewItem(item)
				if err != nil {
					return nil, err
				}
				newIDs = append(newIDs, id)
			}
		}
	}

	res, err := json.Marshal(newIDs)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (itemSrv *ItemServices) GetAllItemsAsCSV() ([]byte, error) {
	byteData, err := itemSrv.GetAllItems()
	if err != nil {
		return nil, err
	}
	var itemsSlice []itemModel.Item
	if err := json.Unmarshal(byteData, &itemsSlice); err != nil {
		return nil, err
	}

	b, err := csvutil.Marshal(itemsSlice)
	if err != nil {
		return nil, err
	}

	return b, nil
}
