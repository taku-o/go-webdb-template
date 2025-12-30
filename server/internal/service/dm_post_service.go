package service

import (
	"context"
	"fmt"

	"github.com/taku-o/go-webdb-template/internal/model"
	"github.com/taku-o/go-webdb-template/internal/repository"
)

// DmPostService は投稿のビジネスロジックを担当
type DmPostService struct {
	dmPostRepo repository.DmPostRepositoryInterface
	dmUserRepo repository.DmUserRepositoryInterface
}

// NewDmPostService は新しいDmPostServiceを作成
func NewDmPostService(dmPostRepo repository.DmPostRepositoryInterface, dmUserRepo repository.DmUserRepositoryInterface) *DmPostService {
	return &DmPostService{
		dmPostRepo: dmPostRepo,
		dmUserRepo: dmUserRepo,
	}
}

// CreateDmPost は投稿を作成
func (s *DmPostService) CreateDmPost(ctx context.Context, req *model.CreateDmPostRequest) (*model.DmPost, error) {
	// バリデーション
	if req.UserID == "" {
		return nil, fmt.Errorf("user id is required")
	}
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if req.Content == "" {
		return nil, fmt.Errorf("content is required")
	}

	// ユーザーの存在確認
	_, err := s.dmUserRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 投稿作成
	dmPost, err := s.dmPostRepo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return dmPost, nil
}

// GetDmPost はIDで投稿を取得
func (s *DmPostService) GetDmPost(ctx context.Context, id string, userID string) (*model.DmPost, error) {
	if id == "" {
		return nil, fmt.Errorf("post id is required")
	}
	if userID == "" {
		return nil, fmt.Errorf("user id is required")
	}

	dmPost, err := s.dmPostRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return dmPost, nil
}

// ListDmPosts は投稿一覧を取得
func (s *DmPostService) ListDmPosts(ctx context.Context, limit, offset int) ([]*model.DmPost, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	dmPosts, err := s.dmPostRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}

	return dmPosts, nil
}

// ListDmPostsByUser はユーザーIDで投稿一覧を取得
func (s *DmPostService) ListDmPostsByUser(ctx context.Context, userID string, limit, offset int) ([]*model.DmPost, error) {
	if userID == "" {
		return nil, fmt.Errorf("user id is required")
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	dmPosts, err := s.dmPostRepo.ListByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts by user: %w", err)
	}

	return dmPosts, nil
}

// GetDmUserPosts はユーザーと投稿をJOINして取得
func (s *DmPostService) GetDmUserPosts(ctx context.Context, limit, offset int) ([]*model.DmUserPost, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	dmUserPosts, err := s.dmPostRepo.GetUserPosts(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user posts: %w", err)
	}

	return dmUserPosts, nil
}

// UpdateDmPost は投稿を更新
func (s *DmPostService) UpdateDmPost(ctx context.Context, id string, userID string, req *model.UpdateDmPostRequest) (*model.DmPost, error) {
	if id == "" {
		return nil, fmt.Errorf("post id is required")
	}
	if userID == "" {
		return nil, fmt.Errorf("user id is required")
	}

	// 更新するフィールドが空の場合はエラー
	if req.Title == "" && req.Content == "" {
		return nil, fmt.Errorf("no fields to update")
	}

	dmPost, err := s.dmPostRepo.Update(ctx, id, userID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return dmPost, nil
}

// DeleteDmPost は投稿を削除
func (s *DmPostService) DeleteDmPost(ctx context.Context, id string, userID string) error {
	if id == "" {
		return fmt.Errorf("post id is required")
	}
	if userID == "" {
		return fmt.Errorf("user id is required")
	}

	if err := s.dmPostRepo.Delete(ctx, id, userID); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}
