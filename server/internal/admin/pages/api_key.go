package pages

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/template/types"

	"github.com/example/go-webdb-template/internal/auth"
	"github.com/example/go-webdb-template/internal/config"
)

// APIKeyPage はAPIキー発行ページを返す
// 注意: RegisterCustomPagesで"/api-key"と登録すると、実際のURLは"/admin/api-key"になる
// HTML内のリンクも"/admin/api-key"とする必要がある
func APIKeyPage(ctx *context.Context, conn db.Connection) (types.Panel, error) {
	// 設定を取得
	cfg, err := config.Load()
	if err != nil {
		return types.Panel{}, err
	}

	// POSTリクエスト: キー生成
	if ctx.Method() == http.MethodPost {
		return handleGenerateKey(ctx, cfg)
	}

	// GETリクエスト: フォーム表示
	return renderAPIKeyPage(ctx, cfg)
}

// handleGenerateKey はAPIキーを生成
func handleGenerateKey(ctx *context.Context, cfg *config.Config) (types.Panel, error) {
	// 現在の環境を取得
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "develop"
	}

	// JWTトークンを生成
	token, err := generatePublicAPIKey(cfg, env)
	if err != nil {
		return types.Panel{}, err
	}

	// ペイロードをデコード
	claims, err := auth.ParseJWTClaims(token)
	if err != nil {
		return types.Panel{}, err
	}

	// 生成結果を表示
	return renderAPIKeyResult(ctx, token, claims)
}

// generatePublicAPIKey はPublic JWTキーを生成
func generatePublicAPIKey(cfg *config.Config, env string) (string, error) {
	now := time.Now()
	return auth.GeneratePublicAPIKey(cfg.API.SecretKey, cfg.API.CurrentVersion, env, now.Unix())
}

// renderAPIKeyPage はAPIキー発行ページをレンダリング
func renderAPIKeyPage(ctx *context.Context, cfg *config.Config) (types.Panel, error) {
	content := `
<div class="box box-primary">
    <div class="box-header with-border">
        <h3 class="box-title">Public APIキー発行</h3>
    </div>
    <div class="box-body">
        <p>新しいPublic APIキーを発行します。</p>
        <form action="/admin/api-key" method="POST">
            <button type="submit" class="btn btn-primary">
                <i class="fa fa-key"></i> APIキーを発行
            </button>
        </form>
    </div>
</div>
`

	return types.Panel{
		Title:       "APIキー発行",
		Description: "Public APIキーを発行します",
		Content:     template.HTML(content),
	}, nil
}

// renderAPIKeyResult は生成結果をレンダリング
func renderAPIKeyResult(ctx *context.Context, token string, claims *auth.JWTClaims) (types.Panel, error) {
	// ペイロードをJSON形式で整形
	payloadJSON, _ := json.MarshalIndent(claims, "", "  ")

	// iatを人間が読める形式に変換
	issuedAt := time.Unix(claims.IssuedAt, 0).Format("2006-01-02 15:04:05")

	content := fmt.Sprintf(`
<div class="box box-success">
    <div class="box-header with-border">
        <h3 class="box-title">APIキー発行結果</h3>
    </div>
    <div class="box-body">
        <div class="form-group">
            <label>JWTトークン</label>
            <textarea class="form-control" rows="3" readonly>%s</textarea>
        </div>
        <div class="form-group">
            <label>JWTペイロード</label>
            <pre class="form-control" style="height: 300px; overflow-y: auto;">%s</pre>
        </div>
        <div class="form-group">
            <label>発行日時</label>
            <p>%s</p>
        </div>
        <div class="form-group">
            <label>バージョン</label>
            <p>%s</p>
        </div>
        <div class="form-group">
            <label>環境</label>
            <p>%s</p>
        </div>
        <div class="form-group">
            <button type="button" class="btn btn-success" onclick="downloadAPIKey()">
                <i class="fa fa-download"></i> ダウンロード
            </button>
        </div>
    </div>
</div>
<script>
function downloadAPIKey() {
    var token = %q;
    var timestamp = new Date().toISOString().replace(/[-:T]/g, '').slice(0, 15);
    var filename = 'api-key-' + timestamp + '.txt';
    var blob = new Blob([token], {type: 'text/plain'});
    var url = window.URL.createObjectURL(blob);
    var a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
}
</script>
`, template.HTMLEscapeString(token), template.HTMLEscapeString(string(payloadJSON)), issuedAt, claims.Version, claims.Env, token)

	return types.Panel{
		Title:       "APIキー発行結果",
		Description: "Public APIキーが発行されました",
		Content:     template.HTML(content),
	}, nil
}
