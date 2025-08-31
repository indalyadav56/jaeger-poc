package domain

import "context"

type ProductUsecase interface {
	CreateProduct(ctx context.Context, product *Product) error
}
