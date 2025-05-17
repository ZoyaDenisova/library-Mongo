package repo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"library-Mongo/internal/domain"
	"time"
)

type (
	BookRepository interface {
		Create(ctx context.Context, b *domain.Book) error
		Update(ctx context.Context, b *domain.Book) error
		Delete(ctx context.Context, id string) error
		GetByID(ctx context.Context, id string) (*domain.Book, error)
		Search(ctx context.Context, filter domain.BookFilter) ([]domain.Book, error)
		Count(ctx context.Context) (int64, error)
	}

	UserRepository interface {
		GetByID(ctx context.Context, id string) (*domain.User, error)
		Login(ctx context.Context, phone, password string) (*domain.User, error)
		Search(ctx context.Context, filter domain.UserFilter) ([]domain.User, error)
		Create(ctx context.Context, u *domain.User) error
		Update(ctx context.Context, u *domain.User) error
		Delete(ctx context.Context, id string) error
		Count(ctx context.Context) (int64, error)
	}

	BorrowRepository interface {
		Create(ctx context.Context, b *domain.Borrow) error
		Close(ctx context.Context, borrowID string, returnTime time.Time) error
		GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Borrow, error)
		GetByClientID(ctx context.Context, clientID primitive.ObjectID) ([]domain.Borrow, error)
		GetOverdue(ctx context.Context, now time.Time) ([]domain.Borrow, error)
		GetDailyStats(ctx context.Context, from, to time.Time) ([]domain.BorrowStat, error)
		CountActive(ctx context.Context) (int64, error)
		HasActiveBorrow(ctx context.Context, bookID primitive.ObjectID) (bool, error)
	}
)
