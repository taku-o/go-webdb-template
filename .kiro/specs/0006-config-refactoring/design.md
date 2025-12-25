# 設定ファイル分割・リファクタリング設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、設定ファイルの分割・リファクタリングの詳細設計を定義する。環境別ディレクトリ構造への移行とデータベース設定の分離を実現し、設定ファイルの管理性と保守性を向上させる。

### 1.2 設計の範囲
- ディレクトリ構造の設計
- 設定ファイル分割の設計
- 設定読み込みロジックの設計
- エラーハンドリング設計
- 後方互換性の維持設計
- テスト戦略

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
config/
├── develop.yaml              # 開発環境設定（全設定を含む）
├── production.yaml.example   # 本番環境設定テンプレート（全設定を含む）
└── staging.yaml              # ステージング環境設定（全設定を含む）
```

#### 2.1.2 変更後の構造
```
config/
├── develop/                  # 開発環境設定ディレクトリ
│   ├── config.yaml           # メイン設定（server, admin, logging, cors）
│   └── database.yaml         # データベース設定
├── production/               # 本番環境設定ディレクトリ
│   ├── config.yaml.example   # メイン設定テンプレート
│   └── database.yaml.example # データベース設定テンプレート
└── staging/                   # ステージング環境設定ディレクトリ
    ├── config.yaml           # メイン設定
    └── database.yaml         # データベース設定
