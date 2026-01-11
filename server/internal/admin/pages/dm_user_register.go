package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/taku-o/go-webdb-template/internal/usecase/admin"
)

// DmUserRegisterPage はユーザー登録ページを返す
func DmUserRegisterPage(ctx *context.Context, dmUserRegisterUsecase *admin.DmUserRegisterUsecase) (types.Panel, error) {
	if ctx.Method() == http.MethodPost {
		return handleDmUserRegisterPost(ctx, dmUserRegisterUsecase)
	}
	return renderDmUserRegisterForm(ctx, "", "", nil)
}

// handleDmUserRegisterPost はPOSTリクエストを処理する
func handleDmUserRegisterPost(ctx *context.Context, dmUserRegisterUsecase *admin.DmUserRegisterUsecase) (types.Panel, error) {
	name := strings.TrimSpace(ctx.FormValue("name"))
	email := strings.TrimSpace(ctx.FormValue("email"))

	// バリデーション
	errors := validateDmUserInput(name, email)
	if len(errors) > 0 {
		return renderDmUserRegisterForm(ctx, name, email, errors)
	}

	// usecase層を呼び出し
	dmUserID, err := dmUserRegisterUsecase.RegisterDmUser(ctx.Request.Context(), name, email)
	if err != nil {
		return renderDmUserRegisterForm(ctx, name, email, []string{err.Error()})
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
