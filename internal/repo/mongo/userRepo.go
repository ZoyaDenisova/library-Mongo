package mongo

import (
	"context"
	"errors"
	"library-Mongo/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepoMongo struct {
	col *mongo.Collection
}

func NewUserRepo(db *mongo.Database) *UserRepoMongo {
	return &UserRepoMongo{
		col: db.Collection("users"),
	}
}

func (r *UserRepoMongo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user domain.User
	err = r.col.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	user.ID = objID.Hex()
	return &user, err
}

func (r *UserRepoMongo) Login(ctx context.Context, phone, password string) (*domain.User, error) {
	var user domain.User
	err := r.col.FindOne(ctx, bson.M{
		"phone":    phone,
		"password": password,
	}).Decode(&user)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if oid, ok := userIDFromUser(user); ok {
		user.ID = oid.Hex()
	}
	return &user, err
}

func (r *UserRepoMongo) Search(ctx context.Context, filter domain.UserFilter) ([]domain.User, error) {
	query := bson.M{}

	if filter.FullNameContains != "" || filter.Phone != "" || filter.Role != "" {
		var orConditions []bson.M

		if filter.FullNameContains != "" {
			orConditions = append(orConditions, bson.M{"fullName": bson.M{"$regex": filter.FullNameContains, "$options": "i"}})
		}
		if filter.Phone != "" {
			orConditions = append(orConditions, bson.M{"phone": bson.M{"$regex": filter.Phone, "$options": "i"}})
		}
		if filter.Role != "" {
			orConditions = append(orConditions, bson.M{"role": bson.M{"$regex": filter.Role, "$options": "i"}})
		}

		if len(orConditions) > 0 {
			query["$or"] = orConditions
		}

		// отдельно оставляем фильтрацию по активности
		if filter.OnlyActive != nil {
			query["isActive"] = *filter.OnlyActive
		}
	}

	cursor, err := r.col.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []domain.User
	for cursor.Next(ctx) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		if oid, ok := userIDFromUser(user); ok {
			user.ID = oid.Hex()
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepoMongo) Create(ctx context.Context, u *domain.User) error {
	// Вставляем без ID — Mongo сам создаст _id
	doc := bson.M{
		"fullName":     u.FullName,
		"phone":        u.Phone,
		"password":     u.Password,
		"role":         u.Role,
		"registeredAt": u.RegisteredAt,
		"isActive":     u.IsActive,
	}

	res, err := r.col.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		u.ID = oid.Hex()
	}
	return nil
}

func (r *UserRepoMongo) Update(ctx context.Context, u *domain.User) error {
	objID, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"fullName": u.FullName,
			"phone":    u.Phone,
			"password": u.Password,
			"role":     u.Role,
			"isActive": u.IsActive,
		},
	}
	_, err = r.col.UpdateByID(ctx, objID, update)
	return err
}

func (r *UserRepoMongo) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.col.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (r *UserRepoMongo) Count(ctx context.Context) (int64, error) {
	return r.col.CountDocuments(ctx, bson.M{})
}

// Вспомогательная функция для попытки извлечь ObjectID
func userIDFromUser(u domain.User) (primitive.ObjectID, bool) {
	id, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		return primitive.NilObjectID, false
	}
	return id, true
}
