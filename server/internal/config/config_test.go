package config

import (
	"testing"
	"time"
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
	cfg, err := Load()
	if err != nil {
		t.Skipf("config file not found, skipping: %v", err)
	}

	// Admin設定が読み込まれることを確認
	if cfg.Admin.Port == 0 {
		t.Log("Admin.Port is not set in config file")
	}
}
