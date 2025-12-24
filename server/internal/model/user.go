package model

import "time"

// User はユーザーのデータモデル
type User struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateUserRequest はユーザー作成リクエスト
type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=100"`
	Email string `json:"email" validate:"required,email,max=255"`
}

// UpdateUserRequest はユーザー更新リクエスト
type UpdateUserRequest struct {
	Name  string `json:"name" validate:"omitempty,min=1,max=100"`
	Email string `json:"email" validate:"omitempty,email,max=255"`
}
