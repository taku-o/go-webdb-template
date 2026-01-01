package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config はアプリケーション全体の設定を保持する構造体
type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Admin       AdminConfig       `mapstructure:"admin"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Logging     LoggingConfig     `mapstructure:"logging"`
	CORS        CORSConfig        `mapstructure:"cors"`
	API         APIConfig         `mapstructure:"api"`
	CacheServer CacheServerConfig `mapstructure:"cache_server"` // キャッシュサーバー設定
	Upload      UploadConfig      `mapstructure:"upload"`       // アップロード設定
}

// CacheServerConfig はキャッシュサーバー設定
type CacheServerConfig struct {
	Redis RedisConfig `mapstructure:"redis"`
}

// RedisConfig はRedis設定
type RedisConfig struct {
	Cluster RedisClusterConfig `mapstructure:"cluster"`
}

// RedisClusterConfig はRedis Cluster設定
type RedisClusterConfig struct {
	Addrs []string `mapstructure:"addrs"` // Redis Clusterのアドレスリスト
}

// ServerConfig はサーバー設定
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// DatabaseConfig はデータベース設定
type DatabaseConfig struct {
	// 後方互換性のため残す（非推奨）
	Shards []ShardConfig `mapstructure:"shards"`

	// 新規: データベースグループ
	Groups DatabaseGroupsConfig `mapstructure:"groups"`
}

// DatabaseGroupsConfig はデータベースグループ設定
type DatabaseGroupsConfig struct {
	Master   []ShardConfig       `mapstructure:"master"`
	Sharding ShardingGroupConfig `mapstructure:"sharding"`
}

// ShardingGroupConfig はshardingグループの設定
type ShardingGroupConfig struct {
	Databases []ShardConfig         `mapstructure:"databases"`
	Tables    []ShardingTableConfig `mapstructure:"tables"`
}

// ShardingTableConfig はshardingグループのテーブル定義
type ShardingTableConfig struct {
	Name        string `mapstructure:"name"`         // テーブル名（例: "users"）
	SuffixCount int    `mapstructure:"suffix_count"` // 分割数（例: 32）
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

	// shardingグループ用: テーブル番号範囲 [min, max]
	TableRange [2]int `mapstructure:"table_range"`
}

// LoggingConfig はロギング設定
type LoggingConfig struct {
	Level           string `mapstructure:"level"`
	Format          string `mapstructure:"format"`
	Output          string `mapstructure:"output"`
	OutputDir       string `mapstructure:"output_dir"`
	SQLLogEnabled   bool   `mapstructure:"sql_log_enabled"`    // SQLログの有効/無効（オプション）
	SQLLogOutputDir string `mapstructure:"sql_log_output_dir"` // SQLログ出力先ディレクトリ（オプション）
}

// CORSConfig はCORS設定
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
	ExposeHeaders  []string `mapstructure:"expose_headers"`
}

// APIConfig はAPIキー設定
type APIConfig struct {
	CurrentVersion     string          `mapstructure:"current_version"`
	PublicKey          string          `mapstructure:"public_key"`
	SecretKey          string          `mapstructure:"secret_key"`
	InvalidVersions    []string        `mapstructure:"invalid_versions"`
	Auth0IssuerBaseURL string          `mapstructure:"auth0_issuer_base_url"` // Auth0のIssuer Base URL
	RateLimit          RateLimitConfig `mapstructure:"rate_limit"`            // レートリミット設定
}

// RateLimitConfig はレートリミット設定
type RateLimitConfig struct {
	Enabled           bool `mapstructure:"enabled"`
	RequestsPerMinute int  `mapstructure:"requests_per_minute"`
	RequestsPerHour   int  `mapstructure:"requests_per_hour"` // オプション
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

// UploadConfig はアップロード機能の設定
type UploadConfig struct {
	BasePath          string        `mapstructure:"base_path"`          // TUSエンドポイントのベースパス
	MaxFileSize       int64         `mapstructure:"max_file_size"`      // 最大ファイルサイズ
	AllowedExtensions []string      `mapstructure:"allowed_extensions"` // 許可された拡張子リスト
	Storage           StorageConfig `mapstructure:"storage"`            // ストレージ設定
}

// StorageConfig はストレージ設定
type StorageConfig struct {
	Type  string             `mapstructure:"type"` // ストレージタイプ（"local" or "s3"）
	Local LocalStorageConfig `mapstructure:"local"`
	S3    S3StorageConfig    `mapstructure:"s3"`
}

// LocalStorageConfig はローカルストレージ設定
type LocalStorageConfig struct {
	Path string `mapstructure:"path"` // ローカル保存パス
}

// S3StorageConfig はS3ストレージ設定
type S3StorageConfig struct {
	Bucket string `mapstructure:"bucket"` // S3バケット名
	Region string `mapstructure:"region"` // AWSリージョン
}

// Load は指定された環境の設定ファイルを読み込む
// 環境変数 APP_ENV が設定されていない場合は "develop" がデフォルト
// 設定ファイルは環境別ディレクトリ（config/{env}/）から読み込む
// メイン設定ファイル（config.yaml）とデータベース設定ファイル（database.yaml）を統合する
func Load() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "develop"
	}

	viper.SetConfigType("yaml")

	// 環境別ディレクトリのパスを追加（複数パスで実行ディレクトリの違いに対応）
	viper.AddConfigPath(fmt.Sprintf("../config/%s", env))
	viper.AddConfigPath(fmt.Sprintf("../../config/%s", env))
	viper.AddConfigPath(fmt.Sprintf("../../../config/%s", env))
	viper.AddConfigPath(fmt.Sprintf("./config/%s", env))

	// 環境変数の自動マッピング
	viper.AutomaticEnv()

	// メイン設定ファイルの読み込み
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read main config file: %w", err)
	}

	// データベース設定ファイルのマージ
	viper.SetConfigName("database")
	if err := viper.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read database config file: %w", err)
	}

	// キャッシュサーバー設定ファイルのマージ（オプショナル）
	viper.SetConfigName("cacheserver")
	if err := viper.MergeInConfig(); err != nil {
		// cacheserver.yamlが存在しない場合はエラーにしない（オプショナル）
		// 設定ファイルが見つからないエラーの場合のみ無視
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read cacheserver config file: %w", err)
		}
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

	// デフォルト値の設定
	if cfg.Logging.OutputDir == "" {
		cfg.Logging.OutputDir = "logs"
	}

	// SQLログ出力先のデフォルト値設定
	if cfg.Logging.SQLLogOutputDir == "" {
		cfg.Logging.SQLLogOutputDir = cfg.Logging.OutputDir
	}

	// SQLログ有効/無効の環境判定（設定ファイルで明示的に指定されていない場合）
	// develop/staging: true, production: false
	// 注意: boolのデフォルトはfalseなので、設定ファイルで明示的にtrueを指定する必要がある
	// 環境判定による自動有効化は行わない（設定ファイルの設定を優先）

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
