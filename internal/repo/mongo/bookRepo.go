package mongo

import (
	"context"
	"errors"
	"fmt"
	"library-Mongo/internal/domain"
	customErr "library-Mongo/internal/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookRepoMongo struct {
	col *mongo.Collection
}

func NewBookRepo(db *mongo.Database) *BookRepoMongo {
	return &BookRepoMongo{
		col: db.Collection("books"),
	}
}

func (r *BookRepoMongo) Create(ctx context.Context, b *domain.Book) error {
	bookDoc := struct {
		Title  string `bson:"title"`
		Author string `bson:"author"`
		Year   int    `bson:"year"`
		Genre  string `bson:"genre"`
	}{
		Title:  b.Title,
		Author: b.Author,
		Year:   b.Year,
		Genre:  b.Genre,
	}

	res, err := r.col.InsertOne(ctx, bookDoc)
	if err != nil {
		return fmt.Errorf("BookRepoMongo.Create: %w", err)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("BookRepoMongo.Create: inserted ID is not ObjectID")
	}
	b.ID = oid.Hex()

	return nil
}

func (r *BookRepoMongo) Update(ctx context.Context, b *domain.Book) error {
	objID, err := primitive.ObjectIDFromHex(b.ID)
	if err != nil {
		return fmt.Errorf("BookRepoMongo.Update: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"title":  b.Title,
			"author": b.Author,
			"year":   b.Year,
			"genre":  b.Genre,
		},
	}

	_, err = r.col.UpdateByID(ctx, objID, update)
	if err != nil {
		return fmt.Errorf("BookRepoMongo.Update: %w", err)
	}
	return nil
}

func (r *BookRepoMongo) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("BookRepoMongo.Delete: %w", err)
	}

	_, err = r.col.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("BookRepoMongo.Delete: %w", err)
	}
	return nil
}

func (r *BookRepoMongo) GetByID(ctx context.Context, id string) (*domain.Book, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("BookRepoMongo.GetByID: %w", err)
	}

	var doc domain.Book
	err = r.col.FindOne(ctx, bson.M{"_id": objID}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("BookRepoMongo.GetByID: %w", customErr.ErrBookNotFound)
		}
		return nil, fmt.Errorf("BookRepoMongo.GetByID: %w", err)
	}

	doc.ID = objID.Hex()
	return &doc, nil
}

func (r *BookRepoMongo) Search(ctx context.Context, filter domain.BookFilter) ([]domain.Book, error) {
	query := bson.M{}

	if filter.Title != "" {
		query["title"] = bson.M{"$regex": filter.Title, "$options": "i"}
	}
	if filter.Author != "" {
		query["author"] = bson.M{"$regex": filter.Author, "$options": "i"}
	}
	if len(filter.Genres) > 0 {
		query["genre"] = bson.M{"$in": filter.Genres}
	}

	cursor, err := r.col.Find(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("BookRepoMongo.Search: %w", err)
	}
	defer cursor.Close(ctx)

	var books []domain.Book
	for cursor.Next(ctx) {
		var book domain.Book
		if err := cursor.Decode(&book); err != nil {
			return nil, fmt.Errorf("BookRepoMongo.Search (decode): %w", err)
		}
		books = append(books, book)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("BookRepoMongo.Search (cursor): %w", err)
	}

	return books, nil
}

func (r *BookRepoMongo) Count(ctx context.Context) (int64, error) {
	count, err := r.col.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("BookRepoMongo.Count: %w", err)
	}
	return count, nil
}