```

### 2.2 設定ファイル分割の設計

#### 2.2.1 メイン設定ファイル（config.yaml）
以下の設定セクションを含む：
- `server`: サーバー設定（ポート、タイムアウト）
- `admin`: 管理画面設定（ポート、認証、セッション）
- `logging`: ロギング設定（レベル、フォーマット、出力先）
- `cors`: CORS設定（許可オリジン、メソッド、ヘッダー）

**注意**: `database`セクションは含まない。

#### 2.2.2 データベース設定ファイル（database.yaml）
以下の設定セクションを含む：
- `database.shards`: シャード設定のリスト
  - 各シャードの設定（id, driver, host, port, name, user, password, dsn, writer_dsn, reader_dsns, reader_policy, max_connections, max_idle_connections, connection_max_lifetime）

### 2.3 設定読み込みフロー

```
┌─────────────────────────────────────────────────────────────┐
│                    Load() 関数の実行                        │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              1. 環境変数 APP_ENV の取得                      │
│              (デフォルト: "develop")                         │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              2. viper の初期設定                             │
│              - SetConfigType("yaml")                        │
│              - AddConfigPath("config/{env}/")               │
│              - AutomaticEnv()                               │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              3. メイン設定ファイルの読み込み                  │
│              - SetConfigName("config")                      │
│              - ReadInConfig()                               │
│              → config/{env}/config.yaml                     │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              4. データベース設定ファイルのマージ              │
│              - SetConfigName("database")                    │
│              - MergeInConfig()                              │
│              → config/{env}/database.yaml                   │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              5. Config 構造体へのマッピング                  │
│              - Unmarshal(&cfg)                              │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              6. 環境変数によるパスワード上書き                │
│              - DB_PASSWORD_SHARD* の処理                     │
└────────────────────────┬──────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              7. Config 構造体の返却                         │
└─────────────────────────────────────────────────────────────┘
```

### 2.4 viperの複数ファイル読み込みメカニズム

viperは以下の手順で複数ファイルを読み込む：

1. **メイン設定ファイルの読み込み**
   ```go
   viper.SetConfigName("config")
   viper.SetConfigType("yaml")
   viper.AddConfigPath("config/develop")
   err := viper.ReadInConfig()
   ```

2. **データベース設定ファイルのマージ**
   ```go
   viper.SetConfigName("database")
   err := viper.MergeInConfig()
   ```

3. **統合された設定の取得**
   ```go
   var cfg Config
   err := viper.Unmarshal(&cfg)
   ```

**重要なポイント**:
- `MergeInConfig()`は既存の設定に新しい設定をマージする
- 同じキーが存在する場合、後から読み込んだ設定が優先される
- データベース設定は`database`セクションとして統合される

## 3. コンポーネント設計

### 3.1 Load()関数の設計

#### 3.1.1 関数シグネチャ
```go
func Load() (*Config, error)
```

変更なし（既存のAPIを維持）。

#### 3.1.2 実装の詳細設計

| ステップ | 処理内容 | エラーハンドリング |
|---------|---------|------------------|
| 1 | 環境変数`APP_ENV`の取得（デフォルト: "develop"） | なし（デフォルト値を使用） |
| 2 | viperの基本設定（SetConfigType, AddConfigPath） | なし |
| 3 | 環境変数の自動マッピング（AutomaticEnv） | なし |
| 4 | メイン設定ファイルの読み込み（ReadInConfig） | エラー時は`fmt.Errorf`で返却 |
| 5 | データベース設定ファイルのマージ（MergeInConfig） | エラー時は`fmt.Errorf`で返却（オプション: ファイルが存在しない場合は警告のみ） |
| 6 | Config構造体へのマッピング（Unmarshal） | エラー時は`fmt.Errorf`で返却 |
| 7 | 環境変数によるパスワード上書き | なし（既存ロジックを維持） |
| 8 | Config構造体の返却 | なし |

#### 3.1.3 実装例（擬似コード）

```go
func Load() (*Config, error) {
	// 1. 環境変数の取得
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "develop"
	}

	// 2. viperの基本設定
	viper.SetConfigType("yaml")
	
	// 環境別ディレクトリのパスを追加（複数パスで実行ディレクトリの違いに対応）
	viper.AddConfigPath(fmt.Sprintf("../config/%s", env))
	viper.AddConfigPath(fmt.Sprintf("../../config/%s", env))
	viper.AddConfigPath(fmt.Sprintf("./config/%s", env))

	// 3. 環境変数の自動マッピング
	viper.AutomaticEnv()

	// 4. メイン設定ファイルの読み込み
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read main config file: %w", err)
	}

	// 5. データベース設定ファイルのマージ
	viper.SetConfigName("database")
	if err := viper.MergeInConfig(); err != nil {
		// オプション: ファイルが存在しない場合は警告のみ（開発環境など）
		// 本番環境では必須とする場合はエラーを返す
		return nil, fmt.Errorf("failed to read database config file: %w", err)
	}

	// 6. Config構造体へのマッピング
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 7. 環境変数によるパスワード上書き（既存ロジックを維持）
	for i := range cfg.Database.Shards {
		shard := &cfg.Database.Shards[i]
		envKey := fmt.Sprintf("DB_PASSWORD_SHARD%d", shard.ID)
		if envPassword := os.Getenv(envKey); envPassword != "" {
			shard.Password = envPassword
		}
	}

	return &cfg, nil
}
```

### 3.2 設定構造体の設計

#### 3.2.1 Config構造体
変更なし（既存の構造体を維持）。

```go
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Admin    AdminConfig    `mapstructure:"admin"`
	Database DatabaseConfig `mapstructure:"database"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	CORS     CORSConfig     `mapstructure:"cors"`
}
```

#### 3.2.2 その他の設定構造体
すべて変更なし（既存の構造体を維持）。

## 4. データモデル

### 4.1 設定ファイルの構造

#### 4.1.1 メイン設定ファイル（config.yaml）の構造
```yaml
server:
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

admin:
  port: 8081
  read_timeout: 30s
  write_timeout: 30s
  auth:
    username: admin
    password: admin123
  session:
    lifetime: 7200

logging:
  level: debug
  format: json
  output: stdout

cors:
  allowed_origins:
    - http://localhost:3000
  allowed_methods:
    - GET
    - POST
    - PUT
    - DELETE
    - OPTIONS
  allowed_headers:
    - Content-Type
    - Authorization
```

