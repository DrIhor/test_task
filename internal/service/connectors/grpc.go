package connectors

import (
	"context"
	"log"

	"github.com/DrIhor/test_task/internal/models/items"
	pb "github.com/DrIhor/test_task/pkg/grpc"
	"google.golang.org/grpc"
)

type GrpcConnector struct {
	itemStor pb.ItemStorageClient

	conn *grpc.ClientConn
}

func NewGRPC(address string) *GrpcConnector {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	return &GrpcConnector{
		itemStor: pb.NewItemStorageClient(conn),
		conn:     conn,
	}
}

func (gr *GrpcConnector) CloseConn() {
	gr.conn.Close()
}

func (gr *GrpcConnector) AddNewItem(item items.Item) (*pb.NoneObjectResp, error) {
	reqItem := pb.Item{
		Name:        item.Name,
		Price:       item.Price,
		ItemsNumber: item.ItemsNumber,
		Description: item.Description,
	}

	return gr.itemStor.AddNewItem(context.Background(), &reqItem)
}

func (gr *GrpcConnector) GetItem(itemName string) ([]byte, error) {
	res, err := gr.itemStor.GetItem(context.Background(), &pb.ItemName{
		Name: itemName,
	})

	return res.Result, err
}

func (gr *GrpcConnector) DeleteItem(itemName string) (*pb.NoneObjectResp, error) {
	return gr.itemStor.DeleteItem(context.Background(), &pb.ItemName{
		Name: itemName,
	})
}

func (gr *GrpcConnector) GetAllItems() ([]byte, error) {
	res, err := gr.itemStor.GetAllItems(context.Background(), &pb.NoneObjectRequest{})
	return res.Result, err
}

func (gr *GrpcConnector) UpdateItem(itemName string) ([]byte, error) {
	res, err := gr.itemStor.UpdateItem(context.Background(), &pb.ItemName{
		Name: itemName,
	})

	return res.Result, err
}
