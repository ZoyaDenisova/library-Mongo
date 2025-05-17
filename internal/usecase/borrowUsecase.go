package usecase

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"library-Mongo/internal/domain"
	customErr "library-Mongo/internal/errors"
	"library-Mongo/internal/repo"
	"library-Mongo/internal/usecase/dto"
	"sort"
	"time"
)

type BorrowUsecase struct {
	borrowRepo repo.BorrowRepository
	bookRepo   repo.BookRepository
	userRepo   repo.UserRepository
}

func NewBorrowUsecase(
	borrowRepo repo.BorrowRepository,
	bookRepo repo.BookRepository,
	userRepo repo.UserRepository,
) *BorrowUsecase {
	return &BorrowUsecase{
		borrowRepo: borrowRepo,
		bookRepo:   bookRepo,
		userRepo:   userRepo,
	}
}

func (uc *BorrowUsecase) GetBorrowHistory(ctx context.Context, userID string) (dto.BorrowHistoryResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return dto.BorrowHistoryResponse{}, fmt.Errorf("GetBorrowHistory: get user: %w", err)
	}
	if user == nil {
		return dto.BorrowHistoryResponse{}, customErr.ErrUserNotFound
	}

	objID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return dto.BorrowHistoryResponse{}, customErr.ErrInvalidID
	}

	borrows, err := uc.borrowRepo.GetByClientID(ctx, objID)
	if err != nil {
		return dto.BorrowHistoryResponse{}, fmt.Errorf("GetBorrowHistory: get borrows: %w", err)
	}

	now := time.Now()
	var overdue, normal []dto.BorrowHistoryItem

	for _, b := range borrows {
		book, err := uc.bookRepo.GetByID(ctx, b.BookID.Hex())
		if err != nil || book == nil {
			continue // можно логировать
		}

		isOverdue := b.ReturnedAt == nil && b.BorrowedAt.Before(now.AddDate(0, 0, -21))
		item := dto.BorrowHistoryItem{
			BorrowID:   b.ID,
			BookID:     book.ID,
			Title:      book.Title,
			Author:     book.Author,
			BorrowedAt: b.BorrowedAt,
			ReturnedAt: b.ReturnedAt,
			Status:     "ok",
		}
		if isOverdue {
			item.Status = "overdue"
			overdue = append(overdue, item)
		} else {
			normal = append(normal, item)
		}
	}

	sort.SliceStable(overdue, func(i, j int) bool {
		return overdue[i].BorrowedAt.Before(overdue[j].BorrowedAt)
	})
	sort.SliceStable(normal, func(i, j int) bool {
		return normal[i].BorrowedAt.Before(normal[j].BorrowedAt)
	})

	return dto.BorrowHistoryResponse{
		UserID:   user.ID,
		FullName: user.FullName,
		Phone:    user.Phone,
		History:  append(overdue, normal...),
	}, nil
}

func (uc *BorrowUsecase) BorrowBook(ctx context.Context, input dto.BorrowBookInput) (domain.Borrow, error) {
	// 1. Проверка валидности ID
	userObjID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return domain.Borrow{}, customErr.ErrInvalidID
	}
	bookObjID, err := primitive.ObjectIDFromHex(input.BookID)
	if err != nil {
		return domain.Borrow{}, customErr.ErrInvalidID
	}

	// 2. Проверка, существует ли пользователь
	user, err := uc.userRepo.GetByID(ctx, input.UserID)
	if err != nil {
		return domain.Borrow{}, fmt.Errorf("BorrowBook: get user: %w", err)
	}
	if user == nil {
		return domain.Borrow{}, customErr.ErrUserNotFound
	}

	// 3. Проверка, существует ли книга
	book, err := uc.bookRepo.GetByID(ctx, input.BookID)
	if err != nil {
		return domain.Borrow{}, fmt.Errorf("BorrowBook: get book: %w", err)
	}
	if book == nil {
		return domain.Borrow{}, customErr.ErrBookNotFound
	}

	// 4. Проверка: книга уже выдана?
	hasActive, err := uc.borrowRepo.HasActiveBorrow(ctx, bookObjID)
	if err != nil {
		return domain.Borrow{}, fmt.Errorf("BorrowBook: check active borrow: %w", err)
	}
	if hasActive {
		return domain.Borrow{}, customErr.ErrBookAlreadyBorrowed
	}

	// 5. Сохраняем новую выдачу
	borrow := domain.Borrow{
		ClientID:   userObjID,
		BookID:     bookObjID,
		BorrowedAt: time.Now(),
	}

	if err := uc.borrowRepo.Create(ctx, &borrow); err != nil {
		return domain.Borrow{}, fmt.Errorf("BorrowBook: insert: %w", err)
	}

	return borrow, nil
}

