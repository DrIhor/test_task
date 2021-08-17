package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strconv"

	itemModel "github.com/DrIhor/test_task/internal/models/items"
	"github.com/DrIhor/test_task/internal/service/connectors"
	"github.com/DrIhor/test_task/internal/service/transport/graphQL/graph/generated"
	"github.com/DrIhor/test_task/internal/service/transport/graphQL/graph/model"
)

func (r *mutationResolver) AddItem(ctx context.Context, item model.Iteminput) (*int, error) {
	itm := itemModel.Item{
		Name:        item.Name,
		Price:       int32(item.Price),
		ItemsNumber: int32(item.ItemsNumber),
		Description: *item.Desc,
	}
	id, err := r.services.AddNewItem(itm)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *mutationResolver) UpdatePerson(ctx context.Context, id string) (*model.Item, error) {
	personID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var errData error
	var res []byte
	switch os.Getenv("STORAGE_TYPE") {
	case "":
		res, errData = r.services.UpdateItem(personID)
	case "grpc":
		grpcConn := connectors.NewGRPC(os.Getenv("GRCP_ADDR"))
		res, errData = grpcConn.UpdateItem(personID)
	}

	// check if struct is empty
	if len(res) < 3 {
		return nil, errors.New("Not found")
	}

	if errData != nil {
		return nil, errData
	}

	var resultItem model.Item
	err = json.Unmarshal(res, &resultItem)
	if err != nil {
		return nil, err
	}

	// add id of person
	resultItem.ID = &personID

	return &resultItem, nil
}

func (r *mutationResolver) DeletePerson(ctx context.Context, id string) (*bool, error) {
	personID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var done bool
	var errData error
	switch os.Getenv("STORAGE_TYPE") {
	case "":
		done, errData = r.services.DeleteItem(personID)
	case "grpc":
		grpcConn := connectors.NewGRPC(os.Getenv("GRCP_ADDR"))
		done, errData = grpcConn.DeleteItem(personID)
	}

	return &done, errData
}

func (r *queryResolver) GetItems(ctx context.Context) ([]*model.Item, error) {
	var errData error
	var res []byte
	switch os.Getenv("STORAGE_TYPE") {
	case "":
		res, errData = r.services.GetAllItems()
	case "grpc":
		grpcConn := connectors.NewGRPC(os.Getenv("GRCP_ADDR"))
		res, errData = grpcConn.GetAllItems()
	}

	if errData != nil {
		return nil, errData
	}

	var resultData []*model.Item
	err := json.Unmarshal(res, &resultData)
	if err != nil {
		return nil, err
	}

	return resultData, nil
}

func (r *queryResolver) GetItem(ctx context.Context, id string) (*model.Item, error) {
	personID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	var errData error
	var res []byte
	switch os.Getenv("STORAGE_TYPE") {
	case "":
		res, errData = r.services.UpdateItem(personID)
	case "grpc":
		grpcConn := connectors.NewGRPC(os.Getenv("GRCP_ADDR"))
		res, errData = grpcConn.GetItem(personID)
	}

	if errData != nil {
		return nil, errData
	}

	var resultData model.Item
	err = json.Unmarshal(res, &resultData)
	if err != nil {
		return nil, err
	}

	resultData.ID = &personID

	return &resultData, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
