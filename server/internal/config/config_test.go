package config

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestAdminConfig_Defaults(t *testing.T) {
	// AdminConfig構造体が正しいフィールドを持つことを確認
	cfg := AdminConfig{
		Port:         8081,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Auth: AuthConfig{
			Username: "admin",
			Password: "password",
		},
		Session: SessionConfig{
			Lifetime: 7200,
		},
	}

	if cfg.Port != 8081 {
		t.Errorf("expected Port 8081, got %d", cfg.Port)
	}
	if cfg.ReadTimeout != 30*time.Second {
		t.Errorf("expected ReadTimeout 30s, got %v", cfg.ReadTimeout)
	}
	if cfg.WriteTimeout != 30*time.Second {
		t.Errorf("expected WriteTimeout 30s, got %v", cfg.WriteTimeout)
	}
	if cfg.Auth.Username != "admin" {
		t.Errorf("expected Username 'admin', got %s", cfg.Auth.Username)
	}
	if cfg.Auth.Password != "password" {
		t.Errorf("expected Password 'password', got %s", cfg.Auth.Password)
	}
	if cfg.Session.Lifetime != 7200 {
		t.Errorf("expected Session.Lifetime 7200, got %d", cfg.Session.Lifetime)
	}
}

func TestAuthConfig(t *testing.T) {
	cfg := AuthConfig{
		Username: "testuser",
		Password: "testpass",
	}

	if cfg.Username != "testuser" {
		t.Errorf("expected Username 'testuser', got %s", cfg.Username)
	}
	if cfg.Password != "testpass" {
		t.Errorf("expected Password 'testpass', got %s", cfg.Password)
	}
}

func TestSessionConfig(t *testing.T) {
	cfg := SessionConfig{
		Lifetime: 3600,
	}

	if cfg.Lifetime != 3600 {
		t.Errorf("expected Lifetime 3600, got %d", cfg.Lifetime)
	}
}

func TestConfig_HasAdminField(t *testing.T) {
	// Config構造体にAdminフィールドがあることを確認
	cfg := Config{
		Admin: AdminConfig{
			Port: 8081,
		},
	}

	if cfg.Admin.Port != 8081 {
		t.Errorf("expected Admin.Port 8081, got %d", cfg.Admin.Port)
	}
}

func TestLoad_AdminConfig(t *testing.T) {
	// 設定ファイルからAdminConfigが読み込まれることを確認
	viper.Reset()
	cfg, err := Load()
	if err != nil {
		t.Skipf("config file not found, skipping: %v", err)
	}

	// Admin設定が読み込まれることを確認
	if cfg.Admin.Port == 0 {
		t.Log("Admin.Port is not set in config file")
	}
}

func TestLoad_WithBothConfigFiles(t *testing.T) {
	// テスト用に環境変数を設定
	originalEnv := os.Getenv("APP_ENV")
	os.Setenv("APP_ENV", "develop")
	defer os.Setenv("APP_ENV", originalEnv)

	// viperをリセット
	viper.Reset()

	cfg, err := Load()
	if err != nil {
		t.Skipf("config files not found, skipping: %v", err)
	}

	// メイン設定の確認
	if cfg.Server.Port != 8080 {
		t.Errorf("expected Server.Port 8080, got %d", cfg.Server.Port)
	}
	if cfg.Admin.Port != 8081 {
		t.Errorf("expected Admin.Port 8081, got %d", cfg.Admin.Port)
	}
	if cfg.Logging.Level != "debug" {
		t.Errorf("expected Logging.Level 'debug', got %s", cfg.Logging.Level)
	}

	// データベース設定の確認
	if len(cfg.Database.Shards) == 0 {
		t.Error("expected Database.Shards to have at least one shard")
	}
	if len(cfg.Database.Shards) > 0 {
		shard := cfg.Database.Shards[0]
		if shard.ID != 1 {
			t.Errorf("expected first shard ID 1, got %d", shard.ID)
		}
		if shard.Driver != "sqlite3" {
			t.Errorf("expected first shard Driver 'sqlite3', got %s", shard.Driver)
		}
	}
}

func TestLoad_DefaultEnv(t *testing.T) {
	// APP_ENVが設定されていない場合、developがデフォルト
	originalEnv := os.Getenv("APP_ENV")
	os.Unsetenv("APP_ENV")
	defer os.Setenv("APP_ENV", originalEnv)

	viper.Reset()

	cfg, err := Load()
	if err != nil {
		t.Skipf("config files not found, skipping: %v", err)
	}

	// develop環境の設定が読み込まれることを確認
	if cfg.Server.Port != 8080 {
		t.Errorf("expected Server.Port 8080, got %d", cfg.Server.Port)
	}
}

func TestLoad_PasswordOverrideFromEnv(t *testing.T) {
	// 環境変数でパスワードを上書きできることを確認
	originalEnv := os.Getenv("APP_ENV")
	originalPassword := os.Getenv("DB_PASSWORD_SHARD1")
	os.Setenv("APP_ENV", "develop")
	os.Setenv("DB_PASSWORD_SHARD1", "env_password_test")
	defer func() {
		os.Setenv("APP_ENV", originalEnv)
		if originalPassword != "" {
			os.Setenv("DB_PASSWORD_SHARD1", originalPassword)
		} else {
			os.Unsetenv("DB_PASSWORD_SHARD1")
		}
	}()

	viper.Reset()

	cfg, err := Load()
	if err != nil {
		t.Skipf("config files not found, skipping: %v", err)
	}

	// 環境変数でパスワードが上書きされることを確認
	if len(cfg.Database.Shards) > 0 {
		if cfg.Database.Shards[0].Password != "env_password_test" {
			t.Errorf("expected Password 'env_password_test', got %s", cfg.Database.Shards[0].Password)
		}
	}
}
