package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config はアプリケーション全体の設定を保持する構造体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Admin    AdminConfig    `mapstructure:"admin"`
	Database DatabaseConfig `mapstructure:"database"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	CORS     CORSConfig     `mapstructure:"cors"`
}

// ServerConfig はサーバー設定
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig はデータベース設定
type DatabaseConfig struct {
	Shards []ShardConfig `mapstructure:"shards"`
}

// ShardConfig は各シャードの設定
type ShardConfig struct {
	ID                    int           `mapstructure:"id"`
	Driver                string        `mapstructure:"driver"`
	Host                  string        `mapstructure:"host"`
	Port                  int           `mapstructure:"port"`
	Name                  string        `mapstructure:"name"`
	User                  string        `mapstructure:"user"`
	Password              string        `mapstructure:"password"`
	DSN                   string        `mapstructure:"dsn"` // SQLite用のDSN
	MaxConnections        int           `mapstructure:"max_connections"`
	MaxIdleConnections    int           `mapstructure:"max_idle_connections"`
	ConnectionMaxLifetime time.Duration `mapstructure:"connection_max_lifetime"`

	// Writer/Reader分離用の設定
	WriterDSN    string   `mapstructure:"writer_dsn"`    // Writer接続用DSN
	ReaderDSNs   []string `mapstructure:"reader_dsns"`   // Reader接続用DSNリスト
	ReaderPolicy string   `mapstructure:"reader_policy"` // "random" or "round_robin"
}

// LoggingConfig はロギング設定
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// CORSConfig はCORS設定
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

// AdminConfig は管理画面設定
type AdminConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	Auth         AuthConfig    `mapstructure:"auth"`
	Session      SessionConfig `mapstructure:"session"`
}

// AuthConfig は認証設定
type AuthConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// SessionConfig はセッション設定
type SessionConfig struct {
	Lifetime int `mapstructure:"lifetime"`
}

// Load は指定された環境の設定ファイルを読み込む
// 環境変数 APP_ENV が設定されていない場合は "develop" がデフォルト
func Load() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "develop"
	}

	viper.SetConfigName(env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config") // サーバー実行ディレクトリから見た相対パス
	viper.AddConfigPath("../../config")
	viper.AddConfigPath("./config")

	// 環境変数の自動マッピング
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 環境変数でパスワードを上書き（セキュリティ向上）
	for i := range cfg.Database.Shards {
		shard := &cfg.Database.Shards[i]
		envKey := fmt.Sprintf("DB_PASSWORD_SHARD%d", shard.ID)
		if envPassword := os.Getenv(envKey); envPassword != "" {
			shard.Password = envPassword
		}
	}

	return &cfg, nil
}

// GetDSN はPostgreSQL/MySQL用のDSN文字列を生成する
func (s *ShardConfig) GetDSN() string {
	if s.DSN != "" {
		// SQLiteの場合はそのまま返す
		return s.DSN
	}

	// PostgreSQL/MySQL用のDSN生成
	switch s.Driver {
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			s.Host, s.Port, s.User, s.Password, s.Name)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			s.User, s.Password, s.Host, s.Port, s.Name)
	default:
		return ""
	}
}

// GetWriterDSN はWriter接続用DSNを取得
func (s *ShardConfig) GetWriterDSN() string {
	if s.WriterDSN != "" {
		return s.WriterDSN
	}
	// 後方互換性: 既存のDSNをWriterとして使用
	return s.GetDSN()
}

// GetReaderDSNs はReader接続用DSNリストを取得
func (s *ShardConfig) GetReaderDSNs() []string {
	if len(s.ReaderDSNs) > 0 {
		return s.ReaderDSNs
	}
	// 後方互換性: Writerと同じDSNをReaderとして使用（開発環境用）
	return []string{s.GetWriterDSN()}
}
