package elk

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"

	itemsModel "github.com/DrIhor/test_task/internal/models/items"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

const indexItems string = "items" // must be lower case

type ElkStorage struct {
	client *elasticsearch.Client
}

func getElkConfig() *elasticsearch.Config {
	return &elasticsearch.Config{
		Addresses: []string{os.Getenv("ELASTIC_ADDR")},
		Username:  os.Getenv("ELASTIC_USER"),
		Password:  os.Getenv("ELASTIC_PASSWORD"),
	}
}

func New() *ElkStorage {
	conf := getElkConfig()
	es, err := elasticsearch.NewClient(*conf)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = es.Info(); err != nil {
		log.Fatal(err)
	}

	return &ElkStorage{client: es}
}

// search not exist id to add to db
// recommend add len of db items as first argument
func (db *ElkStorage) getNewID(ctx context.Context, id int) (int, error) {
	req := esapi.GetRequest{
		Index:      indexItems,
		DocumentID: strconv.Itoa(id),
	}

	resp, err := req.Do(ctx, db.client)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode == 404 {
		return id, nil
	}

	return db.getNewID(ctx, id+1)
}

func (db *ElkStorage) AddNewItem(ctx context.Context, newItem itemsModel.Item) (int, error) {
	data, err := json.Marshal(newItem)
	if err != nil {
		return 0, err
	}

	// get count of items
	countReq := esapi.CountRequest{
		Index: []string{indexItems},
	}

	countInfo := make(map[string]int)
	resp, err := countReq.Do(ctx, db.client)
	if err != nil {
		return 0, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&countInfo); err != nil {
		return 0, err
	}

	// check id vill be valid
	newInfoID := countInfo["count"] + 1
	newID, err := db.getNewID(ctx, newInfoID)
	if err != nil {
		return 0, err
	}

	// add new item
	req := esapi.CreateRequest{
		Index:      indexItems,
		DocumentID: strconv.Itoa(newID),
		Body:       bytes.NewReader(data),
	}

	_, err = req.Do(ctx, db.client)
	if err != nil {
		return 0, err
	}

	return newInfoID, nil
}

func (db *ElkStorage) GetAllItems(ctx context.Context) ([]byte, error) {
	req := esapi.SearchRequest{
		Index: []string{indexItems},
	}

	resp, err := req.Do(ctx, db.client)
	if err != nil {
		return nil, err
	}

	var searchList SearchELKItems
	err = json.NewDecoder(resp.Body).Decode(&searchList)
	if err != nil {
		return nil, err
	}

	var resultList []itemsModel.Item
	for _, val := range searchList.Hits.Hits {
		resultList = append(resultList, val.Source)

	}

	res, err := json.Marshal(resultList)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (db *ElkStorage) GetItem(ctx context.Context, id int) ([]byte, error) {
	req := esapi.GetRequest{
		Index:      indexItems,
		DocumentID: strconv.Itoa(id),
	}

	resp, err := req.Do(ctx, db.client)
	if err != nil {
		return nil, err
	}

	var item GetELKItems
	err = json.NewDecoder(resp.Body).Decode(&item)
	if err != nil {
		return nil, err
	}

	if !item.Found || resp.StatusCode == 404 {
		return nil, nil
	}

	res, err := json.Marshal(item.Source)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (db *ElkStorage) DeleteItem(ctx context.Context, id int) (bool, error) {
	req := esapi.DeleteRequest{
		Index:      indexItems,
		DocumentID: strconv.Itoa(id),
	}

	resp, err := req.Do(ctx, db.client)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == 404 {
		return false, nil
	}

	return true, nil
}

func (db *ElkStorage) createUpdateItemsNumberQuery() ([]byte, error) {
	var query UpdateItemsQuery

	query.Script.Source = "if( ctx._source.itemsNumber >= 1) {ctx._source.itemsNumber -= 1} else { ctx.op = 'delete' }"
	query.Script.Lang = "painless"

	return json.Marshal(query)
}

func (db *ElkStorage) UpdateItem(ctx context.Context, id int) ([]byte, error) {
	query, err := db.createUpdateItemsNumberQuery()
	if err != nil {
		return nil, err
	}

	req := esapi.UpdateRequest{
		Index:      indexItems,
		DocumentID: strconv.Itoa(id),

		Body: bytes.NewReader(query),
	}

	_, err = req.Do(ctx, db.client)
	if err != nil {
		return nil, err
	}

	return db.GetItem(ctx, id)
}
