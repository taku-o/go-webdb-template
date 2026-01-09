package admin

import (
	"fmt"

	goadminConfig "github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/language"

	"github.com/taku-o/go-webdb-template/internal/config"
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

// getDatabaseConfig はGoAdmin用のデータベース設定を返す
func (c *Config) getDatabaseConfig() goadminConfig.DatabaseList {
	// masterグループのデータベースをGoAdmin用データベースとして使用
	if len(c.appConfig.Database.Groups.Master) == 0 {
		panic("no database configuration found: master group is required")
	}

	masterDB := c.appConfig.Database.Groups.Master[0]

	// 接続情報の検証
	if masterDB.Host == "" || masterDB.Port == 0 || masterDB.User == "" || masterDB.Name == "" {
		panic("incomplete database configuration: host, port, user, and name are required")
	}

	// PostgreSQL用のDSN形式を構築
	// host=... port=... user=... password=... dbname=... sslmode=disable
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		masterDB.Host,
		masterDB.Port,
		masterDB.User,
		masterDB.Password,
		masterDB.Name,
	)

	return goadminConfig.DatabaseList{
		"default": {
			Driver: "postgresql",
			Dsn:    dsn,
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
