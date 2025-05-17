package usecase

import (
	"context"
	"fmt"
	"library-Mongo/internal/domain"
	"library-Mongo/internal/repo"
	"library-Mongo/internal/usecase/dto"
)

type BookUsecase struct {
	bookRepo repo.BookRepository
}

func NewBookUsecase(bookRepo repo.BookRepository) *BookUsecase {
	return &BookUsecase{bookRepo: bookRepo}
}

func (uc *BookUsecase) CreateBook(ctx context.Context, input dto.CreateBookInput) (domain.Book, error) {
	if input.Title == "" || input.Author == "" || input.Genre == "" {
		return domain.Book{}, fmt.Errorf("CreateBook: missing required fields")
	}

	book := domain.Book{
		Title:  input.Title,
		Author: input.Author,
		Year:   input.Year,
		Genre:  input.Genre,
	}

	if err := uc.bookRepo.Create(ctx, &book); err != nil {
		return domain.Book{}, fmt.Errorf("CreateBook: %w", err)
	}

	return book, nil
}

func (uc *BookUsecase) UpdateBook(ctx context.Context, input dto.UpdateBookInput) error {
	if input.ID == "" {
		return fmt.Errorf("UpdateBook: missing ID")
	}

	// Получить текущую книгу из репо
	existing, err := uc.bookRepo.GetByID(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("UpdateBook: failed to load existing book: %w", err)
	}

	// Обновить только те поля, которые переданы
	if input.Title != nil {
		existing.Title = *input.Title
	}
	if input.Author != nil {
		existing.Author = *input.Author
	}
	if input.Year != nil {
		existing.Year = *input.Year
	}
	if input.Genre != nil {
		existing.Genre = *input.Genre
	}

	// Сохранить изменения
	if err := uc.bookRepo.Update(ctx, existing); err != nil {
		return fmt.Errorf("UpdateBook: %w", err)
	}

	return nil
}

func (uc *BookUsecase) DeleteBook(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("DeleteBook: missing ID")
	}

	if err := uc.bookRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("DeleteBook: %w", err)
	}

	return nil
}

func (uc *BookUsecase) GetBookByID(ctx context.Context, id string) (domain.Book, error) {
	if id == "" {
		return domain.Book{}, fmt.Errorf("GetBookByID: missing ID")
	}

	book, err := uc.bookRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Book{}, fmt.Errorf("GetBookByID: %w", err)
	}

	return *book, nil
}

func (uc *BookUsecase) SearchBooks(ctx context.Context, filter domain.BookFilter) ([]domain.Book, error) {
	books, err := uc.bookRepo.Search(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("SearchBooks: %w", err)
	}
	return books, nil
}

func (uc *BookUsecase) CountBooks(ctx context.Context) (int64, error) {
	count, err := uc.bookRepo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("CountBooks: %w", err)
	}
	return count, nil
}