func (uc *BorrowUsecase) ReturnBook(ctx context.Context, input dto.ReturnBookInput) error {
	if input.BorrowID == "" {
		return customErr.ErrInvalidID
	}

	// Проверка валидности ID
	objID, err := primitive.ObjectIDFromHex(input.BorrowID)
	if err != nil {
		return customErr.ErrInvalidID
	}

	// Загружаем все выдачи этого пользователя (мог бы быть отдельный метод GetByBorrowID, но допустимо и так, если у тебя нет)
	// Но лучше будет, если у тебя есть такой метод (например, GetBorrowByID)

	// Так что делаем аккуратно: надо добавить метод repo.GetBorrowByID
	borrow, err := uc.borrowRepo.GetByID(ctx, objID)
	if err != nil {
		return fmt.Errorf("ReturnBook: fetch borrow: %w", err)
	}
	if borrow == nil {
		return customErr.ErrBorrowNotFound
	}
	if borrow.ReturnedAt != nil {
		return customErr.ErrAlreadyReturned
	}

	// Помечаем как возвращённую
	now := time.Now()
	err = uc.borrowRepo.Close(ctx, input.BorrowID, now)
	if err != nil {
		return fmt.Errorf("ReturnBook: close borrow: %w", err)
	}

	return nil
}

func (uc *BorrowUsecase) GetOverdueBorrows(ctx context.Context) ([]dto.OverdueReportItem, error) {
	now := time.Now()

	// Получаем список всех просроченных выдач
	borrows, err := uc.borrowRepo.GetOverdue(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("GetOverdueBorrows: %w", err)
	}

	// Счётчик просрочек по пользователям
	overdueCount := make(map[string]int)
	for _, b := range borrows {
		overdueCount[b.ClientID.Hex()]++
	}

	var report []dto.OverdueReportItem

	for _, b := range borrows {
		userID := b.ClientID.Hex()
		bookID := b.BookID.Hex()

		// Получаем пользователя
		user, err := uc.userRepo.GetByID(ctx, userID)
		if err != nil || user == nil {
			continue // можно логировать
		}

		// Получаем книгу
		book, err := uc.bookRepo.GetByID(ctx, bookID)
		if err != nil || book == nil {
			continue // можно логировать
		}

		// Вычисляем просрочку
		daysOverdue := int(now.Sub(b.BorrowedAt).Hours()/24) - 21
		if daysOverdue < 0 {
			daysOverdue = 0 // на всякий случай
		}

		report = append(report, dto.OverdueReportItem{
			UserID:       user.ID,
			FullName:     user.FullName,
			Phone:        user.Phone,
			BookID:       book.ID,
			Title:        book.Title,
			Author:       book.Author,
			BorrowedAt:   b.BorrowedAt,
			DaysOverdue:  daysOverdue,
			TotalOverdue: overdueCount[userID],
		})
	}

	return report, nil
}

func (uc *BorrowUsecase) GetDailyBorrowStats(ctx context.Context, from, to time.Time) ([]domain.BorrowStat, error) {
	if from.After(to) {
		return nil, fmt.Errorf("GetDailyBorrowStats: invalid time range (from > to)")
	}

	stats, err := uc.borrowRepo.GetDailyStats(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("GetDailyBorrowStats: %w", err)
	}

	return stats, nil
}

func (uc *BorrowUsecase) CountActiveBorrows(ctx context.Context) (int64, error) {
	count, err := uc.borrowRepo.CountActive(ctx)
	if err != nil {
		return 0, fmt.Errorf("CountActiveBorrows: %w", err)
	}
	return count, nil
}
