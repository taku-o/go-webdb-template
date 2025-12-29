package admin

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

// GetNewsTable はdm_newsテーブルのGoAdmin設定を返す
// 注意: GoAdminはmasterグループのデータベースのみを使用するため、
// dm_users/dm_postsテーブル（shardingグループ）はGoAdminで管理できません
func GetNewsTable(ctx *context.Context) table.Table {
	newsTable := table.NewDefaultTable(ctx, table.Config{
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
	info := newsTable.GetInfo()
	info.AddField("ID", "id", db.Int).FieldSortable()
	info.AddField("タイトル", "title", db.Varchar).FieldSortable().FieldFilterable()
	info.AddField("内容", "content", db.Text)
	info.AddField("作成者ID", "author_id", db.Int).FieldSortable().FieldFilterable()
	info.AddField("公開日時", "published_at", db.Datetime).FieldSortable().FieldFilterable()
	info.AddField("作成日時", "created_at", db.Datetime).FieldSortable()
	info.AddField("更新日時", "updated_at", db.Datetime).FieldSortable()

	info.SetTable("dm_news").SetTitle("ニュース").SetDescription("ニュース一覧")

	// フォーム設定（新規作成・編集）
	formList := newsTable.GetForm()
	formList.AddField("ID", "id", db.Int, form.Default).FieldNotAllowEdit().FieldNotAllowAdd()
	formList.AddField("タイトル", "title", db.Varchar, form.Text).FieldMust()
	formList.AddField("内容", "content", db.Text, form.TextArea).FieldMust()
	formList.AddField("作成者ID", "author_id", db.Int, form.Number).
		FieldPostFilterFn(func(value types.PostFieldModel) interface{} {
			if value.Value.Value() == "" {
				return nil
			}
			return value.Value.Value()
		})
	formList.AddField("公開日時", "published_at", db.Datetime, form.Datetime)
	formList.AddField("作成日時", "created_at", db.Datetime, form.Datetime).
		FieldHide().
		FieldNowWhenInsert().
		FieldDisableWhenUpdate()
	formList.AddField("更新日時", "updated_at", db.Datetime, form.Datetime).
		FieldHide().
		FieldNowWhenInsert().
		FieldNowWhenUpdate()

	formList.SetTable("dm_news").SetTitle("ニュース").SetDescription("ニュース情報")

	return newsTable
}

// Generators はGoAdminに登録するテーブルジェネレータのマップ
// 注意: dm_usersとdm_postsはシャーディンググループにあるため、
// GoAdmin（masterグループのみ使用）では管理できません
var Generators = map[string]table.Generator{
	"dm_news": GetNewsTable,
}
