package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateIndexes(db *mongo.Database) error {
	ctx := context.TODO()

	_, err := db.Collection("users").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "fullName", Value: 1}}},
		{Keys: bson.D{{Key: "phone", Value: 1}}},
	})
	if err != nil {
		return err
	}

	_, err = db.Collection("books").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "title", Value: 1}}},
		{Keys: bson.D{{Key: "author", Value: 1}}},
		{Keys: bson.D{
			{Key: "author", Value: 1},
			{Key: "title", Value: 1},
		}},
	})
	if err != nil {
		return err
	}

	_, err = db.Collection("borrows").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{
			{Key: "returnedAt", Value: 1},
			{Key: "borrowedAt", Value: 1},
		}},
		{Keys: bson.D{
			{Key: "clientId", Value: 1},
			{Key: "borrowedAt", Value: -1},
		}},
		{Keys: bson.D{
			{Key: "borrowedAt", Value: 1},
			{Key: "clientId", Value: 1},
		}},
	})
	if err != nil {
		return err
	}

	return nil
}
