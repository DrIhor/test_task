package gRPC

import (
	"context"

	"github.com/DrIhor/test_task/internal/models/items"
	pb "github.com/DrIhor/test_task/pkg/grpc"
)

func (s *Server) AddNewItem(ctx context.Context, in *pb.Item) (*pb.NoneObjectResp, error) {
	item := items.Item{
		Name:        in.Name,
		Price:       in.Price,
		ItemsNumber: in.ItemsNumber,
		Description: in.Description,
	}

	err := s.storage.AddNewItem(item)
	if err != nil {
		return nil, err
	}

	return &pb.NoneObjectResp{
		DoneWork: true,
	}, nil
}

func (s *Server) GetAllItems(ctx context.Context, in *pb.NoneObjectRequest) (*pb.EncodeItemResponse, error) {

	res, err := s.storage.GetAllItems()
	if err != nil {
		return nil, err
	}

	return &pb.EncodeItemResponse{
		Result: res,
	}, nil
}

func (s *Server) GetItem(ctx context.Context, in *pb.ItemName) (*pb.EncodeItemResponse, error) {

	res, err := s.storage.GetItem(in.Name)
	if err != nil {
		return nil, err
	}

	return &pb.EncodeItemResponse{
		Result: res,
	}, nil
}

func (s *Server) DeleteItem(ctx context.Context, in *pb.ItemName) (*pb.NoneObjectResp, error) {

	err := s.storage.DeleteItem(in.Name)
	if err != nil {
		return nil, err
	}

	return &pb.NoneObjectResp{
		DoneWork: true,
	}, nil
}

func (s *Server) UpdateItem(ctx context.Context, in *pb.ItemName) (*pb.EncodeItemResponse, error) {

	res, err := s.storage.UpdateItem(in.Name)
	if err != nil {
		return nil, err
	}

	return &pb.EncodeItemResponse{
		Result: res,
	}, nil
}
