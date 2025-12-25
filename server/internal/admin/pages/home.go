package pages

import (
	"fmt"
	"html/template"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// HomePage はダッシュボードページを返す
// 注意: GoAdminはmasterデータベースのみを使用するため、統計情報はmasterグループのテーブル（news）のみ表示
func HomePage(ctx *context.Context, conn db.Connection) (types.Panel, error) {
	// 統計情報を取得（masterグループのテーブルのみ）
	newsCount := getTableCount(conn, "news")

	content := fmt.Sprintf(`
<div class="row">
    <div class="col-lg-3 col-xs-6">
        <div class="small-box bg-yellow">
            <div class="inner">
                <h3>%d</h3>
                <p>ニュース数</p>
            </div>
            <div class="icon">
                <i class="fa fa-newspaper-o"></i>
            </div>
            <a href="/admin/info/news" class="small-box-footer">
                詳細を見る <i class="fa fa-arrow-circle-right"></i>
            </a>
        </div>
    </div>
</div>

<div class="row">
    <div class="col-md-12">
        <div class="box box-primary">
            <div class="box-header with-border">
                <h3 class="box-title">クイックアクション</h3>
            </div>
            <div class="box-body">
                <a href="/admin/info/news/new" class="btn btn-warning">
                    <i class="fa fa-newspaper-o"></i> ニュース作成
                </a>
            </div>
        </div>
    </div>
</div>

<div class="row">
    <div class="col-md-12">
        <div class="box box-info">
            <div class="box-header with-border">
                <h3 class="box-title">システム情報</h3>
            </div>
            <div class="box-body">
                <table class="table table-bordered">
                    <tr>
                        <th style="width: 200px;">プロジェクト名</th>
                        <td>go-webdb-template</td>
                    </tr>
                    <tr>
                        <th>管理画面バージョン</th>
                        <td>GoAdmin v1.2.26</td>
                    </tr>
                </table>
            </div>
        </div>
    </div>
</div>
`, newsCount)

	return types.Panel{
		Title:       "ダッシュボード",
		Description: "管理画面ホーム",
		Content:     template.HTML(content),
	}, nil
}

// getTableCount はテーブルのレコード数を取得する
func getTableCount(conn db.Connection, tableName string) int64 {
	result, err := conn.Query(fmt.Sprintf("SELECT COUNT(*) as count FROM %s", tableName))
	if err != nil || len(result) == 0 {
		return 0
	}

	count, ok := result[0]["count"]
	if !ok {
		return 0
	}

	switch v := count.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	default:
		return 0
	}
}
