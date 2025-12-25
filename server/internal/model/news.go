package model

import "time"

// News はニュースのデータモデル
// masterグループに配置されるテーブル（シャーディング不要）
type News struct {
	ID          int64      `json:"id,string" db:"id" gorm:"primaryKey;autoIncrement"`
	Title       string     `json:"title" db:"title" gorm:"type:varchar(255);not null"`
	Content     string     `json:"content" db:"content" gorm:"type:text;not null"`
	AuthorID    *int64     `json:"author_id,omitempty,string" db:"author_id" gorm:"index:idx_news_author_id"`
	PublishedAt *time.Time `json:"published_at,omitempty" db:"published_at" gorm:"index:idx_news_published_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at" gorm:"autoUpdateTime"`
}

// TableName はテーブル名を明示的に指定
func (News) TableName() string {
	return "news"
}

// CreateNewsRequest はニュース作成リクエスト
type CreateNewsRequest struct {
	Title       string     `json:"title" validate:"required,min=1,max=255"`
	Content     string     `json:"content" validate:"required,min=1"`
	AuthorID    *int64     `json:"author_id,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

// UpdateNewsRequest はニュース更新リクエスト
type UpdateNewsRequest struct {
	Title       string     `json:"title" validate:"omitempty,min=1,max=255"`
	Content     string     `json:"content" validate:"omitempty,min=1"`
	AuthorID    *int64     `json:"author_id,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}
