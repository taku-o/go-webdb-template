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
		t.Fatalf("config file not found: %v", err)
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
		t.Fatalf("config files not found: %v", err)
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
		if master.Driver != "postgres" {
			t.Errorf("expected first master Driver 'postgres', got %s", master.Driver)
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
		t.Fatalf("config files not found: %v", err)
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
		t.Fatalf("config files not found: %v", err)
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
				Driver: "postgres",
				Host:   "localhost",
				Port:   5432,
				User:   "webdb",
				Name:   "webdb_master",
			},
		},
		Sharding: ShardingGroupConfig{
			Databases: []ShardConfig{
				{
					ID:         1,
					Driver:     "postgres",
					Host:       "localhost",
					Port:       5433,
					User:       "webdb",
					Name:       "webdb_sharding_1",
					TableRange: [2]int{0, 7},
				},
				{
					ID:         2,
					Driver:     "postgres",
					Host:       "localhost",
					Port:       5434,
					User:       "webdb",
					Name:       "webdb_sharding_2",
					TableRange: [2]int{8, 15},
				},
			},
			Tables: []ShardingTableConfig{
				{
					Name:        "dm_users",
					SuffixCount: 32,
				},
				{
					Name:        "dm_posts",
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
	if cfg.Master[0].Driver != "postgres" {
		t.Errorf("expected master Driver 'postgres', got %s", cfg.Master[0].Driver)
	}
	if cfg.Master[0].Host != "localhost" {
		t.Errorf("expected master Host 'localhost', got %s", cfg.Master[0].Host)
	}
	if cfg.Master[0].Port != 5432 {
		t.Errorf("expected master Port 5432, got %d", cfg.Master[0].Port)
	}
	if cfg.Master[0].Name != "webdb_master" {
		t.Errorf("expected master Name 'webdb_master', got %s", cfg.Master[0].Name)
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
	if cfg.Sharding.Tables[0].Name != "dm_users" {
		t.Errorf("expected table name 'dm_users', got %s", cfg.Sharding.Tables[0].Name)
	}
	if cfg.Sharding.Tables[0].SuffixCount != 32 {
		t.Errorf("expected suffix count 32, got %d", cfg.Sharding.Tables[0].SuffixCount)
	}
}

func TestShardConfig_TableRange(t *testing.T) {
	// ShardConfigにTableRangeフィールドがあることを確認
	cfg := ShardConfig{
		ID:         1,
		Driver:     "postgres",
		Host:       "localhost",
		Port:       5433,
		User:       "webdb",
		Name:       "webdb_sharding_1",
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
			{ID: 1, Driver: "postgres"},
		},
		Groups: DatabaseGroupsConfig{
			Master: []ShardConfig{
				{ID: 1, Driver: "postgres"},
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
		t.Fatalf("config files not found: %v", err)
	}

	// Masterグループの確認
	if len(cfg.Database.Groups.Master) == 0 {
		t.Error("expected at least one master database")
	} else {
		master := cfg.Database.Groups.Master[0]
		if master.ID != 1 {
			t.Errorf("expected master ID 1, got %d", master.ID)
		}
		if master.Driver != "postgres" {
			t.Errorf("expected master Driver 'postgres', got %s", master.Driver)
		}
		if master.Host != "localhost" {
			t.Errorf("expected master Host 'localhost', got %s", master.Host)
		}
		if master.Port != 5432 {
			t.Errorf("expected master Port 5432, got %d", master.Port)
		}
		if master.Name != "webdb_master" {
			t.Errorf("expected master Name 'webdb_master', got %s", master.Name)
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
		if cfg.Database.Groups.Sharding.Tables[0].Name != "dm_users" {
			t.Errorf("expected first table 'dm_users', got %s", cfg.Database.Groups.Sharding.Tables[0].Name)
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
		StorageType:       "auto",
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
	if cfg.StorageType != "auto" {
		t.Errorf("expected StorageType 'auto', got %s", cfg.StorageType)
	}
}

// タスク2.3: RateLimitConfigのStorageTypeフィールドテスト
func TestRateLimitConfig_StorageType(t *testing.T) {
	tests := []struct {
		name        string
		storageType string
	}{
		{"auto", "auto"},
		{"memory", "memory"},
		{"redis", "redis"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := RateLimitConfig{
				Enabled:     true,
				StorageType: tt.storageType,
			}
			if cfg.StorageType != tt.storageType {
				t.Errorf("expected StorageType '%s', got '%s'", tt.storageType, cfg.StorageType)
			}
		})
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
			JobQueue: RedisSingleConfig{
				Addr: "localhost:6379",
			},
			Default: RedisDefaultConfig{
				Cluster: RedisClusterConfig{
					Addrs: []string{"host1:6379", "host2:6379"},
				},
			},
		},
	}

	// JobQueue設定の確認
	if cfg.Redis.JobQueue.Addr != "localhost:6379" {
		t.Errorf("expected 'localhost:6379', got %s", cfg.Redis.JobQueue.Addr)
	}

	// Default設定の確認
	if len(cfg.Redis.Default.Cluster.Addrs) != 2 {
		t.Errorf("expected 2 addresses, got %d", len(cfg.Redis.Default.Cluster.Addrs))
	}
	if cfg.Redis.Default.Cluster.Addrs[0] != "host1:6379" {
		t.Errorf("expected 'host1:6379', got %s", cfg.Redis.Default.Cluster.Addrs[0])
	}
}

// タスク2.1: RedisSingleConfig構造体のテスト
func TestRedisSingleConfig_Structure(t *testing.T) {
	cfg := RedisSingleConfig{
		Addr: "redis.example.com:6379",
	}

	if cfg.Addr != "redis.example.com:6379" {
		t.Errorf("expected 'redis.example.com:6379', got %s", cfg.Addr)
	}
}

// タスク2.1: RedisDefaultConfig構造体のテスト
func TestRedisDefaultConfig_Structure(t *testing.T) {
	cfg := RedisDefaultConfig{
		Cluster: RedisClusterConfig{
			Addrs: []string{"node1:6379", "node2:6379", "node3:6379"},
		},
	}

	if len(cfg.Cluster.Addrs) != 3 {
		t.Errorf("expected 3 addresses, got %d", len(cfg.Cluster.Addrs))
	}
}

// タスク2.1: ConfigにCacheServerフィールドがあることを確認
func TestConfig_HasCacheServerField(t *testing.T) {
	cfg := Config{
		CacheServer: CacheServerConfig{
			Redis: RedisConfig{
				JobQueue: RedisSingleConfig{
					Addr: "localhost:6379",
				},
				Default: RedisDefaultConfig{
					Cluster: RedisClusterConfig{
						Addrs: []string{},
					},
				},
			},
		},
	}

	if cfg.CacheServer.Redis.JobQueue.Addr == "" {
		t.Error("expected CacheServer.Redis.JobQueue.Addr to be initialized")
	}
	if cfg.CacheServer.Redis.Default.Cluster.Addrs == nil {
		t.Error("expected CacheServer.Redis.Default.Cluster.Addrs to be initialized")
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
		t.Fatalf("config files not found: %v", err)
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
		t.Fatalf("config files not found: %v", err)
	}

	// ジョブキュー用Redis設定の確認
	if cfg.CacheServer.Redis.JobQueue.Addr != "localhost:6379" {
		t.Errorf("expected CacheServer.Redis.JobQueue.Addr 'localhost:6379', got %s", cfg.CacheServer.Redis.JobQueue.Addr)
	}

	// 開発環境ではデフォルト用Redis Clusterのアドレスが空であることを確認
	if len(cfg.CacheServer.Redis.Default.Cluster.Addrs) != 0 {
		t.Errorf("expected CacheServer.Redis.Default.Cluster.Addrs to be empty, got %v", cfg.CacheServer.Redis.Default.Cluster.Addrs)
	}
}

// タスク1.1: RedisClusterConfig構造体に接続オプションフィールドがあることを確認
func TestRedisClusterConfig_ConnectionOptions(t *testing.T) {
	cfg := RedisClusterConfig{
		Addrs:           []string{"host1:6379", "host2:6379"},
		MaxRetries:      2,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
		DialTimeout:     5 * time.Second,
		ReadTimeout:     3 * time.Second,
		PoolSize:        10,
		PoolTimeout:     4 * time.Second,
	}

	if len(cfg.Addrs) != 2 {
		t.Errorf("expected 2 addresses, got %d", len(cfg.Addrs))
	}
	if cfg.MaxRetries != 2 {
		t.Errorf("expected MaxRetries 2, got %d", cfg.MaxRetries)
	}
	if cfg.MinRetryBackoff != 8*time.Millisecond {
		t.Errorf("expected MinRetryBackoff 8ms, got %v", cfg.MinRetryBackoff)
	}
	if cfg.MaxRetryBackoff != 512*time.Millisecond {
		t.Errorf("expected MaxRetryBackoff 512ms, got %v", cfg.MaxRetryBackoff)
	}
	if cfg.DialTimeout != 5*time.Second {
		t.Errorf("expected DialTimeout 5s, got %v", cfg.DialTimeout)
	}
	if cfg.ReadTimeout != 3*time.Second {
		t.Errorf("expected ReadTimeout 3s, got %v", cfg.ReadTimeout)
	}
	if cfg.PoolSize != 10 {
		t.Errorf("expected PoolSize 10, got %d", cfg.PoolSize)
	}
	if cfg.PoolTimeout != 4*time.Second {
		t.Errorf("expected PoolTimeout 4s, got %v", cfg.PoolTimeout)
	}
}

// タスク1.2: RedisSingleConfig構造体に接続オプションフィールドがあることを確認
func TestRedisSingleConfig_ConnectionOptions(t *testing.T) {
	cfg := RedisSingleConfig{
		Addr:            "localhost:6379",
		MaxRetries:      2,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
		DialTimeout:     5 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		PoolSize:        10,
		PoolTimeout:     4 * time.Second,
	}

	if cfg.Addr != "localhost:6379" {
		t.Errorf("expected Addr 'localhost:6379', got %s", cfg.Addr)
	}
	if cfg.MaxRetries != 2 {
		t.Errorf("expected MaxRetries 2, got %d", cfg.MaxRetries)
	}
	if cfg.MinRetryBackoff != 8*time.Millisecond {
		t.Errorf("expected MinRetryBackoff 8ms, got %v", cfg.MinRetryBackoff)
	}
	if cfg.MaxRetryBackoff != 512*time.Millisecond {
		t.Errorf("expected MaxRetryBackoff 512ms, got %v", cfg.MaxRetryBackoff)
	}
	if cfg.DialTimeout != 5*time.Second {
		t.Errorf("expected DialTimeout 5s, got %v", cfg.DialTimeout)
	}
	if cfg.ReadTimeout != 3*time.Second {
		t.Errorf("expected ReadTimeout 3s, got %v", cfg.ReadTimeout)
	}
	if cfg.WriteTimeout != 3*time.Second {
		t.Errorf("expected WriteTimeout 3s, got %v", cfg.WriteTimeout)
	}
	if cfg.PoolSize != 10 {
		t.Errorf("expected PoolSize 10, got %d", cfg.PoolSize)
	}
	if cfg.PoolTimeout != 4*time.Second {
		t.Errorf("expected PoolTimeout 4s, got %v", cfg.PoolTimeout)
	}
}

// タスク1.1: UploadConfig構造体のテスト
func TestUploadConfig_Structure(t *testing.T) {
	// UploadConfig構造体が正しいフィールドを持つことを確認
	cfg := UploadConfig{
		BasePath:          "/api/upload/dm_movie",
		MaxFileSize:       2147483648,
		AllowedExtensions: []string{"mp4"},
		Storage: StorageConfig{
			Type: "local",
			Local: LocalStorageConfig{
				Path: "./uploads",
			},
			S3: S3StorageConfig{
				Bucket: "test-bucket",
				Region: "ap-northeast-1",
			},
		},
	}

	if cfg.BasePath != "/api/upload/dm_movie" {
		t.Errorf("expected BasePath '/api/upload/dm_movie', got %s", cfg.BasePath)
	}
	if cfg.MaxFileSize != 2147483648 {
		t.Errorf("expected MaxFileSize 2147483648, got %d", cfg.MaxFileSize)
	}
	if len(cfg.AllowedExtensions) != 1 {
		t.Errorf("expected 1 allowed extension, got %d", len(cfg.AllowedExtensions))
	}
	if cfg.AllowedExtensions[0] != "mp4" {
		t.Errorf("expected 'mp4', got %s", cfg.AllowedExtensions[0])
	}
	if cfg.Storage.Type != "local" {
		t.Errorf("expected Storage.Type 'local', got %s", cfg.Storage.Type)
	}
	if cfg.Storage.Local.Path != "./uploads" {
		t.Errorf("expected Storage.Local.Path './uploads', got %s", cfg.Storage.Local.Path)
	}
	if cfg.Storage.S3.Bucket != "test-bucket" {
		t.Errorf("expected Storage.S3.Bucket 'test-bucket', got %s", cfg.Storage.S3.Bucket)
	}
	if cfg.Storage.S3.Region != "ap-northeast-1" {
		t.Errorf("expected Storage.S3.Region 'ap-northeast-1', got %s", cfg.Storage.S3.Region)
	}
}

// タスク1.1: StorageConfig構造体のテスト
func TestStorageConfig_Structure(t *testing.T) {
	cfg := StorageConfig{
		Type: "s3",
		Local: LocalStorageConfig{
			Path: "./uploads",
		},
		S3: S3StorageConfig{
			Bucket: "my-bucket",
			Region: "us-east-1",
		},
	}

	if cfg.Type != "s3" {
		t.Errorf("expected Type 's3', got %s", cfg.Type)
	}
	if cfg.Local.Path != "./uploads" {
		t.Errorf("expected Local.Path './uploads', got %s", cfg.Local.Path)
	}
	if cfg.S3.Bucket != "my-bucket" {
		t.Errorf("expected S3.Bucket 'my-bucket', got %s", cfg.S3.Bucket)
	}
	if cfg.S3.Region != "us-east-1" {
		t.Errorf("expected S3.Region 'us-east-1', got %s", cfg.S3.Region)
	}
}

// タスク1.1: LocalStorageConfig構造体のテスト
func TestLocalStorageConfig_Structure(t *testing.T) {
	cfg := LocalStorageConfig{
		Path: "/var/uploads",
	}

	if cfg.Path != "/var/uploads" {
		t.Errorf("expected Path '/var/uploads', got %s", cfg.Path)
	}
}

// タスク1.1: S3StorageConfig構造体のテスト
func TestS3StorageConfig_Structure(t *testing.T) {
	cfg := S3StorageConfig{
		Bucket: "production-bucket",
		Region: "ap-northeast-1",
	}

	if cfg.Bucket != "production-bucket" {
		t.Errorf("expected Bucket 'production-bucket', got %s", cfg.Bucket)
	}
	if cfg.Region != "ap-northeast-1" {
		t.Errorf("expected Region 'ap-northeast-1', got %s", cfg.Region)
	}
}

// タスク1.1: ConfigにUploadフィールドがあることを確認
func TestConfig_HasUploadField(t *testing.T) {
	cfg := Config{
		Upload: UploadConfig{
			BasePath: "/api/upload/dm_movie",
		},
	}

	if cfg.Upload.BasePath != "/api/upload/dm_movie" {
		t.Errorf("expected Upload.BasePath '/api/upload/dm_movie', got %s", cfg.Upload.BasePath)
	}
}

// タスク1.1: LoggingConfig構造体にMailLogOutputDirフィールドがあることを確認
func TestLoggingConfig_HasMailLogOutputDirField(t *testing.T) {
	cfg := LoggingConfig{
		Level:            "debug",
		Format:           "json",
		Output:           "stdout",
		OutputDir:        "logs",
		SQLLogEnabled:    true,
		SQLLogOutputDir:  "logs",
		MailLogOutputDir: "logs", // 新規追加フィールド
	}

	if cfg.MailLogOutputDir != "logs" {
		t.Errorf("expected MailLogOutputDir 'logs', got %s", cfg.MailLogOutputDir)
	}
}

// タスク1.1: Load関数でMailLogOutputDirのデフォルト値が設定されることを確認
func TestLoad_MailLogOutputDirDefault(t *testing.T) {
	originalEnv := os.Getenv("APP_ENV")
	os.Setenv("APP_ENV", "develop")
	defer os.Setenv("APP_ENV", originalEnv)

	viper.Reset()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("config files not found: %v", err)
	}

	// MailLogOutputDirが設定されていない場合はOutputDirをデフォルトとして使用
	if cfg.Logging.MailLogOutputDir == "" {
		t.Error("expected MailLogOutputDir to have a default value")
	}
	// デフォルト値はOutputDirと同じであるべき
	if cfg.Logging.MailLogOutputDir != cfg.Logging.OutputDir {
		t.Errorf("expected MailLogOutputDir to default to OutputDir (%s), got %s",
			cfg.Logging.OutputDir, cfg.Logging.MailLogOutputDir)
	}
}

// タスク1.2: 設定ファイルからUploadConfigが読み込まれることを確認
func TestLoad_UploadConfig(t *testing.T) {
	originalEnv := os.Getenv("APP_ENV")
	os.Setenv("APP_ENV", "develop")
	defer os.Setenv("APP_ENV", originalEnv)

	viper.Reset()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("config files not found: %v", err)
	}

	// アップロード設定が読み込まれることを確認
	if cfg.Upload.BasePath != "/api/upload/dm_movie" {
		t.Errorf("expected Upload.BasePath '/api/upload/dm_movie', got %s", cfg.Upload.BasePath)
	}
	if cfg.Upload.MaxFileSize != 2147483648 {
		t.Errorf("expected Upload.MaxFileSize 2147483648, got %d", cfg.Upload.MaxFileSize)
	}
	if len(cfg.Upload.AllowedExtensions) != 1 {
		t.Errorf("expected 1 allowed extension, got %d", len(cfg.Upload.AllowedExtensions))
	}
	if len(cfg.Upload.AllowedExtensions) > 0 && cfg.Upload.AllowedExtensions[0] != "mp4" {
		t.Errorf("expected 'mp4', got %s", cfg.Upload.AllowedExtensions[0])
	}
	if cfg.Upload.Storage.Type != "local" {
		t.Errorf("expected Storage.Type 'local', got %s", cfg.Upload.Storage.Type)
	}
	if cfg.Upload.Storage.Local.Path != "./uploads" {
		t.Errorf("expected Storage.Local.Path './uploads', got %s", cfg.Upload.Storage.Local.Path)
	}
}

// タスク1.1: GetDSN()メソッドのテスト - PostgreSQL用DSN生成
func TestShardConfig_GetDSN_Postgres(t *testing.T) {
	cfg := ShardConfig{
		Driver:   "postgres",
		Host:     "localhost",
		Port:     5432,
		User:     "webdb",
		Password: "webdb",
		Name:     "webdb_master",
	}

	dsn := cfg.GetDSN()
	expected := "host=localhost port=5432 user=webdb password=webdb dbname=webdb_master sslmode=disable"
	if dsn != expected {
		t.Errorf("expected DSN '%s', got '%s'", expected, dsn)
	}
}

// タスク1.1: GetDSN()メソッドのテスト - MySQL用DSN生成（charset=utf8mb4&loc=Local追加）
func TestShardConfig_GetDSN_MySQL(t *testing.T) {
	cfg := ShardConfig{
		Driver:   "mysql",
		Host:     "localhost",
		Port:     3306,
		User:     "webdb",
		Password: "webdb",
		Name:     "webdb_master",
	}

	dsn := cfg.GetDSN()
	// MySQLのDSNには charset=utf8mb4、parseTime=true、loc=Local が含まれる必要がある
	expected := "webdb:webdb@tcp(localhost:3306)/webdb_master?charset=utf8mb4&parseTime=true&loc=Local"
	if dsn != expected {
		t.Errorf("expected DSN '%s', got '%s'", expected, dsn)
	}
}

// タスク1.1: GetDSN()メソッドのテスト - DSNが直接指定されている場合
func TestShardConfig_GetDSN_DirectDSN(t *testing.T) {
	cfg := ShardConfig{
		DSN: "custom-dsn-string",
	}

	dsn := cfg.GetDSN()
	if dsn != "custom-dsn-string" {
		t.Errorf("expected DSN 'custom-dsn-string', got '%s'", dsn)
	}
}

// タスク1.1: GetDSN()メソッドのテスト - 不明なドライバー
func TestShardConfig_GetDSN_UnknownDriver(t *testing.T) {
	cfg := ShardConfig{
		Driver: "unknown",
		Host:   "localhost",
		Port:   5432,
	}

	dsn := cfg.GetDSN()
	if dsn != "" {
		t.Errorf("expected empty DSN for unknown driver, got '%s'", dsn)
	}
}

// タスク3.1: Load()関数でDB_TYPEがpostgresの場合、database.yamlが読み込まれる
func TestLoad_DBType_PostgreSQL(t *testing.T) {
	originalEnv := os.Getenv("APP_ENV")
	os.Setenv("APP_ENV", "develop")
	defer os.Setenv("APP_ENV", originalEnv)

	viper.Reset()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("config files not found: %v", err)
	}

	// DB_TYPE: postgresの場合、database.yamlが読み込まれる
	// masterのドライバーがpostgresであることを確認
	if len(cfg.Database.Groups.Master) == 0 {
		t.Fatal("expected at least one master database")
	}
	if cfg.Database.Groups.Master[0].Driver != "postgres" {
		t.Errorf("expected driver 'postgres', got '%s'", cfg.Database.Groups.Master[0].Driver)
	}
}
