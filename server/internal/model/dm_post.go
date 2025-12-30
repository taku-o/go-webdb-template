package model

import "time"

// DmPost は投稿のデータモデル（ダミーテーブル）
// ID, UserIDはUUIDv7形式の32文字の16進数文字列（ハイフンなし小文字）
type DmPost struct {
	ID        string    `json:"id" db:"id" gorm:"primaryKey;type:varchar(32)"`
	UserID    string    `json:"user_id" db:"user_id" gorm:"type:varchar(32);not null;index:idx_dm_posts_user_id"`
	Title     string    `json:"title" db:"title" gorm:"type:varchar(200);not null"`
	Content   string    `json:"content" db:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}

// TableName はテーブル名を明示的に指定
func (DmPost) TableName() string {
	return "dm_posts"
}

// CreateDmPostRequest は投稿作成リクエスト
type CreateDmPostRequest struct {
	UserID  string `json:"user_id" validate:"required,len=32"`
	Title   string `json:"title" validate:"required,min=1,max=200"`
	Content string `json:"content" validate:"required,min=1"`
}

// UpdateDmPostRequest は投稿更新リクエスト
type UpdateDmPostRequest struct {
	Title   string `json:"title" validate:"omitempty,min=1,max=200"`
	Content string `json:"content" validate:"omitempty,min=1"`
}

// DmUserPost はユーザーと投稿を結合したモデル（JOIN結果用）
type DmUserPost struct {
	PostID      string    `json:"post_id" db:"post_id" gorm:"column:post_id"`
	PostTitle   string    `json:"post_title" db:"post_title" gorm:"column:post_title"`
	PostContent string    `json:"post_content" db:"post_content" gorm:"column:post_content"`
	UserID      string    `json:"user_id" db:"user_id" gorm:"column:user_id"`
	UserName    string    `json:"user_name" db:"user_name" gorm:"column:user_name"`
	UserEmail   string    `json:"user_email" db:"user_email" gorm:"column:user_email"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" gorm:"column:created_at"`
}

// TableName はテーブル名なし（JOIN結果用）
func (DmUserPost) TableName() string {
	return ""
}
