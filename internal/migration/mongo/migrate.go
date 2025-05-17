// go run -tags migrate ./cmd/app

package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func Migrate(db *mongo.Database) error {
	return CreateIndexes(db)
}
