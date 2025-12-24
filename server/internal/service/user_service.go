package service

import (
	"context"
	"fmt"

	"github.com/example/go-webdb-template/internal/model"
	"github.com/example/go-webdb-template/internal/repository"
)

// UserService はユーザーのビジネスロジックを担当
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService は新しいUserServiceを作成
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser はユーザーを作成
func (s *UserService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	// バリデーション
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	// ユーザー作成
	user, err := s.userRepo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUser はIDでユーザーを取得
func (s *UserService) GetUser(ctx context.Context, id int64) (*model.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user id: %d", id)
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// ListUsers はユーザー一覧を取得
func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*model.User, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// UpdateUser はユーザーを更新
func (s *UserService) UpdateUser(ctx context.Context, id int64, req *model.UpdateUserRequest) (*model.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user id: %d", id)
	}

	// 更新するフィールドが空の場合はエラー
	if req.Name == "" && req.Email == "" {
		return nil, fmt.Errorf("no fields to update")
	}

	user, err := s.userRepo.Update(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// DeleteUser はユーザーを削除
func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid user id: %d", id)
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
