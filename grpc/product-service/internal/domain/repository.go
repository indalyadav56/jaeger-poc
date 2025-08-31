package domain

import "context"

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *Product) error
}
