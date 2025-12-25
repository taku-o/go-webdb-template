package admin

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

// GetUserTable はUsersテーブルのGoAdmin設定を返す
func GetUserTable(ctx *context.Context) table.Table {
	userTable := table.NewDefaultTable(ctx, table.Config{
		Driver:     db.DriverSqlite,
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		Connection: table.DefaultConnectionName,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "id",
		},
	})

	// 一覧表示設定
	info := userTable.GetInfo()
	info.AddField("ID", "id", db.Int).FieldSortable()
	info.AddField("名前", "name", db.Varchar).FieldSortable().FieldFilterable()
	info.AddField("メールアドレス", "email", db.Varchar).FieldSortable().FieldFilterable()
	info.AddField("作成日時", "created_at", db.Datetime).FieldSortable()
	info.AddField("更新日時", "updated_at", db.Datetime).FieldSortable()

	info.SetTable("users").SetTitle("ユーザー").SetDescription("ユーザー一覧")

	// フォーム設定（新規作成・編集）
	formList := userTable.GetForm()
	formList.AddField("ID", "id", db.Int, form.Default).FieldNotAllowEdit().FieldNotAllowAdd()
	formList.AddField("名前", "name", db.Varchar, form.Text).FieldMust()
	formList.AddField("メールアドレス", "email", db.Varchar, form.Email).FieldMust()
	formList.AddField("作成日時", "created_at", db.Datetime, form.Datetime).
		FieldHide().
		FieldNowWhenInsert().
		FieldDisableWhenUpdate()
	formList.AddField("更新日時", "updated_at", db.Datetime, form.Datetime).
		FieldHide().
		FieldNowWhenInsert().
		FieldNowWhenUpdate()

	formList.SetTable("users").SetTitle("ユーザー").SetDescription("ユーザー情報")

	return userTable
}

// GetPostTable はPostsテーブルのGoAdmin設定を返す
func GetPostTable(ctx *context.Context) table.Table {
	postTable := table.NewDefaultTable(ctx, table.Config{
		Driver:     db.DriverSqlite,
		CanAdd:     true,
		Editable:   true,
		Deletable:  true,
		Exportable: true,
		Connection: table.DefaultConnectionName,
		PrimaryKey: table.PrimaryKey{
			Type: db.Int,
			Name: "id",
		},
	})

	// 一覧表示設定
	info := postTable.GetInfo()
	info.AddField("ID", "id", db.Int).FieldSortable()
	info.AddField("ユーザーID", "user_id", db.Int).FieldSortable().FieldFilterable()
	info.AddField("タイトル", "title", db.Varchar).FieldSortable().FieldFilterable()
	info.AddField("内容", "content", db.Text)
	info.AddField("作成日時", "created_at", db.Datetime).FieldSortable()
	info.AddField("更新日時", "updated_at", db.Datetime).FieldSortable()

	info.SetTable("posts").SetTitle("投稿").SetDescription("投稿一覧")

	// フォーム設定（新規作成・編集）
	formList := postTable.GetForm()
	formList.AddField("ID", "id", db.Int, form.Default).FieldNotAllowEdit().FieldNotAllowAdd()
	formList.AddField("ユーザーID", "user_id", db.Int, form.Number).FieldMust()
	formList.AddField("タイトル", "title", db.Varchar, form.Text).FieldMust()
	formList.AddField("内容", "content", db.Text, form.TextArea).FieldMust()
	formList.AddField("作成日時", "created_at", db.Datetime, form.Datetime).
		FieldHide().
		FieldNowWhenInsert().
		FieldDisableWhenUpdate()
	formList.AddField("更新日時", "updated_at", db.Datetime, form.Datetime).
		FieldHide().
		FieldNowWhenInsert().
		FieldNowWhenUpdate()

	formList.SetTable("posts").SetTitle("投稿").SetDescription("投稿情報")

	return postTable
}

// Generators はGoAdminに登録するテーブルジェネレータのマップ
var Generators = map[string]table.Generator{
	"users": GetUserTable,
	"posts": GetPostTable,
}
