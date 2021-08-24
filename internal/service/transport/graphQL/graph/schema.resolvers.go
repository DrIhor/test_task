package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	itemModel "github.com/DrIhor/test_task/internal/models/items"
	"github.com/DrIhor/test_task/internal/service/connectors"
	"github.com/DrIhor/test_task/internal/service/transport/graphQL/graph/generated"
	"github.com/DrIhor/test_task/internal/service/transport/graphQL/graph/model"
	"github.com/google/uuid"
)

func (r *mutationResolver) AddItem(ctx context.Context, item model.Iteminput) (*string, error) {
	itm := itemModel.Item{
		Name:        item.Name,
		Price:       int32(item.Price),
		ItemsNumber: int32(item.ItemsNumber),
		Description: *item.Desc,
	}
	id, err := r.services.AddNewItem(ctx, itm)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *mutationResolver) UpdatePerson(ctx context.Context, id string) (*model.Item, error) {

	var errData error
	var res []byte
	switch os.Getenv("STORAGE_TYPE") {
	case "":
		res, errData = r.services.UpdateItem(ctx, id)
	case "grpc":
		grpcConn := connectors.NewGRPC(ctx, os.Getenv("GRCP_ADDR"))
		res, errData = grpcConn.UpdateItem(ctx, id)
	}

	// check if struct is empty
	if len(res) < 3 {
		return nil, errors.New("Not found")
	}

	if errData != nil {
		return nil, errData
	}

	var resultItem model.Item
	err := json.Unmarshal(res, &resultItem)
	if err != nil {
		return nil, err
	}

	// add id of person
	resultItem.ID = &id

	return &resultItem, nil
}

func (r *mutationResolver) DeletePerson(ctx context.Context, id string) (*bool, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var done bool
	var errData error
	switch os.Getenv("STORAGE_TYPE") {
	case "":
		done, errData = r.services.DeleteItem(ctx, id)
	case "grpc":
		grpcConn := connectors.NewGRPC(ctx, os.Getenv("GRCP_ADDR"))
		done, errData = grpcConn.DeleteItem(ctx, id)
	}

	return &done, errData
}

func (r *queryResolver) GetItems(ctx context.Context) ([]*model.Item, error) {
	var errData error
	var res []byte
	switch os.Getenv("STORAGE_TYPE") {
	case "":
		res, errData = r.services.GetAllItems(ctx)
	case "grpc":
		grpcConn := connectors.NewGRPC(ctx, os.Getenv("GRCP_ADDR"))
		res, errData = grpcConn.GetAllItems(ctx)
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
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var errData error
	var res []byte
	switch os.Getenv("STORAGE_TYPE") {
	case "":
		res, errData = r.services.UpdateItem(ctx, id)
	case "grpc":
		grpcConn := connectors.NewGRPC(ctx, os.Getenv("GRCP_ADDR"))
		res, errData = grpcConn.GetItem(ctx, id)
	}

	if errData != nil {
		return nil, errData
	}

	var resultData model.Item
	err = json.Unmarshal(res, &resultData)
	if err != nil {
		return nil, err
	}

	resultData.ID = &id

	return &resultData, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
