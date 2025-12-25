package pages

import (
	"net/http"

	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/menu"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/types"
)

// PageHandler はカスタムページのハンドラー関数型
type PageHandler func(ctx *context.Context, conn db.Connection) (types.Panel, error)

// MakeHandler はPageHandlerをcontext.Handlerに変換する
func MakeHandler(conn db.Connection, handler PageHandler) context.Handler {
	return func(ctx *context.Context) {
		panel, err := handler(ctx, conn)
		if err != nil {
			panel = types.Panel{
				Title:       "エラー",
				Description: "エラーが発生しました",
				Content:     template.HTML("<div class='alert alert-danger'>" + err.Error() + "</div>"),
			}
		}

		user := auth.Auth(ctx)
		tmpl, tmplName := template.Default(ctx).GetTemplate(ctx.IsPjax())

		buf := template.Execute(ctx, &template.ExecuteParam{
			User:     user,
			TmplName: tmplName,
			Tmpl:     tmpl,
			Panel:    panel,
			Config:   config.Get(),
			Menu:     menu.GetGlobalMenu(user, conn, ctx.Lang()).SetActiveClass(config.URLRemovePrefix(ctx.Path())),
			IsPjax:   ctx.IsPjax(),
			Iframe:   ctx.IsIframe(),
		})

		ctx.HTML(http.StatusOK, buf.String())
	}
}

// RegisterCustomPages はカスタムページのハンドラーを返す
func RegisterCustomPages(conn db.Connection) map[string]context.Handler {
	return map[string]context.Handler{
		"/":                  MakeHandler(conn, HomePage),
		"/user/register":     MakeHandler(conn, UserRegisterPage),
		"/user/register/new": MakeHandler(conn, UserRegisterCompletePage),
	}
}
