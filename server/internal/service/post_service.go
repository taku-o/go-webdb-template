package service

import (
	"context"
	"fmt"

	"github.com/example/go-webdb-template/internal/model"
	"github.com/example/go-webdb-template/internal/repository"
)

// PostService は投稿のビジネスロジックを担当
type PostService struct {
	postRepo repository.PostRepositoryInterface
	userRepo repository.UserRepositoryInterface
}

// NewPostService は新しいPostServiceを作成
func NewPostService(postRepo repository.PostRepositoryInterface, userRepo repository.UserRepositoryInterface) *PostService {
	return &PostService{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

// CreatePost は投稿を作成
func (s *PostService) CreatePost(ctx context.Context, req *model.CreatePostRequest) (*model.Post, error) {
	// バリデーション
	if req.UserID <= 0 {
		return nil, fmt.Errorf("invalid user id: %d", req.UserID)
	}
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if req.Content == "" {
		return nil, fmt.Errorf("content is required")
	}

	// ユーザーの存在確認
	_, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 投稿作成
	post, err := s.postRepo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

// GetPost はIDで投稿を取得
func (s *PostService) GetPost(ctx context.Context, id int64, userID int64) (*model.Post, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid post id: %d", id)
	}
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user id: %d", userID)
	}

	post, err := s.postRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return post, nil
}

// ListPosts は投稿一覧を取得
func (s *PostService) ListPosts(ctx context.Context, limit, offset int) ([]*model.Post, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	posts, err := s.postRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}

	return posts, nil
}

// ListPostsByUser はユーザーIDで投稿一覧を取得
func (s *PostService) ListPostsByUser(ctx context.Context, userID int64, limit, offset int) ([]*model.Post, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user id: %d", userID)
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

	posts, err := s.postRepo.ListByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts by user: %w", err)
	}

	return posts, nil
}

// GetUserPosts はユーザーと投稿をJOINして取得
func (s *PostService) GetUserPosts(ctx context.Context, limit, offset int) ([]*model.UserPost, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	userPosts, err := s.postRepo.GetUserPosts(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user posts: %w", err)
	}

	return userPosts, nil
}

// UpdatePost は投稿を更新
func (s *PostService) UpdatePost(ctx context.Context, id int64, userID int64, req *model.UpdatePostRequest) (*model.Post, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid post id: %d", id)
	}
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user id: %d", userID)
	}

	// 更新するフィールドが空の場合はエラー
	if req.Title == "" && req.Content == "" {
		return nil, fmt.Errorf("no fields to update")
	}

	post, err := s.postRepo.Update(ctx, id, userID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return post, nil
}

// DeletePost は投稿を削除
func (s *PostService) DeletePost(ctx context.Context, id int64, userID int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid post id: %d", id)
	}
	if userID <= 0 {
		return fmt.Errorf("invalid user id: %d", userID)
	}

	if err := s.postRepo.Delete(ctx, id, userID); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}
