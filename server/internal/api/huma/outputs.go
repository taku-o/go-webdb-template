package humaapi

import "github.com/taku-o/go-webdb-template/internal/model"

// DmUserOutput はユーザー単体のレスポンス構造体
type DmUserOutput struct {
	Body model.DmUser
}

// DmUsersOutput はユーザー一覧のレスポンス構造体
type DmUsersOutput struct {
	Body []*model.DmUser
}

// DeleteDmUserOutput はユーザー削除のレスポンス構造体（204 No Content用）
type DeleteDmUserOutput struct {
}

// DmPostOutput は投稿単体のレスポンス構造体
type DmPostOutput struct {
	Body model.DmPost
}

// DmPostsOutput は投稿一覧のレスポンス構造体
type DmPostsOutput struct {
	Body []*model.DmPost
}

// DmUserPostsOutput はユーザーと投稿のJOIN結果のレスポンス構造体
type DmUserPostsOutput struct {
	Body []*model.DmUserPost
}

// DeleteDmPostOutput は投稿削除のレスポンス構造体（204 No Content用）
type DeleteDmPostOutput struct {
}

// TodayOutput は今日の日付のレスポンス構造体
type TodayOutput struct {
	Body struct {
		Date string `json:"date" doc:"今日の日付（YYYY-MM-DD形式）"`
	}
}