#### 4.1.2 データベース設定ファイル（database.yaml）の構造
```yaml
database:
  shards:
    - id: 1
      driver: sqlite3
      dsn: ./data/shard1.db
      writer_dsn: ./data/shard1.db
      reader_dsns:
        - ./data/shard1.db
      reader_policy: random
      max_connections: 10
      max_idle_connections: 5
      connection_max_lifetime: 300s
    - id: 2
      driver: sqlite3
      dsn: ./data/shard2.db
      writer_dsn: ./data/shard2.db
      reader_dsns:
        - ./data/shard2.db
      reader_policy: random
      max_connections: 10
      max_idle_connections: 5
      connection_max_lifetime: 300s
```

### 4.2 設定の統合結果

メイン設定ファイルとデータベース設定ファイルをマージした結果、以下の構造になる：

```yaml
server:
  # ... (メイン設定ファイルから)

admin:
  # ... (メイン設定ファイルから)

database:
  # ... (データベース設定ファイルから)

logging:
  # ... (メイン設定ファイルから)

cors:
  # ... (メイン設定ファイルから)
```

この統合された構造が`Config`構造体にマッピングされる。

## 5. エラーハンドリング

### 5.1 メイン設定ファイルの読み込みエラー

**エラーケース**:
- ファイルが存在しない
- ファイルの形式が不正（YAML構文エラー）
- ファイルの読み込み権限がない

**処理**:
- `viper.ReadInConfig()`がエラーを返す
- `fmt.Errorf("failed to read main config file: %w", err)`でエラーを返却
- アプリケーションの起動を停止

### 5.2 データベース設定ファイルの読み込みエラー

**エラーケース**:
- ファイルが存在しない
- ファイルの形式が不正（YAML構文エラー）
- ファイルの読み込み権限がない

**処理方針**:
- **必須とする場合**: `viper.MergeInConfig()`がエラーを返したら、`fmt.Errorf("failed to read database config file: %w", err)`でエラーを返却
- **オプションとする場合**: ファイルが存在しない場合は警告ログを出力し、データベース設定を空の状態で続行（開発環境など）

**推奨**: 本実装では必須とする（データベース設定は必須のため）。

### 5.3 設定のマッピングエラー

**エラーケース**:
- 設定値の型が不一致
- 必須フィールドが不足

**処理**:
- `viper.Unmarshal()`がエラーを返す
- `fmt.Errorf("failed to unmarshal config: %w", err)`でエラーを返却
- アプリケーションの起動を停止

### 5.4 環境変数の処理エラー

**エラーケース**:
- 環境変数の値が不正（型変換エラーなど）

**処理**:
- 既存のロジックを維持（エラーハンドリングなし）
- 環境変数の値が不正な場合は、設定ファイルの値が使用される

## 6. テスト戦略

### 6.1 ユニットテスト

#### 6.1.1 Load()関数のテスト

**テストケース**:
1. **正常系**: メイン設定ファイルとデータベース設定ファイルの両方が存在する場合
2. **正常系**: 環境変数`APP_ENV`が設定されている場合
3. **正常系**: 環境変数`APP_ENV`が設定されていない場合（デフォルト値の使用）
4. **正常系**: 環境変数によるパスワード上書きが正常に動作する場合
5. **異常系**: メイン設定ファイルが存在しない場合
6. **異常系**: データベース設定ファイルが存在しない場合
7. **異常系**: メイン設定ファイルのYAML構文が不正な場合
8. **異常系**: データベース設定ファイルのYAML構文が不正な場合
9. **異常系**: 設定値の型が不一致の場合

**テスト実装場所**:
- `server/internal/config/config_test.go`

**テストデータ**:
- `testdata/develop/`: テスト用の設定ファイル
  - `config.yaml`: テスト用メイン設定
  - `database.yaml`: テスト用データベース設定

