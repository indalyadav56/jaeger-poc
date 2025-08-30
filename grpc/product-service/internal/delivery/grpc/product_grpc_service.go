package grpcservices

import (
	"context"
	"product-service/pb"
)

type ProductGrpcService struct {
	pb.UnimplementedProductServiceServer
}

func NewProductGrpcService() *ProductGrpcService {
	return &ProductGrpcService{}
}

func (s *ProductGrpcService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	return &pb.GetProductResponse{Product: &pb.Product{Id: req.Id, Name: "Product 1", Description: "Product 1 Description", Price: 100, Quantity: 100}}, nil
}
