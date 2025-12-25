package admin

import (
	goadminConfig "github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/language"

	"github.com/example/go-webdb-template/internal/config"
)

// Config はGoAdmin管理画面の設定を管理する構造体
type Config struct {
	appConfig *config.Config
}

// NewConfig は新しいConfig構造体を作成する
func NewConfig(appConfig *config.Config) *Config {
	return &Config{
		appConfig: appConfig,
	}
}

// GetGoAdminConfig はGoAdmin用の設定を返す
func (c *Config) GetGoAdminConfig() *goadminConfig.Config {
	return &goadminConfig.Config{
		Env:       c.getEnv(),
		Databases: c.getDatabaseConfig(),
		UrlPrefix: "admin",
		IndexUrl:  "/",
		Store: goadminConfig.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Language:        language.JP,
		Title:           "GoAdmin 管理画面",
		Logo:            "GoAdmin",
		MiniLogo:        "GA",
		Theme:           "adminlte",
		Debug:           c.isDebug(),
		SessionLifeTime: c.getSessionLifetime(),
	}
}

// getSessionLifetime はセッション有効期限を返す
func (c *Config) getSessionLifetime() int {
	if c.appConfig.Admin.Session.Lifetime <= 0 {
		return 7200 // デフォルト2時間
	}
	return c.appConfig.Admin.Session.Lifetime
}

// getDatabaseConfig はGORM Managerを使用したデータベース設定を返す
func (c *Config) getDatabaseConfig() goadminConfig.DatabaseList {
	// 最初のシャードをGoAdmin用データベースとして使用
	dsn := ""
	if len(c.appConfig.Database.Shards) > 0 {
		dsn = c.appConfig.Database.Shards[0].DSN
	}

	return goadminConfig.DatabaseList{
		"default": {
			Driver: "sqlite",
			File:   dsn,
		},
	}
}

// getEnv は環境設定を返す
func (c *Config) getEnv() string {
	// ロギングレベルに基づいて環境を判断
	switch c.appConfig.Logging.Level {
	case "debug":
		return goadminConfig.EnvLocal
	case "info":
		return goadminConfig.EnvTest
	case "warn", "error":
		return goadminConfig.EnvProd
	default:
		return goadminConfig.EnvLocal
	}
}

// isDebug はデバッグモードかどうかを返す
func (c *Config) isDebug() bool {
	return c.appConfig.Logging.Level == "debug"
}

// GetAdminPort は管理画面のポート番号を返す
func (c *Config) GetAdminPort() int {
	return c.appConfig.Admin.Port
}

// GetReadTimeout は読み取りタイムアウトを返す
func (c *Config) GetReadTimeout() int {
	return int(c.appConfig.Admin.ReadTimeout.Seconds())
}

// GetWriteTimeout は書き込みタイムアウトを返す
func (c *Config) GetWriteTimeout() int {
	return int(c.appConfig.Admin.WriteTimeout.Seconds())
}
