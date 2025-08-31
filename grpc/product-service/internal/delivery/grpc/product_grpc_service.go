package grpcservices

import (
	"context"
	"product-service/internal/domain"
	"product-service/pb"
)

type ProductGrpcService struct {
	pb.UnimplementedProductServiceServer
	productUsecase domain.ProductUsecase
}

func NewProductGrpcService(productUsecase domain.ProductUsecase) *ProductGrpcService {
	return &ProductGrpcService{
		productUsecase: productUsecase,
	}
}

func (s *ProductGrpcService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {

	err := s.productUsecase.CreateProduct(ctx, &domain.Product{
		ID:          "1",
		Name:        "Product 1",
		Description: "Product 1 Description",
		Price:       100,
		Quantity:    100,
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetProductResponse{Product: &pb.Product{Id: "1", Name: "Product 1", Description: "Product 1 Description", Price: 100, Quantity: 100}}, nil
}
