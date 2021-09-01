package gRPC

import (
	"context"

	"github.com/DrIhor/test_task/internal/models/items"
	pb "github.com/DrIhor/test_task/pkg/grpc"
)

func (s *Server) AddNewItem(ctx context.Context, in *pb.Item) (*pb.ItemID, error) {
	item := items.Item{
		Name:        in.Name,
		Price:       in.Price,
		ItemsNumber: in.ItemsNumber,
		Description: in.Description,
	}

	id, err := s.storage.AddNewItem(ctx, item)
	if err != nil {
		return nil, err
	}

	return &pb.ItemID{
		ID: id,
	}, nil
}

func (s *Server) GetAllItems(ctx context.Context, in *pb.NoneObjectRequest) (*pb.EncodeItemResponse, error) {

	res, err := s.storage.GetAllItems(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.EncodeItemResponse{
		Result: res,
	}, nil
}

func (s *Server) GetItem(ctx context.Context, in *pb.ItemID) (*pb.EncodeItemResponse, error) {

	res, err := s.storage.GetItem(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	return &pb.EncodeItemResponse{
		Result: res,
	}, nil
}

func (s *Server) DeleteItem(ctx context.Context, in *pb.ItemID) (*pb.NoneObjectResp, error) {

	done, err := s.storage.DeleteItem(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	return &pb.NoneObjectResp{
		DoneWork: done,
	}, nil
}

func (s *Server) UpdateItem(ctx context.Context, in *pb.ItemID) (*pb.EncodeItemResponse, error) {
	res, err := s.storage.UpdateItem(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	return &pb.EncodeItemResponse{
		Result: res,
	}, nil
}
