package model

import "time"

// DmUser はユーザーのデータモデル（ダミーテーブル）
// IDはUUIDv7形式の32文字の16進数文字列（ハイフンなし小文字）
type DmUser struct {
	ID        string    `json:"id" db:"id" gorm:"primaryKey;type:varchar(32)"`
	Name      string    `json:"name" db:"name" gorm:"type:varchar(100);not null"`
	Email     string    `json:"email" db:"email" gorm:"type:varchar(255);not null;uniqueIndex:idx_dm_users_email"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}

// TableName はテーブル名を明示的に指定
func (DmUser) TableName() string {
	return "dm_users"
}

// CreateDmUserRequest はユーザー作成リクエスト
type CreateDmUserRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=100"`
	Email string `json:"email" validate:"required,email,max=255"`
}

// UpdateDmUserRequest はユーザー更新リクエスト
type UpdateDmUserRequest struct {
	Name  string `json:"name" validate:"omitempty,min=1,max=100"`
	Email string `json:"email" validate:"omitempty,email,max=255"`
}
