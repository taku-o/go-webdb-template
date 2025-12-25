package auth

import (
	"github.com/GoAdminGroup/go-admin/modules/config"
)

// ConfigureSession はセッション設定を構成する
func ConfigureSession(cfg *config.Config, lifetime int) {
	// GoAdminのセッション設定
	// セッションの有効期限はGoAdminの設定で管理される
	cfg.SessionLifeTime = lifetime
}

// GetSessionLifetime は設定からセッション有効期限を取得する
func GetSessionLifetime(lifetime int) int {
	if lifetime <= 0 {
		return 7200 // デフォルト2時間
	}
	return lifetime
}
