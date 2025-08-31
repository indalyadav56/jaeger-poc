package repository

import (
	"common-service/pkg/db/mongodb"
	"common-service/pkg/trace"
	"context"
	"product-service/internal/domain"
)

type mongodbProductRepository struct {
	mongodbClient *mongodb.MongoClient
}

func NewMongodbProductRepository(mongodbClient *mongodb.MongoClient) *mongodbProductRepository {
	return &mongodbProductRepository{
		mongodbClient: mongodbClient,
	}
}

func (r *mongodbProductRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
	ctx, span := trace.StartSpan(ctx, "MongodbProductRepository.CreateProduct")
	defer span.End()

	collection := r.mongodbClient.DB.Collection("products")

	_, err := collection.InsertOne(ctx, product)
	if err != nil {
		return err
	}

	return nil
}
