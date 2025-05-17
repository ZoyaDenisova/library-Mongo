package usecase

import (
	"context"
	"library-Mongo/internal/domain"
	"library-Mongo/internal/usecase/dto"
	"time"
)

type BookUC interface {
	CreateBook(ctx context.Context, input dto.CreateBookInput) (domain.Book, error)
	UpdateBook(ctx context.Context, input dto.UpdateBookInput) error
	DeleteBook(ctx context.Context, id string) error
	GetBookByID(ctx context.Context, id string) (domain.Book, error)
	SearchBooks(ctx context.Context, filter domain.BookFilter) ([]domain.Book, error)
	CountBooks(ctx context.Context) (int64, error)
}

type UserUC interface {
	RegisterUser(ctx context.Context, input dto.RegisterUserInput) (domain.User, error)
	Login(ctx context.Context, phone, password string) (domain.User, error)
	GetUserByID(ctx context.Context, id string) (domain.User, error)
	SearchUsers(ctx context.Context, filter domain.UserFilter) ([]domain.User, error)
	UpdateUser(ctx context.Context, input dto.UpdateUserInput) error
	DeleteUser(ctx context.Context, id string) error
	CountUsers(ctx context.Context, filter *domain.UserFilter) (int64, error)
	BlockUser(ctx context.Context, id string) error
	UnblockUser(ctx context.Context, id string) error
}

type BorrowUC interface {
	// Оформить выдачу книги (librarian)
	BorrowBook(ctx context.Context, input dto.BorrowBookInput) (domain.Borrow, error)
	// Оформить возврат книги (librarian)
	ReturnBook(ctx context.Context, input dto.ReturnBookInput) error
	// История всех выдач конкретного читателя (reader/librarian)
	GetBorrowHistory(ctx context.Context, userID string) (dto.BorrowHistoryResponse, error)
	//Список всех просроченных выдач (librarian)
	GetOverdueBorrows(ctx context.Context) ([]dto.OverdueReportItem, error)
	// Статистика уникальных читателей по дням/месяцам (для отчёта 3)
	GetDailyBorrowStats(ctx context.Context, from, to time.Time) ([]domain.BorrowStat, error)
	// Подсчитать число активных (не возвращённых) выдач
	CountActiveBorrows(ctx context.Context) (int64, error)
}
