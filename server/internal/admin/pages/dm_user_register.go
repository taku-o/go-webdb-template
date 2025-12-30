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
	"github.com/taku-o/go-webdb-template/internal/util/idgen"
)

// DmUserRegisterPage はユーザー登録ページを返す
func DmUserRegisterPage(ctx *context.Context, groupManager *appdb.GroupManager) (types.Panel, error) {
	if ctx.Method() == http.MethodPost {
		return handleDmUserRegisterPost(ctx, groupManager)
	}
	return renderDmUserRegisterForm(ctx, "", "", nil)
}

// handleDmUserRegisterPost はPOSTリクエストを処理する
func handleDmUserRegisterPost(ctx *context.Context, groupManager *appdb.GroupManager) (types.Panel, error) {
	name := strings.TrimSpace(ctx.FormValue("name"))
	email := strings.TrimSpace(ctx.FormValue("email"))

	// バリデーション
	errors := validateDmUserInput(name, email)
	if len(errors) > 0 {
		return renderDmUserRegisterForm(ctx, name, email, errors)
	}

	// メールアドレスの重複チェック（全シャードを検索）
	exists, err := checkEmailExistsSharded(groupManager, email)
	if err != nil {
		return renderDmUserRegisterForm(ctx, name, email, []string{"データベースエラーが発生しました"})
	}
	if exists {
		return renderDmUserRegisterForm(ctx, name, email, []string{"このメールアドレスは既に登録されています"})
	}

	// dm_user登録（シャーディング対応）
	dmUserID, err := insertDmUserSharded(groupManager, name, email)
	if err != nil {
		return renderDmUserRegisterForm(ctx, name, email, []string{"ユーザー登録に失敗しました: " + err.Error()})
	}

	// 登録完了ページへリダイレクト（クエリパラメータで情報を渡す）
	redirectURL := fmt.Sprintf("/admin/dm-user/register/new?id=%s&name=%s&email=%s",
		url.QueryEscape(dmUserID),
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

// validateDmUserInput は入力値をバリデーションする
func validateDmUserInput(name, email string) []string {
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

// insertDmUserSharded はdm_userを登録する（シャーディング対応）
func insertDmUserSharded(groupManager *appdb.GroupManager, name, email string) (string, error) {
	now := time.Now()

	// UUIDv7でIDを生成
	dmUserID, err := idgen.GenerateUUIDv7()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUIDv7: %w", err)
	}

	// UUIDからテーブル番号を計算
	selector := appdb.NewTableSelector(appdb.DBShardingTableCount, appdb.DBShardingTablesPerDB)
	tableNumber, err := selector.GetTableNumberFromUUID(dmUserID)
	if err != nil {
		return "", fmt.Errorf("failed to get table number: %w", err)
	}
	tableName := fmt.Sprintf("dm_users_%03d", tableNumber)

	// 接続の取得
	conn, err := groupManager.GetShardingConnection(tableNumber)
	if err != nil {
		return "", fmt.Errorf("failed to get sharding connection: %w", err)
	}

	// sql.DBを取得
	sqlDB, err := conn.DB.DB()
	if err != nil {
		return "", fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// dm_userを挿入
	query := fmt.Sprintf(`
		INSERT INTO %s (id, name, email, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, tableName)

	_, err = sqlDB.Exec(query, dmUserID, name, email, now, now)
	if err != nil {
		return "", fmt.Errorf("failed to insert dm_user: %w", err)
	}

	return dmUserID, nil
}

// renderDmUserRegisterForm はユーザー登録フォームをレンダリングする
func renderDmUserRegisterForm(ctx *context.Context, name, email string, errors []string) (types.Panel, error) {
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
