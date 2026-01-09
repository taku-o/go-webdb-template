package admin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/taku-o/go-webdb-template/internal/config"
)

func TestGetDatabaseConfig_Success(t *testing.T) {
	// 正常系: masterグループのデータベース設定が存在する場合、PostgreSQL設定が正しく構築されることを確認
	appConfig := &config.Config{
		Admin: config.AdminConfig{
			Port: 8081,
		},
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Master: []config.ShardConfig{
					{
						ID:       1,
						Driver:   "postgresql",
						Host:     "localhost",
						Port:     5432,
						Name:     "webdb_master",
						User:     "webdb",
						Password: "webdb",
					},
				},
			},
		},
		Logging: config.LoggingConfig{
			Level: "debug",
		},
	}

	cfg := NewConfig(appConfig)
	goadminCfg := cfg.GetGoAdminConfig()

	require.NotNil(t, goadminCfg)
	require.NotNil(t, goadminCfg.Databases)

	// PostgreSQL設定が正しく構築されていることを確認
	defaultDB, ok := goadminCfg.Databases["default"]
	require.True(t, ok, "default database configuration should exist")
	assert.Equal(t, "postgresql", defaultDB.Driver, "Driver should be postgresql")

	// DSN形式が正しく構築されていることを確認
	expectedDsn := "host=localhost port=5432 user=webdb password=webdb dbname=webdb_master sslmode=disable"
	assert.Equal(t, expectedDsn, defaultDB.Dsn, "DSN should be correctly formatted for PostgreSQL")
}

func TestGetDatabaseConfig_NoMasterGroup(t *testing.T) {
	// 異常系: masterグループのデータベース設定が存在しない場合、エラーが発生することを確認
	appConfig := &config.Config{
		Admin: config.AdminConfig{
			Port: 8081,
		},
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Master: []config.ShardConfig{}, // 空のmasterグループ
			},
		},
		Logging: config.LoggingConfig{
			Level: "debug",
		},
	}

	cfg := NewConfig(appConfig)

	// panicが発生することを確認
	assert.Panics(t, func() {
		cfg.GetGoAdminConfig()
	}, "should panic when master group is empty")
}

func TestGetDatabaseConfig_IncompleteConfig_MissingHost(t *testing.T) {
	// 異常系: 接続情報が不完全な場合（hostが空）、エラーが発生することを確認
	appConfig := &config.Config{
		Admin: config.AdminConfig{
			Port: 8081,
		},
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Master: []config.ShardConfig{
					{
						ID:       1,
						Driver:   "postgresql",
						Host:     "", // 空のhost
						Port:     5432,
						Name:     "webdb_master",
						User:     "webdb",
						Password: "webdb",
					},
				},
			},
		},
		Logging: config.LoggingConfig{
			Level: "debug",
		},
	}

	cfg := NewConfig(appConfig)

	// panicが発生することを確認
	assert.Panics(t, func() {
		cfg.GetGoAdminConfig()
	}, "should panic when host is empty")
}

func TestGetDatabaseConfig_IncompleteConfig_MissingPort(t *testing.T) {
	// 異常系: 接続情報が不完全な場合（portが0）、エラーが発生することを確認
	appConfig := &config.Config{
		Admin: config.AdminConfig{
			Port: 8081,
		},
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Master: []config.ShardConfig{
					{
						ID:       1,
						Driver:   "postgresql",
						Host:     "localhost",
						Port:     0, // portが0
						Name:     "webdb_master",
						User:     "webdb",
						Password: "webdb",
					},
				},
			},
		},
		Logging: config.LoggingConfig{
			Level: "debug",
		},
	}

	cfg := NewConfig(appConfig)

	// panicが発生することを確認
	assert.Panics(t, func() {
		cfg.GetGoAdminConfig()
	}, "should panic when port is 0")
}

func TestGetDatabaseConfig_IncompleteConfig_MissingUser(t *testing.T) {
	// 異常系: 接続情報が不完全な場合（userが空）、エラーが発生することを確認
	appConfig := &config.Config{
		Admin: config.AdminConfig{
			Port: 8081,
		},
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Master: []config.ShardConfig{
					{
						ID:       1,
						Driver:   "postgresql",
						Host:     "localhost",
						Port:     5432,
						Name:     "webdb_master",
						User:     "", // 空のuser
						Password: "webdb",
					},
				},
			},
		},
		Logging: config.LoggingConfig{
			Level: "debug",
		},
	}

	cfg := NewConfig(appConfig)

	// panicが発生することを確認
	assert.Panics(t, func() {
		cfg.GetGoAdminConfig()
	}, "should panic when user is empty")
}

func TestGetDatabaseConfig_IncompleteConfig_MissingName(t *testing.T) {
	// 異常系: 接続情報が不完全な場合（nameが空）、エラーが発生することを確認
	appConfig := &config.Config{
		Admin: config.AdminConfig{
			Port: 8081,
		},
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Master: []config.ShardConfig{
					{
						ID:       1,
						Driver:   "postgresql",
						Host:     "localhost",
						Port:     5432,
						Name:     "", // 空のname
						User:     "webdb",
						Password: "webdb",
					},
				},
			},
		},
		Logging: config.LoggingConfig{
			Level: "debug",
		},
	}

	cfg := NewConfig(appConfig)

	// panicが発生することを確認
	assert.Panics(t, func() {
		cfg.GetGoAdminConfig()
	}, "should panic when name is empty")
}

func TestGetDatabaseConfig_EmptyPassword(t *testing.T) {
	// 正常系: パスワードが空でもDSNは正しく構築されることを確認
	appConfig := &config.Config{
		Admin: config.AdminConfig{
			Port: 8081,
		},
		Database: config.DatabaseConfig{
			Groups: config.DatabaseGroupsConfig{
				Master: []config.ShardConfig{
					{
						ID:       1,
						Driver:   "postgresql",
						Host:     "localhost",
						Port:     5432,
						Name:     "webdb_master",
						User:     "webdb",
						Password: "", // 空のpassword（許容される）
					},
				},
			},
		},
		Logging: config.LoggingConfig{
			Level: "debug",
		},
	}

	cfg := NewConfig(appConfig)
	goadminCfg := cfg.GetGoAdminConfig()

	require.NotNil(t, goadminCfg)
	defaultDB, ok := goadminCfg.Databases["default"]
	require.True(t, ok)
	assert.Equal(t, "postgresql", defaultDB.Driver)

	// DSN形式が正しく構築されていることを確認（パスワードは空）
	expectedDsn := "host=localhost port=5432 user=webdb password= dbname=webdb_master sslmode=disable"
	assert.Equal(t, expectedDsn, defaultDB.Dsn)
}
