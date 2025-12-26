package humaapi

import "github.com/example/go-webdb-template/internal/model"

// UserOutput はユーザー単体のレスポンス構造体
type UserOutput struct {
	Body model.User
}

// UsersOutput はユーザー一覧のレスポンス構造体
type UsersOutput struct {
	Body []*model.User
}

// DeleteUserOutput はユーザー削除のレスポンス構造体（204 No Content用）
type DeleteUserOutput struct {
}

// PostOutput は投稿単体のレスポンス構造体
type PostOutput struct {
	Body model.Post
}

// PostsOutput は投稿一覧のレスポンス構造体
type PostsOutput struct {
	Body []*model.Post
}

// UserPostsOutput はユーザーと投稿のJOIN結果のレスポンス構造体
type UserPostsOutput struct {
	Body []*model.UserPost
}

// DeletePostOutput は投稿削除のレスポンス構造体（204 No Content用）
type DeletePostOutput struct {
}
