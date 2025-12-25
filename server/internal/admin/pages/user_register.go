package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// UserRegisterPage はユーザー登録ページを返す
func UserRegisterPage(ctx *context.Context, conn db.Connection) (types.Panel, error) {
	if ctx.Method() == http.MethodPost {
		return handleUserRegisterPost(ctx, conn)
	}
	return renderUserRegisterForm(ctx, "", "", nil)
}

// handleUserRegisterPost はPOSTリクエストを処理する
func handleUserRegisterPost(ctx *context.Context, conn db.Connection) (types.Panel, error) {
	name := strings.TrimSpace(ctx.FormValue("name"))
	email := strings.TrimSpace(ctx.FormValue("email"))

	// バリデーション
	errors := validateUserInput(name, email)
	if len(errors) > 0 {
		return renderUserRegisterForm(ctx, name, email, errors)
	}

	// メールアドレスの重複チェック
	exists, err := checkEmailExists(conn, email)
	if err != nil {
		return renderUserRegisterForm(ctx, name, email, []string{"データベースエラーが発生しました"})
	}
	if exists {
		return renderUserRegisterForm(ctx, name, email, []string{"このメールアドレスは既に登録されています"})
	}

	// ユーザー登録
	userID, err := insertUser(conn, name, email)
	if err != nil {
		return renderUserRegisterForm(ctx, name, email, []string{"ユーザー登録に失敗しました: " + err.Error()})
	}

	// 登録完了ページへリダイレクト
	ctx.SetCookie(&http.Cookie{
		Name:  "registered_user_id",
		Value: fmt.Sprintf("%d", userID),
		Path:  "/",
	})
	ctx.SetCookie(&http.Cookie{
		Name:  "registered_user_name",
		Value: name,
		Path:  "/",
	})
	ctx.SetCookie(&http.Cookie{
		Name:  "registered_user_email",
		Value: email,
		Path:  "/",
	})

	ctx.Redirect("/admin/user/register/new")
	return types.Panel{}, nil
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

// checkEmailExists はメールアドレスが既に存在するかチェックする
func checkEmailExists(conn db.Connection, email string) (bool, error) {
	result, err := conn.Query("SELECT COUNT(*) as count FROM users WHERE email = ?", email)
	if err != nil {
		return false, err
	}
	if len(result) == 0 {
		return false, nil
	}

	count, ok := result[0]["count"]
	if !ok {
		return false, nil
	}

	switch v := count.(type) {
	case int64:
		return v > 0, nil
	case int:
		return v > 0, nil
	case float64:
		return v > 0, nil
	default:
		return false, nil
	}
}

// insertUser はユーザーを登録する
func insertUser(conn db.Connection, name, email string) (int64, error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	result, err := conn.Exec(
		"INSERT INTO users (name, email, created_at, updated_at) VALUES (?, ?, ?, ?)",
		name, email, now, now,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
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
    <form action="/admin/user/register" method="POST">
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
            <a href="/admin/info/users" class="btn btn-default">キャンセル</a>
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
