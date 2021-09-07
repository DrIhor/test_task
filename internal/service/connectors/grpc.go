package connectors

import (
	"context"
	"time"

	"github.com/DrIhor/test_task/internal/models/items"
	pb "github.com/DrIhor/test_task/pkg/grpc"
	"google.golang.org/grpc"
)

type GrpcConnector struct {
	itemStor pb.ItemStorageClient
	conn     *grpc.ClientConn
}

// create new GRPC connection with items storage services
func NewGRPC(ctx context.Context, address string) (*GrpcConnector, error) {
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &GrpcConnector{
		itemStor: pb.NewItemStorageClient(conn),
		conn:     conn,
	}, nil
}

/**
 * service calls
 */

func (gr *GrpcConnector) AddNewItem(ctx context.Context, item items.Item) (string, error) {
	reqItem := pb.Item{
		Name:        item.Name,
		Price:       item.Price,
		ItemsNumber: item.ItemsNumber,
		Description: item.Description,
	}

	res, err := gr.itemStor.AddNewItem(ctx, &reqItem)

	return res.ID, err
}

func (gr *GrpcConnector) GetItem(ctx context.Context, id string) ([]byte, error) {
	res, err := gr.itemStor.GetItem(ctx, &pb.ItemID{
		ID: id,
	})

	if err != nil {
		return nil, err
	}

	return res.Result, err
}

func (gr *GrpcConnector) DeleteItem(ctx context.Context, id string) (bool, error) {
	res, err := gr.itemStor.DeleteItem(ctx, &pb.ItemID{
		ID: id,
	})
	return res.DoneWork, err
}

func (gr *GrpcConnector) GetAllItems(ctx context.Context) ([]byte, error) {
	res, err := gr.itemStor.GetAllItems(ctx, &pb.NoneObjectRequest{})
	if err != nil {
		return nil, err
	}

	return res.Result, err
}

func (gr *GrpcConnector) UpdateItem(ctx context.Context, id string) ([]byte, error) {
	res, err := gr.itemStor.UpdateItem(ctx, &pb.ItemID{
		ID: id,
	})
	if err != nil {
		return nil, err
	}

	return res.Result, err
}
