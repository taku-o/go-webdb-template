package humaapi

// CreateDmUserInput はユーザー作成リクエストの入力構造体
type CreateDmUserInput struct {
	Body struct {
		Name  string `json:"name" required:"true" maxLength:"100" doc:"ユーザー名"`
		Email string `json:"email" required:"true" format:"email" maxLength:"255" doc:"メールアドレス"`
	}
}

// GetDmUserInput はユーザー取得リクエストの入力構造体
type GetDmUserInput struct {
	ID int64 `path:"id" doc:"ユーザーID"`
}

// ListDmUsersInput はユーザー一覧取得リクエストの入力構造体
type ListDmUsersInput struct {
	Limit  int `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"取得件数"`
	Offset int `query:"offset" default:"0" minimum:"0" doc:"オフセット"`
}

// UpdateDmUserInput はユーザー更新リクエストの入力構造体
type UpdateDmUserInput struct {
	ID   int64 `path:"id" doc:"ユーザーID"`
	Body struct {
		Name  string `json:"name,omitempty" maxLength:"100" doc:"ユーザー名"`
		Email string `json:"email,omitempty" format:"email" maxLength:"255" doc:"メールアドレス"`
	}
}

// DeleteDmUserInput はユーザー削除リクエストの入力構造体
type DeleteDmUserInput struct {
	ID int64 `path:"id" doc:"ユーザーID"`
}

// CreateDmPostInput は投稿作成リクエストの入力構造体
type CreateDmPostInput struct {
	Body struct {
		UserID  int64  `json:"user_id" required:"true" minimum:"1" doc:"ユーザーID"`
		Title   string `json:"title" required:"true" maxLength:"200" doc:"タイトル"`
		Content string `json:"content" required:"true" doc:"内容"`
	}
}

// GetDmPostInput は投稿取得リクエストの入力構造体
type GetDmPostInput struct {
	ID     int64 `path:"id" doc:"投稿ID"`
	UserID int64 `query:"user_id" required:"true" minimum:"1" doc:"ユーザーID"`
}

// ListDmPostsInput は投稿一覧取得リクエストの入力構造体
type ListDmPostsInput struct {
	Limit  int   `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"取得件数"`
	Offset int   `query:"offset" default:"0" minimum:"0" doc:"オフセット"`
	UserID int64 `query:"user_id" default:"0" doc:"ユーザーID（0の場合は全件取得）"`
}

// UpdateDmPostInput は投稿更新リクエストの入力構造体
type UpdateDmPostInput struct {
	ID     int64 `path:"id" doc:"投稿ID"`
	UserID int64 `query:"user_id" required:"true" minimum:"1" doc:"ユーザーID"`
	Body   struct {
		Title   string `json:"title,omitempty" maxLength:"200" doc:"タイトル"`
		Content string `json:"content,omitempty" doc:"内容"`
	}
}

// DeleteDmPostInput は投稿削除リクエストの入力構造体
type DeleteDmPostInput struct {
	ID     int64 `path:"id" doc:"投稿ID"`
	UserID int64 `query:"user_id" required:"true" minimum:"1" doc:"ユーザーID"`
}

// GetDmUserPostsInput はユーザー投稿一覧取得リクエストの入力構造体
type GetDmUserPostsInput struct {
	Limit  int `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"取得件数"`
	Offset int `query:"offset" default:"0" minimum:"0" doc:"オフセット"`
}

// GetTodayInput は今日の日付取得リクエストの入力構造体
type GetTodayInput struct {
}
