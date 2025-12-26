package humaapi

// CreateUserInput はユーザー作成リクエストの入力構造体
type CreateUserInput struct {
	Body struct {
		Name  string `json:"name" maxLength:"100" doc:"ユーザー名"`
		Email string `json:"email" format:"email" maxLength:"255" doc:"メールアドレス"`
	}
}

// GetUserInput はユーザー取得リクエストの入力構造体
type GetUserInput struct {
	ID int64 `path:"id" doc:"ユーザーID"`
}

// ListUsersInput はユーザー一覧取得リクエストの入力構造体
type ListUsersInput struct {
	Limit  int `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"取得件数"`
	Offset int `query:"offset" default:"0" minimum:"0" doc:"オフセット"`
}

// UpdateUserInput はユーザー更新リクエストの入力構造体
type UpdateUserInput struct {
	ID   int64 `path:"id" doc:"ユーザーID"`
	Body struct {
		Name  string `json:"name,omitempty" maxLength:"100" doc:"ユーザー名"`
		Email string `json:"email,omitempty" format:"email" maxLength:"255" doc:"メールアドレス"`
	}
}

// DeleteUserInput はユーザー削除リクエストの入力構造体
type DeleteUserInput struct {
	ID int64 `path:"id" doc:"ユーザーID"`
}

// CreatePostInput は投稿作成リクエストの入力構造体
type CreatePostInput struct {
	Body struct {
		UserID  string `json:"user_id" doc:"ユーザーID"`
		Title   string `json:"title" maxLength:"200" doc:"タイトル"`
		Content string `json:"content" doc:"内容"`
	}
}

// GetPostInput は投稿取得リクエストの入力構造体
type GetPostInput struct {
	ID     int64 `path:"id" doc:"投稿ID"`
	UserID int64 `query:"user_id" required:"true" minimum:"1" doc:"ユーザーID"`
}

// ListPostsInput は投稿一覧取得リクエストの入力構造体
type ListPostsInput struct {
	Limit  int   `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"取得件数"`
	Offset int   `query:"offset" default:"0" minimum:"0" doc:"オフセット"`
	UserID int64 `query:"user_id" default:"0" doc:"ユーザーID（0の場合は全件取得）"`
}

// UpdatePostInput は投稿更新リクエストの入力構造体
type UpdatePostInput struct {
	ID     int64 `path:"id" doc:"投稿ID"`
	UserID int64 `query:"user_id" required:"true" minimum:"1" doc:"ユーザーID"`
	Body   struct {
		Title   string `json:"title,omitempty" maxLength:"200" doc:"タイトル"`
		Content string `json:"content,omitempty" doc:"内容"`
	}
}

// DeletePostInput は投稿削除リクエストの入力構造体
type DeletePostInput struct {
	ID     int64 `path:"id" doc:"投稿ID"`
	UserID int64 `query:"user_id" required:"true" minimum:"1" doc:"ユーザーID"`
}

// GetUserPostsInput はユーザー投稿一覧取得リクエストの入力構造体
type GetUserPostsInput struct {
	Limit  int `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"取得件数"`
	Offset int `query:"offset" default:"0" minimum:"0" doc:"オフセット"`
}
