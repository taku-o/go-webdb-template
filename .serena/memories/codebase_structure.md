# コードベース構造

## ディレクトリ構成

```
go-webdb-template/
├── server/                      # Golangサーバー
│   ├── cmd/
│   │   ├── server/
│   │   │   └── main.go         # メインサーバーエントリーポイント
│   │   └── admin/
│   │       └── main.go         # 管理画面エントリーポイント
│   ├── internal/
│   │   ├── api/                # API定義層
│   │   │   ├── handler/        # HTTPハンドラー
│   │   │   └── router/         # ルーティング
│   │   ├── service/            # ビジネスロジック層
│   │   ├── repository/         # データベース処理層
│   │   ├── model/              # データモデル
│   │   ├── db/                 # DB接続管理
│   │   │   ├── connection.go   # DB接続プール管理
│   │   │   ├── sharding.go     # Sharding戦略
│   │   │   └── manager.go      # DBマネージャー（GORM）
│   │   ├── config/             # 設定読み込み
│   │   └── admin/              # GoAdmin関連
│   │       ├── config.go       # GoAdmin設定
│   │       ├── tables.go       # テーブル定義
│   │       ├── sharding.go     # シャーディング対応
│   │       ├── auth/           # 認証
│   │       └── pages/          # カスタムページ
│   ├── test/                   # テストユーティリティ
│   │   ├── integration/        # 統合テスト
│   │   └── e2e/                # E2Eテスト
│   └── data/                   # SQLiteデータファイル
│
├── client/                      # Next.js + TypeScript
│   ├── src/
│   │   ├── app/                # App Router
│   │   ├── components/         # Reactコンポーネント
│   │   ├── lib/                # API呼び出し等
│   │   └── types/              # TypeScript型定義
│   └── e2e/                    # Playwrightテスト
│
├── config/                      # 環境別設定ファイル
│   ├── develop.yaml
│   ├── staging.yaml
│   └── production.yaml.example
│
├── db/
│   └── migrations/             # マイグレーションSQL
│       ├── shard1/
│       └── shard2/
│
├── docs/                       # ドキュメント
│   ├── Architecture.md
│   ├── API.md
│   ├── Sharding.md
│   ├── Testing.md
│   └── Project-Structure.md
│
└── .kiro/                      # 仕様管理
    ├── steering/               # プロジェクト全体ルール
    └── specs/                  # 機能仕様
```

## レイヤー構成（サーバー側）

1. **API定義層** (`internal/api/`)
   - HTTPリクエスト/レスポンスの処理
   - ルーティング定義

2. **ビジネスロジック層** (`internal/service/`)
   - アプリケーションのコアロジック

3. **データベース処理層** (`internal/repository/`)
   - データベースへのアクセス
   - GORM版とraw SQL版が存在

4. **DB接続管理層** (`internal/db/`)
   - 複数DBシャードへの接続プール管理
   - GORMManagerでWriter/Reader分離をサポート

5. **設定管理層** (`internal/config/`)
   - 環境別設定ファイルの読み込み

## データモデル
- **User**: id, name, email, created_at, updated_at
- **Post**: id, user_id, title, content, created_at, updated_at
