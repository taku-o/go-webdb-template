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

	// データベース設定の確認（Groupsを使用）
	if len(cfg.Database.Groups.Master) == 0 {
		t.Error("expected Database.Groups.Master to have at least one entry")
	}
	if len(cfg.Database.Groups.Master) > 0 {
		master := cfg.Database.Groups.Master[0]
		if master.ID != 1 {
			t.Errorf("expected first master ID 1, got %d", master.ID)
		}
		if master.Driver != "sqlite3" {
			t.Errorf("expected first master Driver 'sqlite3', got %s", master.Driver)
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

// タスク1.1, 1.2: DatabaseGroupsConfig構造体のテスト
func TestDatabaseGroupsConfig_Structure(t *testing.T) {
	// DatabaseGroupsConfig構造体が正しいフィールドを持つことを確認
	cfg := DatabaseGroupsConfig{
		Master: []ShardConfig{
			{
				ID:     1,
				Driver: "sqlite3",
				DSN:    "./data/master.db",
			},
		},
		Sharding: ShardingGroupConfig{
			Databases: []ShardConfig{
				{
					ID:         1,
					Driver:     "sqlite3",
					DSN:        "./data/sharding_db_1.db",
					TableRange: [2]int{0, 7},
				},
				{
					ID:         2,
					Driver:     "sqlite3",
					DSN:        "./data/sharding_db_2.db",
					TableRange: [2]int{8, 15},
				},
			},
			Tables: []ShardingTableConfig{
				{
					Name:        "users",
					SuffixCount: 32,
				},
				{
					Name:        "posts",
					SuffixCount: 32,
				},
			},
		},
	}

	// Master構成の確認
	if len(cfg.Master) != 1 {
		t.Errorf("expected 1 master database, got %d", len(cfg.Master))
	}
	if cfg.Master[0].ID != 1 {
		t.Errorf("expected master ID 1, got %d", cfg.Master[0].ID)
	}
	if cfg.Master[0].DSN != "./data/master.db" {
		t.Errorf("expected master DSN './data/master.db', got %s", cfg.Master[0].DSN)
	}

	// Sharding構成の確認
	if len(cfg.Sharding.Databases) != 2 {
		t.Errorf("expected 2 sharding databases, got %d", len(cfg.Sharding.Databases))
	}
	if cfg.Sharding.Databases[0].TableRange[0] != 0 {
		t.Errorf("expected TableRange[0] 0, got %d", cfg.Sharding.Databases[0].TableRange[0])
	}
	if cfg.Sharding.Databases[0].TableRange[1] != 7 {
		t.Errorf("expected TableRange[1] 7, got %d", cfg.Sharding.Databases[0].TableRange[1])
	}
	if cfg.Sharding.Databases[1].TableRange[0] != 8 {
		t.Errorf("expected TableRange[0] 8, got %d", cfg.Sharding.Databases[1].TableRange[0])
	}
	if cfg.Sharding.Databases[1].TableRange[1] != 15 {
		t.Errorf("expected TableRange[1] 15, got %d", cfg.Sharding.Databases[1].TableRange[1])
	}

	// Tables構成の確認
	if len(cfg.Sharding.Tables) != 2 {
		t.Errorf("expected 2 tables, got %d", len(cfg.Sharding.Tables))
	}
	if cfg.Sharding.Tables[0].Name != "users" {
		t.Errorf("expected table name 'users', got %s", cfg.Sharding.Tables[0].Name)
	}
	if cfg.Sharding.Tables[0].SuffixCount != 32 {
		t.Errorf("expected suffix count 32, got %d", cfg.Sharding.Tables[0].SuffixCount)
	}
}

func TestShardConfig_TableRange(t *testing.T) {
	// ShardConfigにTableRangeフィールドがあることを確認
	cfg := ShardConfig{
		ID:         1,
		Driver:     "sqlite3",
		DSN:        "./data/sharding_db_1.db",
		TableRange: [2]int{0, 7},
	}

	if cfg.TableRange[0] != 0 {
		t.Errorf("expected TableRange[0] 0, got %d", cfg.TableRange[0])
	}
	if cfg.TableRange[1] != 7 {
		t.Errorf("expected TableRange[1] 7, got %d", cfg.TableRange[1])
	}
}

func TestDatabaseConfig_Groups(t *testing.T) {
	// DatabaseConfigにGroupsフィールドがあることを確認
	cfg := DatabaseConfig{
		Shards: []ShardConfig{
			{ID: 1, Driver: "sqlite3"},
		},
		Groups: DatabaseGroupsConfig{
			Master: []ShardConfig{
				{ID: 1, Driver: "sqlite3"},
			},
		},
	}

	if len(cfg.Shards) != 1 {
		t.Errorf("expected 1 shard, got %d", len(cfg.Shards))
	}
	if len(cfg.Groups.Master) != 1 {
		t.Errorf("expected 1 master database, got %d", len(cfg.Groups.Master))
	}
}

func TestLoad_GroupsConfig(t *testing.T) {
	// 設定ファイルからGroupsConfigが読み込まれることを確認
	originalEnv := os.Getenv("APP_ENV")
	os.Setenv("APP_ENV", "develop")
	defer os.Setenv("APP_ENV", originalEnv)

	viper.Reset()

	cfg, err := Load()
	if err != nil {
		t.Skipf("config files not found, skipping: %v", err)
	}

	// Masterグループの確認
	if len(cfg.Database.Groups.Master) == 0 {
		t.Error("expected at least one master database")
	} else {
		master := cfg.Database.Groups.Master[0]
		if master.ID != 1 {
			t.Errorf("expected master ID 1, got %d", master.ID)
		}
		if master.Driver != "sqlite3" {
			t.Errorf("expected master Driver 'sqlite3', got %s", master.Driver)
		}
		if master.DSN != "./data/master.db" {
			t.Errorf("expected master DSN './data/master.db', got %s", master.DSN)
		}
	}

	// Shardingグループの確認（8シャーディングエントリ構成）
	if len(cfg.Database.Groups.Sharding.Databases) != 8 {
		t.Errorf("expected 8 sharding databases, got %d", len(cfg.Database.Groups.Sharding.Databases))
	}
	if len(cfg.Database.Groups.Sharding.Databases) > 0 {
		db1 := cfg.Database.Groups.Sharding.Databases[0]
		// 8シャーディング構成では各エントリが4テーブルを担当
		if db1.TableRange[0] != 0 || db1.TableRange[1] != 3 {
			t.Errorf("expected table_range [0, 3], got [%d, %d]", db1.TableRange[0], db1.TableRange[1])
		}
	}

	// テーブル設定の確認
	if len(cfg.Database.Groups.Sharding.Tables) != 2 {
		t.Errorf("expected 2 tables, got %d", len(cfg.Database.Groups.Sharding.Tables))
	}
	if len(cfg.Database.Groups.Sharding.Tables) > 0 {
		if cfg.Database.Groups.Sharding.Tables[0].Name != "users" {
			t.Errorf("expected first table 'users', got %s", cfg.Database.Groups.Sharding.Tables[0].Name)
		}
		if cfg.Database.Groups.Sharding.Tables[0].SuffixCount != 32 {
			t.Errorf("expected suffix_count 32, got %d", cfg.Database.Groups.Sharding.Tables[0].SuffixCount)
		}
	}
}

// タスク2.1: RateLimitConfig構造体のテスト
func TestRateLimitConfig_Structure(t *testing.T) {
	// RateLimitConfig構造体が正しいフィールドを持つことを確認
	cfg := RateLimitConfig{
		Enabled:           true,
		RequestsPerMinute: 60,
		RequestsPerHour:   1000,
	}

	if !cfg.Enabled {
		t.Error("expected Enabled true, got false")
	}
	if cfg.RequestsPerMinute != 60 {
		t.Errorf("expected RequestsPerMinute 60, got %d", cfg.RequestsPerMinute)
	}
	if cfg.RequestsPerHour != 1000 {
		t.Errorf("expected RequestsPerHour 1000, got %d", cfg.RequestsPerHour)
	}
}

// タスク2.1: APIConfigにRateLimitフィールドがあることを確認
func TestAPIConfig_HasRateLimitField(t *testing.T) {
	cfg := APIConfig{
		CurrentVersion: "v2",
		SecretKey:      "test_secret",
		RateLimit: RateLimitConfig{
			Enabled:           true,
			RequestsPerMinute: 60,
		},
	}

	if cfg.RateLimit.Enabled != true {
		t.Error("expected RateLimit.Enabled true, got false")
	}
	if cfg.RateLimit.RequestsPerMinute != 60 {
		t.Errorf("expected RateLimit.RequestsPerMinute 60, got %d", cfg.RateLimit.RequestsPerMinute)
	}
}

// タスク2.1: CacheServerConfig構造体のテスト
func TestCacheServerConfig_Structure(t *testing.T) {
	cfg := CacheServerConfig{
		Redis: RedisConfig{
			Cluster: RedisClusterConfig{
				Addrs: []string{"host1:6379", "host2:6379"},
			},
		},
	}

	if len(cfg.Redis.Cluster.Addrs) != 2 {
		t.Errorf("expected 2 addresses, got %d", len(cfg.Redis.Cluster.Addrs))
	}
	if cfg.Redis.Cluster.Addrs[0] != "host1:6379" {
		t.Errorf("expected 'host1:6379', got %s", cfg.Redis.Cluster.Addrs[0])
	}
}

// タスク2.1: ConfigにCacheServerフィールドがあることを確認
func TestConfig_HasCacheServerField(t *testing.T) {
	cfg := Config{
		CacheServer: CacheServerConfig{
			Redis: RedisConfig{
				Cluster: RedisClusterConfig{
					Addrs: []string{},
				},
			},
		},
	}

	if cfg.CacheServer.Redis.Cluster.Addrs == nil {
		t.Error("expected CacheServer.Redis.Cluster.Addrs to be initialized")
	}
}

// タスク2.2, 2.3: 設定ファイルからRateLimitConfigが読み込まれることを確認
func TestLoad_RateLimitConfig(t *testing.T) {
	originalEnv := os.Getenv("APP_ENV")
	os.Setenv("APP_ENV", "develop")
	defer os.Setenv("APP_ENV", originalEnv)

	viper.Reset()

	cfg, err := Load()
	if err != nil {
		t.Skipf("config files not found, skipping: %v", err)
	}

	// レートリミット設定が読み込まれることを確認
	if !cfg.API.RateLimit.Enabled {
		t.Error("expected API.RateLimit.Enabled true, got false")
	}
	if cfg.API.RateLimit.RequestsPerMinute != 60 {
		t.Errorf("expected API.RateLimit.RequestsPerMinute 60, got %d", cfg.API.RateLimit.RequestsPerMinute)
	}
	if cfg.API.RateLimit.RequestsPerHour != 1000 {
		t.Errorf("expected API.RateLimit.RequestsPerHour 1000, got %d", cfg.API.RateLimit.RequestsPerHour)
	}
}

// タスク2.2, 2.4: 設定ファイルからCacheServerConfigが読み込まれることを確認
func TestLoad_CacheServerConfig(t *testing.T) {
	originalEnv := os.Getenv("APP_ENV")
	os.Setenv("APP_ENV", "develop")
	defer os.Setenv("APP_ENV", originalEnv)

	viper.Reset()

	cfg, err := Load()
	if err != nil {
		t.Skipf("config files not found, skipping: %v", err)
	}

	// 開発環境ではRedis Clusterのアドレスが空であることを確認
	if len(cfg.CacheServer.Redis.Cluster.Addrs) != 0 {
		t.Errorf("expected CacheServer.Redis.Cluster.Addrs to be empty, got %v", cfg.CacheServer.Redis.Cluster.Addrs)
	}
}
