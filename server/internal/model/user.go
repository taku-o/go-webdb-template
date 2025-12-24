package model

import "time"

// User はユーザーのデータモデル
type User struct {
	ID        int64     `json:"id,string" db:"id" gorm:"primaryKey"`
	Name      string    `json:"name" db:"name" gorm:"type:varchar(100);not null"`
	Email     string    `json:"email" db:"email" gorm:"type:varchar(255);not null;uniqueIndex:idx_users_email"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}

// TableName はテーブル名を明示的に指定
func (User) TableName() string {
	return "users"
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
