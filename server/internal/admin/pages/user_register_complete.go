package pages

import (
	"fmt"
	"html/template"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// UserRegisterCompletePage はユーザー登録完了ページを返す
func UserRegisterCompletePage(ctx *context.Context, conn db.Connection) (types.Panel, error) {
	_ = conn // connは将来の拡張用
	// Cookieから登録情報を取得
	userID := ctx.Cookie("registered_user_id")
	userName := ctx.Cookie("registered_user_name")
	userEmail := ctx.Cookie("registered_user_email")

	if userID == "" || userName == "" || userEmail == "" {
		// 情報がない場合は一覧へリダイレクト
		ctx.Redirect("/admin/info/users")
		return types.Panel{}, nil
	}

	registeredAt := time.Now().Format("2006-01-02 15:04:05")

	content := fmt.Sprintf(`
<div class="callout callout-success">
    <h4><i class="fa fa-check"></i> 登録が正常に処理されました</h4>
    <p>ユーザー情報が正常に登録されました。</p>
</div>

<div class="box box-success">
    <div class="box-header with-border">
        <h3 class="box-title">登録されたユーザー情報</h3>
    </div>
    <div class="box-body">
        <table class="table table-bordered">
            <tr>
                <th style="width: 150px;">ID</th>
                <td>%s</td>
            </tr>
            <tr>
                <th>名前</th>
                <td>%s</td>
            </tr>
            <tr>
                <th>メールアドレス</th>
                <td>%s</td>
            </tr>
            <tr>
                <th>登録日時</th>
                <td>%s</td>
            </tr>
        </table>
    </div>
    <div class="box-footer">
        <a href="/admin/info/users" class="btn btn-primary">
            <i class="fa fa-list"></i> ユーザー一覧に戻る
        </a>
        <a href="/admin/user/register" class="btn btn-success">
            <i class="fa fa-plus"></i> 新規登録を続ける
        </a>
    </div>
</div>
`, template.HTMLEscapeString(userID),
		template.HTMLEscapeString(userName),
		template.HTMLEscapeString(userEmail),
		registeredAt)

	return types.Panel{
		Title:       "登録処理中",
		Description: "ユーザー登録結果",
		Content:     template.HTML(content),
	}, nil
}
