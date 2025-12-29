package humaapi

import "github.com/taku-o/go-webdb-template/internal/model"

// UserOutput はユーザー単体のレスポンス構造体
type UserOutput struct {
	Body model.DmUser
}

// UsersOutput はユーザー一覧のレスポンス構造体
type UsersOutput struct {
	Body []*model.DmUser
}

// DeleteUserOutput はユーザー削除のレスポンス構造体（204 No Content用）
type DeleteUserOutput struct {
}

// PostOutput は投稿単体のレスポンス構造体
type PostOutput struct {
	Body model.DmPost
}

// PostsOutput は投稿一覧のレスポンス構造体
type PostsOutput struct {
	Body []*model.DmPost
}

// UserPostsOutput はユーザーと投稿のJOIN結果のレスポンス構造体
type UserPostsOutput struct {
	Body []*model.DmUserPost
}

// DeletePostOutput は投稿削除のレスポンス構造体（204 No Content用）
type DeletePostOutput struct {
}

// TodayOutput は今日の日付のレスポンス構造体
type TodayOutput struct {
	Body struct {
		Date string `json:"date" doc:"今日の日付（YYYY-MM-DD形式）"`
	}
}
