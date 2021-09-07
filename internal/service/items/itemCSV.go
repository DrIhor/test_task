package items

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"strings"

	itemModel "github.com/DrIhor/test_task/internal/models/items"
	"github.com/jszwec/csvutil"
)

// item services using CSV files
type ItemCSVServ interface {
	AddFromCSV(context.Context, *csv.Reader) ([]byte, error)
	GetAllItemsAsCSV(context.Context) ([]byte, error)
}

/**
 * main services logic
 */

// read csv file and write into DB as records
func (itemSrv *ItemServices) AddFromCSV(ctx context.Context, rd *csv.Reader) ([]byte, error) {
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

		// transform data
		csvIteam := []byte(itemHeader + strings.Join(row, ","))
		var items []itemModel.Item
		if err := csvutil.Unmarshal(csvIteam, &items); err != nil {
			return nil, err
		}

		// save data
		for _, item := range items {
			if item != (itemModel.Item{}) {
				id, err := itemSrv.AddNewItem(ctx, item)
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

// read data from DB and return to user as csv file
func (itemSrv *ItemServices) GetAllItemsAsCSV(ctx context.Context) ([]byte, error) {
	byteData, err := itemSrv.GetAllItems(ctx)
	if err != nil {
		return nil, err
	}

	// convert data into items to return using `csvutil` marshal
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
