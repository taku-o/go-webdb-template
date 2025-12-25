package admin

import (
	"testing"

	"github.com/example/go-webdb-template/internal/config"
)

func TestNewConfig(t *testing.T) {
	appConfig := &config.Config{
		Admin: config.AdminConfig{
			Port: 8081,
			Auth: config.AuthConfig{
				Username: "admin",
				Password: "password",
			},
			Session: config.SessionConfig{
				Lifetime: 7200,
			},
		},
		Database: config.DatabaseConfig{
			Shards: []config.ShardConfig{
				{
					ID:  1,
					DSN: "test.db",
				},
			},
		},
		Logging: config.LoggingConfig{
			Level: "debug",
		},
	}

	cfg := NewConfig(appConfig)

	if cfg == nil {
		t.Error("NewConfig returned nil")
	}
}

func TestGetGoAdminConfig(t *testing.T) {
	appConfig := &config.Config{
		Admin: config.AdminConfig{
			Port: 8081,
			Session: config.SessionConfig{
				Lifetime: 3600,
			},
		},
		Database: config.DatabaseConfig{
			Shards: []config.ShardConfig{
				{
					ID:  1,
					DSN: "test.db",
				},
			},
		},
		Logging: config.LoggingConfig{
			Level: "debug",
		},
	}

	cfg := NewConfig(appConfig)
	goadminCfg := cfg.GetGoAdminConfig()

	if goadminCfg == nil {
		t.Error("GetGoAdminConfig returned nil")
	}

	if goadminCfg.SessionLifeTime != 3600 {
		t.Errorf("expected SessionLifeTime 3600, got %d", goadminCfg.SessionLifeTime)
	}

	if goadminCfg.Debug != true {
		t.Error("expected Debug to be true for debug logging level")
	}
}

func TestGetAdminPort(t *testing.T) {
	appConfig := &config.Config{
		Admin: config.AdminConfig{
			Port: 9090,
		},
	}

	cfg := NewConfig(appConfig)

	if cfg.GetAdminPort() != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.GetAdminPort())
	}
}

func TestGeneratorsMap(t *testing.T) {
	if Generators == nil {
		t.Error("Generators map is nil")
	}

	if _, ok := Generators["users"]; !ok {
		t.Error("users generator not found in Generators map")
	}

	if _, ok := Generators["posts"]; !ok {
		t.Error("posts generator not found in Generators map")
	}
}