### 6.2 統合テスト

#### 6.2.1 設定読み込みの統合テスト

**テストケース**:
1. 実際の設定ファイルを使用した読み込みテスト
2. 複数環境（develop, staging, production）での読み込みテスト
3. 環境変数による上書きの統合テスト

**テスト実装場所**:
- `server/internal/config/config_integration_test.go`（必要に応じて）

### 6.3 E2Eテスト

#### 6.3.1 アプリケーション起動テスト

**テストケース**:
1. 新しい設定ファイル構造でアプリケーションが正常に起動することを確認
2. 既存の設定取得コード（`cfg.Database.Shards`など）が正常に動作することを確認

**テスト実装場所**:
- 既存のE2Eテストスイートに追加

## 7. 実装上の注意事項

### 7.1 viperの設定パス解決

**問題点**:
- 実行ディレクトリによって設定ファイルのパスが変わる可能性がある

**解決策**:
- 複数のパスを`viper.AddConfigPath()`で追加
  - `../config/{env}/`: サーバー実行ディレクトリから見た相対パス
  - `../../config/{env}/`: プロジェクトルートから見た相対パス
  - `./config/{env}/`: カレントディレクトリからの相対パス

### 7.2 viperのインスタンス管理

**注意点**:
- viperはグローバルインスタンスを使用するため、複数の設定ファイルを読み込む際は順序が重要
- `SetConfigName()`でファイル名を変更した後、`MergeInConfig()`を呼び出す

**実装時の注意**:
- メイン設定ファイルを読み込んだ後、データベース設定ファイルをマージする
- ファイル名の変更（`SetConfigName()`）を忘れない

### 7.3 設定ファイルの命名規則

**規則**:
- メイン設定ファイル: `config.yaml`（または`config.yaml.example`）
- データベース設定ファイル: `database.yaml`（または`database.yaml.example`）

**理由**:
- viperは`SetConfigName()`で指定した名前のファイルを検索する
- 拡張子（`.yaml`）は`SetConfigType()`で指定した形式に基づいて自動的に付与される

### 7.4 後方互換性の維持

**重要なポイント**:
- `Config`構造体の定義は変更しない
- 既存の設定取得コード（`cfg.Database.Shards`など）はそのまま動作する
- 環境変数による上書き機能（`DB_PASSWORD_SHARD*`）は維持する

**互換性チェック**:
- 既存のテストがすべて正常に動作することを確認
- 既存のAPIエンドポイントが正常に動作することを確認

### 7.5 設定ファイルの移行手順

**移行手順**:
1. 新しいディレクトリ構造を作成
2. 既存の設定ファイルからデータベース設定を抽出
3. メイン設定ファイルとデータベース設定ファイルを作成
4. `config.go`の`Load()`関数を修正
5. テストを実行して動作確認
6. 既存の設定ファイルを削除

**注意点**:
- 移行中は既存の設定ファイルも残しておき、動作確認後に削除する
- バージョン管理システムで変更履歴を確認できるようにする

## 8. 参考情報

### 8.1 viper公式ドキュメント
- Viper Configuration: https://github.com/spf13/viper
- Viper Multiple Config Files: https://github.com/spf13/viper#reading-multiple-config-files
- Viper MergeInConfig: https://pkg.go.dev/github.com/spf13/viper#MergeInConfig

### 8.2 関連Issue
- GitHub Issue #8: 設定ファイルの分割・リファクタリング

### 8.3 既存ドキュメント
- `server/internal/config/config.go`: 現在の設定読み込み実装
- `config/develop.yaml`: 開発環境設定ファイル
- `config/production.yaml.example`: 本番環境設定テンプレート
- `config/staging.yaml`: ステージング環境設定ファイル

### 8.4 既存実装
- `server/internal/config/config.go`: 設定読み込みロジック
- `server/internal/db/manager.go`: データベース接続管理（設定を使用）

