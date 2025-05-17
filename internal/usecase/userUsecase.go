package usecase

import (
	"context"
	"fmt"
	"library-Mongo/internal/domain"
	customErr "library-Mongo/internal/errors"
	"library-Mongo/internal/repo"
	"library-Mongo/internal/usecase/dto"
	"time"
)

type UserUsecase struct {
	userRepo repo.UserRepository
}

func NewUserUsecase(userRepo repo.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

func (uc *UserUsecase) RegisterUser(ctx context.Context, input dto.RegisterUserInput) (domain.User, error) {
	if input.FullName == "" || input.Phone == "" || input.Password == "" || input.Role == "" {
		return domain.User{}, fmt.Errorf("RegisterUser: missing required fields")
	}

	user := domain.User{
		FullName:     input.FullName,
		Phone:        input.Phone,
		Password:     input.Password,
		Role:         input.Role,
		RegisteredAt: time.Now().Format("2006-01-02 15:04:05"),
		IsActive:     true,
	}

	if err := uc.userRepo.Create(ctx, &user); err != nil {
		return domain.User{}, fmt.Errorf("RegisterUser: %w", err)
	}

	return user, nil
}

func (uc *UserUsecase) Login(ctx context.Context, phone, password string) (domain.User, error) {
	if phone == "" || password == "" {
		return domain.User{}, fmt.Errorf("Login: phone and password required")
	}

	user, err := uc.userRepo.Login(ctx, phone, password)
	if err != nil {
		return domain.User{}, fmt.Errorf("Login: %w", err)
	}
	if user == nil {
		return domain.User{}, customErr.ErrUserNotFound
	}
	if !user.IsActive {
		return domain.User{}, customErr.ErrUserBlocked
	}

	return *user, nil
}

func (uc *UserUsecase) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	if id == "" {
		return domain.User{}, customErr.ErrInvalidID
	}

	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("GetUserByID: %w", err)
	}
	if user == nil {
		return domain.User{}, customErr.ErrUserNotFound
	}

	return *user, nil
}

func (uc *UserUsecase) SearchUsers(ctx context.Context, filter domain.UserFilter) ([]domain.User, error) {
	users, err := uc.userRepo.Search(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("SearchUsers: %w", err)
	}
	return users, nil
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, input dto.UpdateUserInput) error {
	if input.ID == "" {
		return customErr.ErrInvalidID
	}

	user, err := uc.userRepo.GetByID(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("UpdateUser: %w", err)
	}
	if user == nil {
		return customErr.ErrUserNotFound
	}

	if input.FullName != nil {
		user.FullName = *input.FullName
	}
	if input.Phone != nil {
		user.Phone = *input.Phone
	}
	if input.Password != nil {
		user.Password = *input.Password
	}
	if input.Role != nil {
		user.Role = *input.Role
	}
	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("UpdateUser: %w", err)
	}

	return nil
}

func (uc *UserUsecase) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return customErr.ErrInvalidID
	}
	if err := uc.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("DeleteUser: %w", err)
	}
	return nil
}

func (uc *UserUsecase) CountUsers(ctx context.Context, filter *domain.UserFilter) (int64, error) {
	return uc.userRepo.Count(ctx)
}

func (uc *UserUsecase) BlockUser(ctx context.Context, id string) error {
	active := false
	return uc.UpdateUser(ctx, dto.UpdateUserInput{ID: id, IsActive: &active})
}

func (uc *UserUsecase) UnblockUser(ctx context.Context, id string) error {
	active := true
	return uc.UpdateUser(ctx, dto.UpdateUserInput{ID: id, IsActive: &active})
}
