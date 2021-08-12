package connectors

import (
	"context"
	"log"
	"time"

	"github.com/DrIhor/test_task/internal/models/items"
	pb "github.com/DrIhor/test_task/pkg/grpc"
	"google.golang.org/grpc"
)

type GrpcConnector struct {
	itemStor pb.ItemStorageClient

	conn *grpc.ClientConn
}

func NewGRPC(address string) *GrpcConnector {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	return &GrpcConnector{
		itemStor: pb.NewItemStorageClient(conn),
		conn:     conn,
	}
}

func (gr *GrpcConnector) AddNewItem(item items.Item) (int, error) {
	reqItem := pb.Item{
		Name:        item.Name,
		Price:       item.Price,
		ItemsNumber: item.ItemsNumber,
		Description: item.Description,
	}

	res, err := gr.itemStor.AddNewItem(context.Background(), &reqItem)

	time.Sleep(2 * time.Second)

	return int(res.ID), err
}

func (gr *GrpcConnector) GetItem(id int) ([]byte, error) {
	res, err := gr.itemStor.GetItem(context.Background(), &pb.ItemID{
		ID: int64(id),
	})

	if err != nil {
		return nil, err
	}

	return res.Result, err
}

func (gr *GrpcConnector) DeleteItem(id int) (bool, error) {
	res, err := gr.itemStor.DeleteItem(context.Background(), &pb.ItemID{
		ID: int64(id),
	})
	return res.DoneWork, err
}

func (gr *GrpcConnector) GetAllItems() ([]byte, error) {
	res, err := gr.itemStor.GetAllItems(context.Background(), &pb.NoneObjectRequest{})
	if err != nil {
		return nil, err
	}

	return res.Result, err
}

func (gr *GrpcConnector) UpdateItem(id int) ([]byte, error) {
	res, err := gr.itemStor.UpdateItem(context.Background(), &pb.ItemID{
		ID: int64(id),
	})
	if err != nil {
		return nil, err
	}

	return res.Result, err
}
