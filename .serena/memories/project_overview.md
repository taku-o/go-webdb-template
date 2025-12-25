# プロジェクト概要

## 目的
go-webdb-templateは、Go + Next.js + Database Sharding対応のサンプルプロジェクトです。大規模プロジェクト向けの構成を採用しています。

## 技術スタック

### サーバー側（Go）
- **言語**: Go 1.21+
- **ORM**: GORM v1.25.12（Writer/Reader分離対応）
- **ルーティング**: gorilla/mux
- **データベース**: SQLite（開発環境）、PostgreSQL/MySQL（本番想定）
- **設定管理**: spf13/viper（YAML設定ファイル読み込み）
- **管理画面**: GoAdmin

### クライアント側（Next.js）
- **フレームワーク**: Next.js 14+ (App Router)
- **言語**: TypeScript 5+
- **スタイリング**: TailwindCSS

## 主要機能
- Sharding対応: Hash-based shardingで複数DBへデータ分散
- GORM対応: Writer/Reader分離をサポート
- レイヤー分離: API層、Service層、Repository層、DB層で責務を明確化
- 環境別設定: develop/staging/production環境で設定切り替え
- GoAdmin管理画面: ポート8081で管理画面を提供

## 環境設定
環境変数 `APP_ENV` の値に基づいて設定ファイルを読み込み:
- `APP_ENV=develop` → `config/develop.yaml`
- `APP_ENV=staging` → `config/staging.yaml`
- `APP_ENV=production` → `config/production.yaml`
