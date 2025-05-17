package mongo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"library-Mongo/internal/domain"
	customErr "library-Mongo/internal/errors"
	"time"
)

type BorrowRepoMongo struct {
	col *mongo.Collection
}

func NewBorrowRepo(db *mongo.Database) *BorrowRepoMongo {
	return &BorrowRepoMongo{
		col: db.Collection("borrows"),
	}
}

func (r *BorrowRepoMongo) Create(ctx context.Context, b *domain.Borrow) error {
	filter := bson.M{
		"bookId":     b.BookID,
		"returnedAt": bson.M{"$exists": false},
	}

	count, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("BorrowRepoMongo.Create (count check): %w", err)
	}

	if count > 0 {
		return customErr.ErrBookAlreadyBorrowed
	}

	doc := bson.M{
		"clientId":   b.ClientID,
		"bookId":     b.BookID,
		"borrowedAt": b.BorrowedAt,
	}
	if b.ReturnedAt != nil {
		doc["returnedAt"] = b.ReturnedAt
	}

	res, err := r.col.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("BorrowRepoMongo.Create (insert): %w", err)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("BorrowRepoMongo.Create (cast ID): inserted ID is not ObjectID")
	}
	b.ID = oid.Hex()

	return nil
}

func (r *BorrowRepoMongo) Close(ctx context.Context, borrowID string, returnTime time.Time) error {
	objID, err := primitive.ObjectIDFromHex(borrowID)
	if err != nil {
		return fmt.Errorf("BorrowRepoMongo.Close (parse ID): %w", err)
	}
	update := bson.M{
		"$set": bson.M{
			"returnedAt": returnTime,
		},
	}
	_, err = r.col.UpdateByID(ctx, objID, update)
	if err != nil {
		return fmt.Errorf("BorrowRepoMongo.Close (update): %w", err)
	}
	return nil
}

func (r *BorrowRepoMongo) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Borrow, error) {
	var b domain.Borrow
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&b)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("BorrowRepoMongo.GetByID: %w", err)
	}
	b.ID = id.Hex()
	return &b, nil
}

// Отчет №1 (Вернуть список выдачи книг конкретному пользователю)
func (r *BorrowRepoMongo) GetByClientID(ctx context.Context, clientID primitive.ObjectID) ([]domain.Borrow, error) {
	cursor, err := r.col.Find(ctx, bson.M{"clientId": clientID})
	if err != nil {
		return nil, fmt.Errorf("BorrowRepoMongo.GetByClientID (find): %w", err)
	}
	defer cursor.Close(ctx)

	var borrows []domain.Borrow
	if err := cursor.All(ctx, &borrows); err != nil {
		return nil, fmt.Errorf("BorrowRepoMongo.GetByClientID (decode): %w", err)
	}
	return borrows, nil
}

// Отчет №2 (Вернуть список просроченных книг)
func (r *BorrowRepoMongo) GetOverdue(ctx context.Context, now time.Time) ([]domain.Borrow, error) {
	limit := now.AddDate(0, 0, -21)
	filter := bson.M{
		"returnedAt": bson.M{"$eq": nil},
		"borrowedAt": bson.M{"$lt": limit},
	}

	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("BorrowRepoMongo.GetOverdue (find): %w", err)
	}
	defer cursor.Close(ctx)

	var results []domain.Borrow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("BorrowRepoMongo.GetOverdue (decode): %w", err)
	}
	return results, nil
}

// Отчет №3 (Вернуть кол-во пришедших читателей по дням за период)
func (r *BorrowRepoMongo) GetDailyStats(ctx context.Context, from, to time.Time) ([]domain.BorrowStat, error) {
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{
			"borrowedAt": bson.M{
				"$gte": from,
				"$lte": to,
			},
		}}},
		{{"$group", bson.M{
			"_id": bson.M{
				"$dateToString": bson.M{
					"format": "%Y-%m-%d",
					"date":   "$borrowedAt",
				},
			},
			"clients": bson.M{"$addToSet": "$clientId"},
		}}},
		{{"$project", bson.M{
			"date":          "$_id",
			"uniqueReaders": bson.M{"$size": "$clients"},
			"_id":           0,
		}}},
		{{"$sort", bson.M{"date": 1}}},
	}

	cursor, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("BorrowRepoMongo.GetDailyStats (aggregate): %w", err)
	}
	defer cursor.Close(ctx)

	var stats []domain.BorrowStat
	if err := cursor.All(ctx, &stats); err != nil {
		return nil, fmt.Errorf("BorrowRepoMongo.GetDailyStats (decode): %w", err)
	}
	return stats, nil
}

func (r *BorrowRepoMongo) CountActive(ctx context.Context) (int64, error) {
	count, err := r.col.CountDocuments(ctx, bson.M{"returnedAt": nil})
	if err != nil {
		return 0, fmt.Errorf("BorrowRepoMongo.CountActive: %w", err)
	}
	return count, nil
}

func (r *BorrowRepoMongo) HasActiveBorrow(ctx context.Context, bookID primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"bookId":     bookID,
		"returnedAt": bson.M{"$exists": false},
	}
	count, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("BorrowRepoMongo.IsBookCurrentlyBorrowed: %w", err)
	}
	return count > 0, nil
}
