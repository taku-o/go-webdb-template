package model

import "time"

// Post は投稿のデータモデル
type Post struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreatePostRequest は投稿作成リクエスト
type CreatePostRequest struct {
	UserID  int64  `json:"user_id" validate:"required,gt=0"`
	Title   string `json:"title" validate:"required,min=1,max=200"`
	Content string `json:"content" validate:"required,min=1"`
}

// UpdatePostRequest は投稿更新リクエスト
type UpdatePostRequest struct {
	Title   string `json:"title" validate:"omitempty,min=1,max=200"`
	Content string `json:"content" validate:"omitempty,min=1"`
}

// UserPost はユーザーと投稿を結合したモデル（JOIN結果用）
type UserPost struct {
	PostID      int64     `json:"post_id" db:"post_id"`
	PostTitle   string    `json:"post_title" db:"post_title"`
	PostContent string    `json:"post_content" db:"post_content"`
	UserID      int64     `json:"user_id" db:"user_id"`
	UserName    string    `json:"user_name" db:"user_name"`
	UserEmail   string    `json:"user_email" db:"user_email"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
