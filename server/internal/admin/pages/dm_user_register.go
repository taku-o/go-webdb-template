package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/template/types"
	appdb "github.com/taku-o/go-webdb-template/internal/db"
)

// UserRegisterPage はユーザー登録ページを返す
func UserRegisterPage(ctx *context.Context, groupManager *appdb.GroupManager) (types.Panel, error) {
	if ctx.Method() == http.MethodPost {
		return handleUserRegisterPost(ctx, groupManager)
	}
	return renderUserRegisterForm(ctx, "", "", nil)
}

// handleUserRegisterPost はPOSTリクエストを処理する
func handleUserRegisterPost(ctx *context.Context, groupManager *appdb.GroupManager) (types.Panel, error) {
	name := strings.TrimSpace(ctx.FormValue("name"))
	email := strings.TrimSpace(ctx.FormValue("email"))

	// バリデーション
	errors := validateUserInput(name, email)
	if len(errors) > 0 {
		return renderUserRegisterForm(ctx, name, email, errors)
	}

	// メールアドレスの重複チェック（全シャードを検索）
	exists, err := checkEmailExistsSharded(groupManager, email)
	if err != nil {
		return renderUserRegisterForm(ctx, name, email, []string{"データベースエラーが発生しました"})
	}
	if exists {
		return renderUserRegisterForm(ctx, name, email, []string{"このメールアドレスは既に登録されています"})
	}

	// ユーザー登録（シャーディング対応）
	userID, err := insertUserSharded(groupManager, name, email)
	if err != nil {
		return renderUserRegisterForm(ctx, name, email, []string{"ユーザー登録に失敗しました: " + err.Error()})
	}

	// 登録完了ページへリダイレクト（クエリパラメータで情報を渡す）
	redirectURL := fmt.Sprintf("/admin/dm-user/register/new?id=%d&name=%s&email=%s",
		userID,
		url.QueryEscape(name),
		url.QueryEscape(email),
	)

	// GoAdminのContent wrapperはctx.Redirectを上書きするため、
	// JavaScriptリダイレクトを使用
	return types.Panel{
		Title:       "リダイレクト中",
		Description: "",
		Content:     template.HTML(fmt.Sprintf(`<script>window.location.href='%s';</script>`, redirectURL)),
	}, nil
}

// validateUserInput は入力値をバリデーションする
func validateUserInput(name, email string) []string {
	var errors []string

	if name == "" {
		errors = append(errors, "名前は必須です")
	} else if len(name) > 100 {
		errors = append(errors, "名前は100文字以内で入力してください")
	}

	if email == "" {
		errors = append(errors, "メールアドレスは必須です")
	} else if !strings.Contains(email, "@") {
		errors = append(errors, "有効なメールアドレスを入力してください")
	} else if len(email) > 255 {
		errors = append(errors, "メールアドレスは255文字以内で入力してください")
	}

	return errors
}

// checkEmailExistsSharded はメールアドレスが既に存在するかチェックする（全シャード検索）
func checkEmailExistsSharded(groupManager *appdb.GroupManager, email string) (bool, error) {
	// 全テーブルを検索
	for tableNum := 0; tableNum < appdb.DBShardingTableCount; tableNum++ {
		conn, err := groupManager.GetShardingConnection(tableNum)
		if err != nil {
			return false, fmt.Errorf("failed to get connection for table %d: %w", tableNum, err)
		}

		sqlDB, err := conn.DB.DB()
		if err != nil {
			return false, fmt.Errorf("failed to get sql.DB for table %d: %w", tableNum, err)
		}

		tableName := fmt.Sprintf("dm_users_%03d", tableNum)
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE email = ?", tableName)

		var count int
		err = sqlDB.QueryRow(query, email).Scan(&count)
		if err != nil {
			return false, fmt.Errorf("failed to check email in %s: %w", tableName, err)
		}

		if count > 0 {
			return true, nil
		}
	}

	return false, nil
}

// insertUserSharded はユーザーを登録する（シャーディング対応）
func insertUserSharded(groupManager *appdb.GroupManager, name, email string) (int64, error) {
	now := time.Now()

	// IDを生成（タイムスタンプベース）
	userID := now.UnixNano()

	// テーブル番号を計算
	tableNumber := int(userID % appdb.DBShardingTableCount)
	tableName := fmt.Sprintf("dm_users_%03d", tableNumber)

	// 接続の取得
	conn, err := groupManager.GetShardingConnection(tableNumber)
	if err != nil {
		return 0, fmt.Errorf("failed to get sharding connection: %w", err)
	}

	// sql.DBを取得
	sqlDB, err := conn.DB.DB()
	if err != nil {
		return 0, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// ユーザーを挿入
	query := fmt.Sprintf(`
		INSERT INTO %s (id, name, email, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, tableName)

	_, err = sqlDB.Exec(query, userID, name, email, now, now)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	return userID, nil
}

// renderUserRegisterForm はユーザー登録フォームをレンダリングする
func renderUserRegisterForm(ctx *context.Context, name, email string, errors []string) (types.Panel, error) {
	errorHTML := ""
	if len(errors) > 0 {
		errorHTML = `<div class="alert alert-danger"><ul>`
		for _, e := range errors {
			errorHTML += fmt.Sprintf("<li>%s</li>", e)
		}
		errorHTML += `</ul></div>`
	}

	content := fmt.Sprintf(`
%s
<div class="box box-primary">
    <div class="box-header with-border">
        <h3 class="box-title">ユーザー情報入力</h3>
    </div>
    <form action="/admin/dm-user/register" method="POST">
        <div class="box-body">
            <div class="form-group">
                <label for="name">名前 <span class="text-red">*</span></label>
                <input type="text" class="form-control" id="name" name="name" value="%s" placeholder="名前を入力" required maxlength="100">
            </div>
            <div class="form-group">
                <label for="email">メールアドレス <span class="text-red">*</span></label>
                <input type="email" class="form-control" id="email" name="email" value="%s" placeholder="メールアドレスを入力" required maxlength="255">
            </div>
        </div>
        <div class="box-footer">
            <button type="submit" class="btn btn-primary">登録</button>
            <a href="/admin" class="btn btn-default">キャンセル</a>
        </div>
    </form>
</div>
`, errorHTML, template.HTMLEscapeString(name), template.HTMLEscapeString(email))

	return types.Panel{
		Title:       "ユーザー登録",
		Description: "新規ユーザー情報を入力してください",
		Content:     template.HTML(content),
	}, nil
}
