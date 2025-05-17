package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"library-Mongo/internal/config"
)

func Connect(ctx context.Context, cfg *config.Config) (*mongo.Database, error) {
	clientOpts := options.Client().ApplyURI(cfg.MongoURI)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB:", cfg.MongoURI)
	return client.Database(cfg.Database), nil
}
