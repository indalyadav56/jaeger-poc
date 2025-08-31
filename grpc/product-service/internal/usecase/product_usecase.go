package usecase

import (
	"common-service/pkg/trace"
	"context"
	"product-service/internal/domain"
)

type productUsecase struct {
	productRepository domain.ProductRepository
}

func NewProductUsecase(productRepository domain.ProductRepository) *productUsecase {
	return &productUsecase{productRepository: productRepository}
}

func (u *productUsecase) CreateProduct(ctx context.Context, product *domain.Product) error {
	ctx, span := trace.StartSpan(ctx, "ProductUsecase.CreateProduct")
	defer span.End()

	return u.productRepository.CreateProduct(ctx, product)
}
