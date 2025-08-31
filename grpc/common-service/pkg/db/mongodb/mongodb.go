package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

type Config struct {
	URI            string        `json:"uri" yaml:"uri"`
	Database       string        `json:"database" yaml:"database"`
	MaxPoolSize    int           `json:"max_pool_size" yaml:"max_pool_size"`
	ConnectTimeout time.Duration `json:"connect_timeout" yaml:"connect_timeout"`
}

type MongoClient struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewMongoDBClient(cfg Config) (*MongoClient, error) {
	// Set up client options
	opts := options.Client().
		ApplyURI(cfg.URI).
		SetMaxPoolSize(uint64(cfg.MaxPoolSize)).
		SetConnectTimeout(cfg.ConnectTimeout).
		SetMonitor(otelmongo.NewMonitor())

	// Create client
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Verify connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return &MongoClient{
		Client: client,
		DB:     client.Database(cfg.Database),
	}, nil
}

// Close terminates the MongoDB connection
func (c *MongoClient) Close() error {
	if c.Client == nil {
		return nil
	}
	return c.Client.Disconnect(context.TODO())
}

// Ping checks the health of the MongoDB connection
func (c *MongoClient) Ping() error {
	return c.Client.Ping(context.TODO(), readpref.Primary())
}

func (c *MongoClient) CreateCollection(name string) error {
	if c.DB == nil {
		return fmt.Errorf("database is not initialized")
	}
	if err := c.DB.CreateCollection(context.TODO(), name); err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}
	return nil
}
