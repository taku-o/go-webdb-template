package service

import (
	"context"
	"fmt"

	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/repository"
)

// DmUserService はユーザーのビジネスロジックを担当
type DmUserService struct {
	dmUserRepo repository.DmUserRepositoryInterface
}

// NewDmUserService は新しいDmUserServiceを作成
func NewDmUserService(dmUserRepo repository.DmUserRepositoryInterface) *DmUserService {
	return &DmUserService{
		dmUserRepo: dmUserRepo,
	}
}

// CreateDmUser はユーザーを作成
func (s *DmUserService) CreateDmUser(ctx context.Context, req *model.CreateDmUserRequest) (*model.DmUser, error) {
	// バリデーション
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	// ユーザー作成
	dmUser, err := s.dmUserRepo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return dmUser, nil
}

// GetDmUser はIDでユーザーを取得
func (s *DmUserService) GetDmUser(ctx context.Context, id string) (*model.DmUser, error) {
	if id == "" {
		return nil, fmt.Errorf("user id is required")
	}

	dmUser, err := s.dmUserRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return dmUser, nil
}

// ListDmUsers はユーザー一覧を取得
func (s *DmUserService) ListDmUsers(ctx context.Context, limit, offset int) ([]*model.DmUser, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	dmUsers, err := s.dmUserRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return dmUsers, nil
}

// UpdateDmUser はユーザーを更新
func (s *DmUserService) UpdateDmUser(ctx context.Context, id string, req *model.UpdateDmUserRequest) (*model.DmUser, error) {
	if id == "" {
		return nil, fmt.Errorf("user id is required")
	}

	// 更新するフィールドが空の場合はエラー
	if req.Name == "" && req.Email == "" {
		return nil, fmt.Errorf("no fields to update")
	}

	dmUser, err := s.dmUserRepo.Update(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return dmUser, nil
}

// DeleteDmUser はユーザーを削除
func (s *DmUserService) DeleteDmUser(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("user id is required")
	}

	if err := s.dmUserRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
